package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/veronica-agent/cans/internal/audio"
	"github.com/veronica-agent/cans/internal/keep"
)

func sayBin(text string, cur keep.Current) (Result, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return Result{}, fmt.Errorf("say: empty text")
	}
	bin := os.Getenv("CANS_SAY_BIN")
	if bin == "" {
		return Result{}, fmt.Errorf("say: empty CANS_SAY_BIN")
	}
	cmd := exec.Command(bin, "--text", text, "--ref", cur.Wav, "--ref-text", cur.RefText)
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
		return Result{}, fmt.Errorf("say: say bin printed no json (%q)", stdout.String())
	}
	var r Result
	if err := json.Unmarshal(line, &r); err != nil {
		return Result{}, fmt.Errorf("say: bad say bin json: %w (%q)", err, stdout.String())
	}
	if r.Wav == "" {
		return Result{}, fmt.Errorf("say: say bin returned no wav")
	}
	if err := audio.HeaderOK(r.Wav); err != nil {
		return Result{}, fmt.Errorf("say: %w", err)
	}
	return r, nil
}

func lastJSONLine(raw []byte) []byte {
	var last []byte
	for _, line := range bytes.Split(bytes.TrimSpace(raw), []byte("\n")) {
		line = bytes.TrimSpace(line)
		var r Result
		if json.Unmarshal(line, &r) == nil && r.Wav != "" {
			last = line
		}
	}
	return last
}
