package tts

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const workerProtocol = "qwen3-tts-worker/v1"

// Client is a live qwen3-tts-worker process.
type Client struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
}

type pcmResult struct {
	sampleRate int
	samples    []float32
	wall       time.Duration
}

// StartWorker launches workerBin with modelDir. Honor ctx.
func StartWorker(ctx context.Context, workerBin, modelDir string) (*Client, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	workerBin, err := filepath.Abs(workerBin)
	if err != nil {
		return nil, err
	}
	modelDir, err = filepath.Abs(modelDir)
	if err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(ctx, workerBin, modelDir)
	libDir := filepath.Dir(workerBin)
	cmd.Env = withLibPath(os.Environ(), libDir)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = io.Discard
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	c := &Client{cmd: cmd, stdin: stdin, stdout: bufio.NewReader(stdout)}
	if err := c.readReady(ctx); err != nil {
		_ = c.Close()
		return nil, err
	}
	return c, nil
}

func (c *Client) readReady(ctx context.Context) error {
	line, err := readLine(ctx, c.stdout)
	if err != nil {
		return fmt.Errorf("ready: %w", err)
	}
	var head struct {
		Type     string `json:"type"`
		Protocol string `json:"protocol"`
	}
	if err := json.Unmarshal([]byte(line), &head); err != nil {
		return fmt.Errorf("ready json: %w (%q)", err, line)
	}
	if head.Type != "ready" {
		return fmt.Errorf("expected ready, got %s", head.Type)
	}
	if head.Protocol != "" && head.Protocol != workerProtocol {
		return fmt.Errorf("protocol %s", head.Protocol)
	}
	return nil
}

func (c *Client) synthesize(ctx context.Context, id, text, refWAV string) (*pcmResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	req := map[string]string{"type": "synthesize", "id": id, "text": text}
	if refWAV != "" {
		req["ref_wav"] = refWAV
	}
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	start := time.Now()
	if _, err := c.stdin.Write(append(b, '\n')); err != nil {
		return nil, err
	}
	return c.readPCM(ctx, start)
}

func (c *Client) Close() error {
	if c.stdin != nil {
		_, _ = c.stdin.Write([]byte(`{"type":"shutdown"}` + "\n"))
		_ = c.stdin.Close()
	}
	if c.cmd != nil && c.cmd.Process != nil {
		return c.cmd.Wait()
	}
	return nil
}

func withLibPath(env []string, libDir string) []string {
	out := make([]string, 0, len(env)+2)
	for _, e := range env {
		key, _, ok := strings.Cut(e, "=")
		if !ok || key == "DYLD_LIBRARY_PATH" || key == "LD_LIBRARY_PATH" {
			if !ok {
				out = append(out, e)
			}
			continue
		}
		out = append(out, e)
	}
	sep := string(os.PathListSeparator)
	for _, key := range []string{"DYLD_LIBRARY_PATH", "LD_LIBRARY_PATH"} {
		prev := os.Getenv(key)
		if prev != "" {
			out = append(out, key+"="+libDir+sep+prev)
		} else {
			out = append(out, key+"="+libDir)
		}
	}
	return out
}

func readLine(ctx context.Context, r *bufio.Reader) (string, error) {
	type result struct {
		line string
		err  error
	}
	ch := make(chan result, 1)
	go func() {
		line, err := r.ReadString('\n')
		ch <- result{line, err}
	}()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-ch:
		return strings.TrimSpace(res.line), res.err
	}
}
