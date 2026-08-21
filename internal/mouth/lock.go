// Package mouth serializes the native worker: one flock, one resident mouth.
package mouth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/veronica-agent/cans/internal/ship"
)

// ErrBusy is a held mouth when the caller asked not to wait, or the wait ran out.
var ErrBusy = errors.New("mouth busy")

const pollEvery = 100 * time.Millisecond

// Lock is an exclusive flock on mouth.lock. The kernel drops it if we die.
type Lock struct{ f *os.File }

// Path is CANS_HOME/mouth.lock.
func Path() string {
	return filepath.Join(ship.Home(), "mouth.lock")
}

// Acquire takes an exclusive flock on path.
// wait < 0 waits forever; wait == 0 tries once; wait > 0 gives up after wait.
// onWait fires once only when the caller is about to block.
func Acquire(ctx context.Context, path string, wait time.Duration, onWait func()) (*Lock, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("mouth: lock %s: %w", path, err)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, fmt.Errorf("mouth: lock %s: %w", path, err)
	}
	l := &Lock{f: f}
	if err := tryUntil(ctx, path, l, wait, onWait); err != nil {
		_ = f.Close()
		return nil, err
	}
	return l, nil
}

// Release drops the flock and closes the fd. The lock file is never deleted.
func (l *Lock) Release() error {
	if l == nil || l.f == nil {
		return nil
	}
	flockErr := syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
	closeErr := l.f.Close()
	l.f = nil
	if flockErr != nil {
		return fmt.Errorf("mouth: unlock: %w", flockErr)
	}
	return closeErr
}

func tryUntil(ctx context.Context, path string, l *Lock, wait time.Duration, onWait func()) error {
	var deadline time.Time
	if wait > 0 {
		deadline = time.Now().Add(wait)
	}
	called := false
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		err := syscall.Flock(int(l.f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			return nil
		}
		if err != syscall.EWOULDBLOCK && err != syscall.EAGAIN {
			return fmt.Errorf("mouth: lock %s: %w", path, err)
		}
		if wait == 0 {
			return ErrBusy
		}
		if wait > 0 && !time.Now().Before(deadline) {
			return ErrBusy
		}
		if !called && onWait != nil {
			onWait()
			called = true
		}
		d := pollEvery
		if wait > 0 {
			if rem := time.Until(deadline); rem <= 0 {
				return ErrBusy
			} else if rem < d {
				d = rem
			}
		}
		if err := waitPoll(ctx, d); err != nil {
			return err
		}
	}
}

func waitPoll(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}
