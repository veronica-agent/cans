package tts

import (
	"context"
	"fmt"
	"os"
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

// Say clones text with the current throat.
func Say(ctx context.Context, text string) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return Result{}, fmt.Errorf("say: empty text")
	}
	cur, err := keep.Load()
	if err != nil {
		return Result{}, err
	}
	return SayWith(ctx, text, cur)
}

// SayWith clones text using a frozen throat.
func SayWith(ctx context.Context, text string, cur keep.Current) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if os.Getenv("CANS_SAY_BIN") != "" {
		return sayBin(text, cur)
	}
	sess, err := Open(ctx)
	if err != nil {
		return Result{}, err
	}
	defer sess.Close()
	return sess.Say(ctx, text, cur)
}

// RemoveTemp deletes a synth wav if it lives under the process temp dir.
func RemoveTemp(path string) {
	if path == "" {
		return
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return
	}
	tmp, err := filepath.Abs(os.TempDir())
	if err != nil {
		return
	}
	rel, err := filepath.Rel(tmp, abs)
	if err != nil || strings.HasPrefix(rel, "..") {
		return
	}
	_ = os.Remove(abs)
}
