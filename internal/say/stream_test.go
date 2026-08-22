package say

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/veronica-agent/cans/internal/mouth"
)

func TestStreamFailContinues(t *testing.T) {
	fakeWorkerEnv(t)
	dir := t.TempDir()
	o := DefaultOptions()
	o.Stream = true
	o.Out = filepath.Join(dir, "%03d.wav")
	var out, errBuf bytes.Buffer
	code := Run(context.Background(), o, strings.NewReader("a\nfail\n\nb\n"), &out, &errBuf)
	if code != ExitFail {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "001.wav")); err != nil {
		t.Fatalf("001.wav: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "002.wav")); err != nil {
		t.Fatalf("002.wav: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "003.wav")); err == nil {
		t.Fatal("003.wav should not exist")
	}
	paths := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(paths) != 2 {
		t.Fatalf("stdout %q", out.String())
	}
	if !strings.Contains(errBuf.String(), "line 2:") {
		t.Fatalf("stderr %q", errBuf.String())
	}
}

func TestStreamJSONRecords(t *testing.T) {
	fakeWorkerEnv(t)
	dir := t.TempDir()
	o := DefaultOptions()
	o.Stream = true
	o.JSON = true
	o.Out = filepath.Join(dir, "%03d.wav")
	var out, errBuf bytes.Buffer
	code := Run(context.Background(), o, strings.NewReader("a\nfail\n\nb\n"), &out, &errBuf)
	if code != ExitFail {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("stdout %q", out.String())
	}
	var recs []map[string]any
	for i, line := range lines {
		var rec map[string]any
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			t.Fatalf("record %d %q: %v", i, line, err)
		}
		recs = append(recs, rec)
	}
	if recs[0]["line"] != float64(1) || recs[1]["line"] != float64(2) || recs[2]["line"] != float64(4) {
		t.Fatalf("lines %+v", recs)
	}
	if recs[1]["error"] == nil {
		t.Fatalf("middle record %+v", recs[1])
	}
	if recs[0]["wav"] != filepath.Join(dir, "001.wav") {
		t.Fatalf("wav %+v", recs[0]["wav"])
	}
}

