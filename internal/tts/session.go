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
	"github.com/veronica-agent/cans/internal/mouth"
	"github.com/veronica-agent/cans/internal/ship"
)

// Session is a warm worker for one booth (or one say). The mouth lock lives
// as long as the session does.
type Session struct {
	c    *Client
	lock *mouth.Lock
	// done is the session context's Done channel, kept so Close can tell an
	// expected terminate-on-cancel (D014) from a real worker failure.
	done <-chan struct{}
}

// Options bound how long Open waits for the mouth.
type Options struct {
	Wait   time.Duration
	OnWait func()
}

// DefaultOptions wait forever and print a stderr line while blocked.
func DefaultOptions() Options {
	return Options{Wait: -1, OnWait: defaultOnWait}
}

func defaultOnWait() {
	fmt.Fprintln(os.Stderr, "waiting for the mouth…")
}

// Open starts the native worker, waiting forever for the mouth if needed.
func Open(ctx context.Context) (*Session, error) {
	return OpenWith(ctx, DefaultOptions())
}

// OpenWith acquires the mouth lock, then starts the worker.
func OpenWith(ctx context.Context, o Options) (*Session, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	lk, err := mouth.Acquire(ctx, mouth.Path(), o.Wait, o.OnWait)
	if err != nil {
		return nil, err
	}
	sess, err := startSession(ctx, lk)
	if err != nil {
		_ = lk.Release()
		return nil, err
	}
	return sess, nil
}

func startSession(ctx context.Context, lock *mouth.Lock) (*Session, error) {
	bin := ship.WorkerBin()
	models := ship.WorkerModels()
	if _, err := os.Stat(bin); err != nil {
		return nil, fmt.Errorf("say: native mouth missing — cans doctor")
	}
	c, err := StartWorker(ctx, bin, models)
	if err != nil {
		return nil, fmt.Errorf("say: worker: %w", err)
	}
	return &Session{c: c, lock: lock, done: ctx.Done()}, nil
}

// Say clones text using the frozen throat into a temp wav.
func (s *Session) Say(ctx context.Context, text string, cur keep.Current) (Result, error) {
	return s.SayTo(ctx, text, cur, "")
}

// SayTo clones text using the frozen throat. out == "" writes a temp wav;
// otherwise the wav is written at out, parent directories included.
func (s *Session) SayTo(ctx context.Context, text string, cur keep.Current, out string) (Result, error) {
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
	if out == "" {
		out = filepath.Join(os.TempDir(), fmt.Sprintf("cans-%d.wav", time.Now().UnixNano()))
	} else if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return Result{}, fmt.Errorf("say: %w", err)
	}
	pcm, err := s.c.synthesize(ctx, "cans", text, cur.Wav)
	if err != nil {
		return Result{}, fmt.Errorf("say: %w", err)
	}
	rate := pcm.sampleRate
	if rate <= 0 {
		rate = 24000
	}
	pcm.samples = audio.Clean(pcm.samples, rate)
	if audio.Silent(pcm.samples, rate) {
		pcm, err = s.c.synthesize(ctx, "cans", text, cur.Wav)
		if err != nil {
			return Result{}, fmt.Errorf("say: %w", err)
		}
		if rate != pcm.sampleRate && pcm.sampleRate > 0 {
			rate = pcm.sampleRate
		}
		pcm.samples = audio.Clean(pcm.samples, rate)
		if audio.Silent(pcm.samples, rate) {
			return Result{}, fmt.Errorf("say: silent mouth")
		}
	}
	pcm.samples = audio.Normalize(pcm.samples, 0.5)
	if err := audio.WritePCM16(out, rate, pcm.samples); err != nil {
		return Result{}, err
	}
	ms := int(pcm.wall.Milliseconds())
	return Result{Wav: out, TTFAMs: ms, SampleRate: rate}, nil
}

// Close shuts the worker down, then releases the mouth lock. After a cancel
// the worker is signalled rather than asked (D014), so the refused shutdown
// write and the `signal: terminated` exit are expected, not failures.
func (s *Session) Close() error {
	if s == nil {
		return nil
	}
	var err error
	if s.c != nil {
		err = s.c.Close()
		s.c = nil
	}
	if s.cancelled() {
		err = nil
	}
	if s.lock != nil {
		rerr := s.lock.Release()
		s.lock = nil
		if err == nil {
			err = rerr
		}
	}
	return err
}

// cancelled reports whether the context the session was opened with is done.
func (s *Session) cancelled() bool {
	select {
	case <-s.done:
		return true
	default:
		return false
	}
}
