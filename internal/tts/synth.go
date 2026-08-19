package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/veronica-agent/cans/internal/keep"
)

// Result is one spoken line.
type Result struct {
	Wav        string `json:"wav"`
	TTFAMs     int    `json:"ttfa_ms"`
	SampleRate int    `json:"sample_rate"`
}

func sidecarBin() string {
	if b := os.Getenv("CANS_SAY_BIN"); b != "" {
		return b
	}
	root := os.Getenv("CANS_ROOT")
	if root == "" {
		wd, _ := os.Getwd()
		root = wd
	}
	return filepath.Join(root, "sidecar", "say.py")
}

// Say clones text with the current throat.
func Say(text string) (Result, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return Result{}, fmt.Errorf("say: empty text")
	}
	cur, err := keep.Load()
	if err != nil {
		return Result{}, err
	}
	bin := sidecarBin()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".py") {
		cmd = exec.Command("uv", "run", "python", bin, "--text", text, "--ref", cur.Wav, "--ref-text", cur.RefText)
		if root := os.Getenv("CANS_ROOT"); root != "" {
			cmd.Dir = root
		}
	} else {
		cmd = exec.Command(bin, "--text", text, "--ref", cur.Wav, "--ref-text", cur.RefText)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return Result{}, fmt.Errorf("say: %s", msg)
	}
	line := lastJSONLine(stdout.Bytes())
	if line == nil {
		return Result{}, fmt.Errorf("say: sidecar printed no json (%q)", stdout.String())
	}
	var r Result
	if err := json.Unmarshal(line, &r); err != nil {
		return Result{}, fmt.Errorf("say: bad sidecar json: %w (%q)", err, stdout.String())
	}
	if r.Wav == "" {
		return Result{}, fmt.Errorf("say: sidecar returned no wav")
	}
	return r, nil
}

func lastJSONLine(raw []byte) []byte {
	var last []byte
	for _, line := range bytes.Split(bytes.TrimSpace(raw), []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) > 0 && line[0] == '{' {
			last = line
		}
	}
	return last
}