func TestStreamJSONWithoutOutOmitsWav(t *testing.T) {
	fakeWorkerEnv(t)
	o := DefaultOptions()
	o.Stream = true
	o.JSON = true
	var out, errBuf bytes.Buffer
	code := Run(context.Background(), o, strings.NewReader("Put the cans on.\n"), &out, &errBuf)
	if code != ExitOK {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if strings.Contains(out.String(), `"wav"`) {
		t.Fatalf("wav in record without -o: %s", out.String())
	}
	var rec map[string]any
	if err := json.Unmarshal(out.Bytes(), &rec); err != nil {
		t.Fatalf("stdout %q: %v", out.String(), err)
	}
	if rec["line"] != float64(1) || rec["ttfa_ms"] == nil {
		t.Fatalf("record %+v", rec)
	}
}

func TestStreamOneWorker(t *testing.T) {
	fakeWorkerEnv(t)
	dir := t.TempDir()
	count := countingWorker(t, dir)
	o := DefaultOptions()
	o.Stream = true
	o.Out = filepath.Join(dir, "%03d.wav")
	var out, errBuf bytes.Buffer
	code := Run(context.Background(), o, strings.NewReader("a\nb\nc\nd\ne\n"), &out, &errBuf)
	if code != ExitOK {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	wantOneStart(t, count)
}

func TestStreamBusyDoesNotReadStdin(t *testing.T) {
	fakeWorkerEnv(t)
	lk, err := mouth.Acquire(context.Background(), mouth.Path(), 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer lk.Release()
	o := DefaultOptions()
	o.Stream = true
	o.Wait = 0
	var out, errBuf bytes.Buffer
	code := Run(context.Background(), o, failingReader{t}, &out, &errBuf)
	if code != ExitBusy {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
}

type syncBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (w *syncBuffer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}

func (w *syncBuffer) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.String()
}

func TestStreamCancel(t *testing.T) {
	fakeWorkerEnv(t)
	dir := t.TempDir()
	count := countingWorker(t, dir)

	pr, pw := io.Pipe()
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })
	var out syncBuffer
	var errBuf bytes.Buffer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	o := DefaultOptions()
	o.Stream = true
	o.Out = filepath.Join(dir, "%03d.wav")

	done := make(chan int, 1)
	go func() {
		done <- Run(ctx, o, pr, &out, &errBuf)
	}()
	if _, err := pw.Write([]byte("a\n")); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool { return out.String() != "" }, "no first record")
	cancel()
	// Off the test goroutine: once Run has returned nothing reads the pipe,
	// so this write would block forever.
	go func() { _, _ = pw.Write([]byte("b\n")) }()
	code := <-done
	_ = pw.CloseWithError(io.EOF)
	if code != ExitInterrupted {
		t.Fatalf("code %d stderr %q stdout %q", code, errBuf.String(), out.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "001.wav")); err != nil {
		t.Fatalf("001.wav: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "002.wav")); err == nil {
		t.Fatal("002.wav should not exist")
	}
	if !strings.Contains(errBuf.String(), "interrupted after line 1") {
		t.Fatalf("stderr %q", errBuf.String())
	}
	lk, err := mouth.Acquire(context.Background(), mouth.Path(), 0, nil)
	if err != nil {
		t.Fatalf("lock after cancel: %v", err)
	}
	lk.Release()
	wantOneStart(t, count)
}

func TestStreamEmptyStdinOK(t *testing.T) {
	fakeWorkerEnv(t)
	o := DefaultOptions()
	o.Stream = true
	var out, errBuf bytes.Buffer
	code := Run(context.Background(), o, strings.NewReader(""), &out, &errBuf)
	if code != ExitOK {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if out.Len() != 0 {
		t.Fatalf("stdout %q", out.String())
	}
}

func TestStreamLongLineReportsStdin(t *testing.T) {
	fakeWorkerEnv(t)
	dir := t.TempDir()
	o := DefaultOptions()
	o.Stream = true
	o.Out = filepath.Join(dir, "%03d.wav")
	long := strings.Repeat("x", 1<<20+1)
	var out, errBuf bytes.Buffer
	code := Run(context.Background(), o, strings.NewReader("a\n"+long+"\n"), &out, &errBuf)
	if code != ExitFail {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if !strings.Contains(errBuf.String(), "say: stdin: ") {
		t.Fatalf("stderr %q", errBuf.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "001.wav")); err != nil {
		t.Fatalf("the line before the bad one should survive: %v", err)
	}
}

// TestStreamCancelBeforeAnyLine holds the worker mid-synthesis on the first
// line, so cancel lands before anything is spoken (D014).
func TestStreamCancelBeforeAnyLine(t *testing.T) {
	fakeWorkerEnv(t)
	dir := t.TempDir()
	t.Setenv("CANS_FAKE_BLOCK_FILE", filepath.Join(dir, "blocked"))
	pr, pw := io.Pipe()
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	o := DefaultOptions()
	o.Stream = true
	o.Out = filepath.Join(dir, "%03d.wav")
	var out syncBuffer
	var errBuf bytes.Buffer

	done := make(chan int, 1)
	go func() { done <- Run(ctx, o, pr, &out, &errBuf) }()
	if _, err := pw.Write([]byte("block\n")); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool {
		_, err := os.Stat(filepath.Join(dir, "blocked"))
		return err == nil
	}, "worker never reached synthesis")
	cancel()
	code := <-done
	if code != ExitInterrupted {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if !strings.Contains(errBuf.String(), "interrupted before the first line") {
		t.Fatalf("stderr %q", errBuf.String())
	}
	if out.String() != "" {
		t.Fatalf("stdout %q", out.String())
	}
	if _, err := os.Stat(filepath.Join(dir, "001.wav")); err == nil {
		t.Fatal("001.wav should not exist")
	}
	lk, err := mouth.Acquire(context.Background(), mouth.Path(), 0, nil)
	if err != nil {
		t.Fatalf("lock after cancel: %v", err)
	}
	lk.Release()
}

// TestStreamCancelAfterFailedLine pins the other half of D014's N: a line that
// was reported as failed still counts as fully processed.
func TestStreamCancelAfterFailedLine(t *testing.T) {
	fakeWorkerEnv(t)
	dir := t.TempDir()
	t.Setenv("CANS_FAKE_BLOCK_FILE", filepath.Join(dir, "blocked"))
	pr, pw := io.Pipe()
	t.Cleanup(func() { _ = pw.Close(); _ = pr.Close() })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	o := DefaultOptions()
	o.Stream = true
	o.Out = filepath.Join(dir, "%03d.wav")
	var out syncBuffer
	var errBuf syncBuffer

	done := make(chan int, 1)
	go func() { done <- Run(ctx, o, pr, &out, &errBuf) }()
	if _, err := pw.Write([]byte("fail\nblock\n")); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool {
		_, err := os.Stat(filepath.Join(dir, "blocked"))
		return err == nil
	}, "worker never reached synthesis")
	cancel()
	if code := <-done; code != ExitInterrupted {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if !strings.Contains(errBuf.String(), "line 1:") {
		t.Fatalf("stderr %q", errBuf.String())
	}
	if !strings.Contains(errBuf.String(), "interrupted after line 1") {
		t.Fatalf("stderr %q", errBuf.String())
	}
}

// TestStreamNoOutPlaysEachLine is D007: no -o means play every line and keep no
// wav behind. CANS_NOPLAY makes play a header check.
func TestStreamNoOutPlaysEachLine(t *testing.T) {
	fakeWorkerEnv(t)
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	o := DefaultOptions()
	o.Stream = true
	var out, errBuf bytes.Buffer
	code := Run(context.Background(), o, strings.NewReader("a\n\nb\n"), &out, &errBuf)
	if code != ExitOK {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if errBuf.Len() != 0 {
		t.Fatalf("stderr %q", errBuf.String())
	}
	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("stdout %q", out.String())
	}
	for i, line := range lines {
		if !strings.HasPrefix(line, "ttfa_ms=") {
			t.Fatalf("line %d %q", i, line)
		}
	}
	left, err := filepath.Glob(filepath.Join(tmp, "cans-*.wav"))
	if err != nil {
		t.Fatal(err)
	}
	if len(left) != 0 {
		t.Fatalf("temp wavs left: %v", left)
	}
}

// countingWorker wraps the fake worker in a script that appends one line per
// start, and returns the path of that counter.
func countingWorker(t *testing.T, dir string) string {
	t.Helper()
	bin := os.Getenv("CANS_WORKER_BIN")
	count := filepath.Join(dir, "starts")
	wrap := filepath.Join(filepath.Dir(bin), "wrap")
	body := "#!/bin/sh\nprintf 'x\\n' >> " + count + "\nexec " + bin + " \"$@\"\n"
	if err := os.WriteFile(wrap, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CANS_WORKER_BIN", wrap)
	return count
}

// wantOneStart asserts the mouth was opened exactly once: one worker, ever.
func wantOneStart(t *testing.T, count string) {
	t.Helper()
	got, err := os.ReadFile(count)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(got), "\n") != 1 {
		t.Fatalf("starts %q", got)
	}
}
