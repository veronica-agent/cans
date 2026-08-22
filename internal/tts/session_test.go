package tts

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/veronica-agent/cans/internal/mouth"
)

func TestOpenWithLockBeforeWorker(t *testing.T) {
	bin := buildFakeWorker(t)
	home := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_WORKER_BIN", bin)
	t.Setenv("CANS_WORKER_MODELS", t.TempDir())
	t.Setenv("CANS_NOPLAY", "1")
	t.Setenv("CANS_SAY_BIN", "")

	a, err := OpenWith(context.Background(), Options{Wait: -1})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(home, "mouth.lock")); err != nil {
		t.Fatalf("lock file: %v", err)
	}

	t.Setenv("CANS_WORKER_BIN", filepath.Join(t.TempDir(), "missing"))
	_, err = OpenWith(context.Background(), Options{Wait: 0})
	if !errors.Is(err, mouth.ErrBusy) {
		t.Fatalf("want ErrBusy, got %v", err)
	}

	if err := a.Close(); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CANS_WORKER_BIN", bin)
	b, err := OpenWith(context.Background(), Options{Wait: 0})
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestOpenWithOnWaitOnce(t *testing.T) {
	bin := buildFakeWorker(t)
	t.Setenv("CANS_HOME", t.TempDir())
	t.Setenv("CANS_WORKER_BIN", bin)
	t.Setenv("CANS_WORKER_MODELS", t.TempDir())
	t.Setenv("CANS_NOPLAY", "1")
	t.Setenv("CANS_SAY_BIN", "")

	a, err := OpenWith(context.Background(), Options{Wait: -1})
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()

	var n atomic.Int32
	_, err = OpenWith(context.Background(), Options{
		Wait:   300 * time.Millisecond,
		OnWait: func() { n.Add(1) },
	})
	if !errors.Is(err, mouth.ErrBusy) {
		t.Fatalf("want ErrBusy, got %v", err)
	}
	if n.Load() != 1 {
		t.Fatalf("onWait %d, want 1", n.Load())
	}
}
