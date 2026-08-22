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

// SayWith clones text using a frozen throat into a temp wav.
func SayWith(ctx context.Context, text string, cur keep.Current) (Result, error) {
	return SayTo(ctx, text, cur, "")
}

// SayTo clones text using a frozen throat. out == "" writes a temp wav;
// otherwise the wav is written at out and is the caller's to delete.
func SayTo(ctx context.Context, text string, cur keep.Current, out string) (Result, error) {
	return SayToWith(ctx, text, cur, out, DefaultOptions())
}

// SayToWith is SayTo with lock wait options. CANS_SAY_BIN takes no lock.
func SayToWith(ctx context.Context, text string, cur keep.Current, out string, o Options) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if os.Getenv("CANS_SAY_BIN") != "" {
		return sayBinTo(text, cur, out)
	}
	sess, err := OpenWith(ctx, o)
	if err != nil {
		return Result{}, err
	}
	defer sess.Close()
	return sess.SayTo(ctx, text, cur, out)
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
