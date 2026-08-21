package tts

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/veronica-agent/cans/internal/audio"
	"github.com/veronica-agent/cans/internal/keep"
	"github.com/veronica-agent/cans/internal/ship"
)

// Session is a warm worker for one booth (or one say).
type Session struct {
	c *Client
}

// Open starts the native worker.
func Open(ctx context.Context) (*Session, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	bin := ship.WorkerBin()
	models := ship.WorkerModels()
	if _, err := os.Stat(bin); err != nil {
		return nil, fmt.Errorf("say: native mouth missing — cans doctor")
	}
	c, err := StartWorker(ctx, bin, models)
	if err != nil {
		return nil, fmt.Errorf("say: worker: %w", err)
	}
	return &Session{c: c}, nil
}

// Say clones text using the frozen throat.
func (s *Session) Say(ctx context.Context, text string, cur keep.Current) (Result, error) {
	if s == nil || s.c == nil {
		return Result{}, fmt.Errorf("say: no worker")
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return Result{}, fmt.Errorf("say: empty text")
	}
	if strings.TrimSpace(cur.Wav) == "" {
		return Result{}, fmt.Errorf("say: empty ref wav")
	}
	pcm, err := s.c.synthesize(ctx, "cans", text, cur.Wav)
	if err != nil {
		return Result{}, fmt.Errorf("say: %w", err)
	}
	out := filepath.Join(os.TempDir(), fmt.Sprintf("cans-%d.wav", time.Now().UnixNano()))
	rate := pcm.sampleRate
	if rate <= 0 {
		rate = 24000
	}
	pcm.samples = audio.Clean(pcm.samples, rate)
	if err := audio.WritePCM16(out, rate, pcm.samples); err != nil {
		return Result{}, err
	}
	ms := int(pcm.wall.Milliseconds())
	return Result{Wav: out, TTFAMs: ms, SampleRate: rate}, nil
}

// Close shuts the worker down.
func (s *Session) Close() error {
	if s == nil || s.c == nil {
		return nil
	}
	return s.c.Close()
}
