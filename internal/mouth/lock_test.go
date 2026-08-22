package mouth

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestAcquireWaitZeroBusy(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mouth.lock")
	held := mustAcquire(t, path, -1)
	defer held.Release()
	_, err := Acquire(context.Background(), path, 0, nil)
	if !errors.Is(err, ErrBusy) {
		t.Fatalf("err %v, want ErrBusy", err)
	}
}

func TestAcquireBoundedWaitBusy(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mouth.lock")
	held := mustAcquire(t, path, -1)
	defer held.Release()
	start := time.Now()
	_, err := Acquire(context.Background(), path, 150*time.Millisecond, nil)
	elapsed := time.Since(start)
	if !errors.Is(err, ErrBusy) {
		t.Fatalf("err %v, want ErrBusy", err)
	}
	if elapsed < 150*time.Millisecond {
		t.Fatalf("elapsed %v, want >= 150ms", elapsed)
	}
}

func TestAcquireCtxCancel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mouth.lock")
	held := mustAcquire(t, path, -1)
	defer held.Release()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := Acquire(ctx, path, -1, nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err %v, want ctx deadline", err)
	}
}

func TestOnWaitFiresOnce(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mouth.lock")
	held := mustAcquire(t, path, -1)
	defer held.Release()
	var n atomic.Int32
	_, err := Acquire(context.Background(), path, 350*time.Millisecond, func() { n.Add(1) })
	if !errors.Is(err, ErrBusy) {
		t.Fatalf("err %v, want ErrBusy", err)
	}
	if n.Load() != 1 {
		t.Fatalf("onWait %d, want 1", n.Load())
	}
}

func TestReleaseThenReacquireKeepsFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mouth.lock")
	held := mustAcquire(t, path, 0)
	if err := held.Release(); err != nil {
		t.Fatal(err)
	}
	again := mustAcquire(t, path, 0)
	defer again.Release()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("lock file missing: %v", err)
	}
}

func TestReleaseNil(t *testing.T) {
	var l *Lock
	if err := l.Release(); err != nil {
		t.Fatal(err)
	}
}

func TestHelperHoldLock(t *testing.T) {
	if os.Getenv("MOUTH_HELPER") != "1" {
		return
	}
	l, err := Acquire(context.Background(), os.Getenv("MOUTH_LOCK"), -1, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "helper acquire: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("held")
	_ = os.Stdout.Sync()
	defer runtime.KeepAlive(l)
	select {}
}

func TestKillDropsLock(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mouth.lock")
	cmd := exec.Command(os.Args[0], "-test.run=^TestHelperHoldLock$")
	cmd.Env = append(os.Environ(), "MOUTH_HELPER=1", "MOUTH_LOCK="+path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	line, err := bufio.NewReader(stdout).ReadString('\n')
	if err != nil || line != "held\n" {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		t.Fatalf("helper hello %q err %v", line, err)
	}
	_, err = Acquire(context.Background(), path, 0, nil)
	if !errors.Is(err, ErrBusy) {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		t.Fatalf("while held: %v", err)
	}
	if err := cmd.Process.Kill(); err != nil {
		t.Fatal(err)
	}
	_ = cmd.Wait()
	l, err := Acquire(context.Background(), path, 0, nil)
	if err != nil {
		t.Fatalf("after kill: %v", err)
	}
	defer l.Release()
}

func TestPathUsesHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CANS_HOME", home)
	if Path() != filepath.Join(home, "mouth.lock") {
		t.Fatalf("Path() %q", Path())
	}
}

func mustAcquire(t *testing.T, path string, wait time.Duration) *Lock {
	t.Helper()
	l, err := Acquire(context.Background(), path, wait, nil)
	if err != nil {
		t.Fatal(err)
	}
	return l
}
