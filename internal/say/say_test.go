package say

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"time"

	"github.com/veronica-agent/cans/internal/audio"
	"github.com/veronica-agent/cans/internal/mouth"
	"github.com/veronica-agent/cans/internal/tts"
)

func TestRunMissingText(t *testing.T) {
	o := DefaultOptions()
	o.StdinTTY = true
	var out, errBuf bytes.Buffer
	code := Run(context.Background(), o, failingReader{t}, &out, &errBuf)
	if code != ExitUsage {
		t.Fatalf("code %d, want %d", code, ExitUsage)
	}
	if errBuf.String() != "say: missing text\n" {
		t.Fatalf("stderr %q", errBuf.String())
	}
	if out.Len() != 0 {
		t.Fatalf("stdout %q", out.String())
	}
}

func TestRunOutUnderAFileFails(t *testing.T) {
	fake := sayBinEnv(t)
	blocker := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	o := DefaultOptions()
	o.Text = "Put the cans on."
	o.Out = filepath.Join(blocker, "take.wav")
	var out, errBuf bytes.Buffer
	if code := Run(context.Background(), o, nil, &out, &errBuf); code != ExitFail {
		t.Fatalf("code %d, want %d", code, ExitFail)
	}
	if errBuf.Len() == 0 {
		t.Fatal("expected an error on stderr")
	}
	if out.Len() != 0 {
		t.Fatalf("stdout %q", out.String())
	}
	if _, err := os.Stat(fake.argv); err == nil {
		t.Fatal("say bin ran before mkdir failed")
	}
}

func TestRunOneShotPrintsTTFA(t *testing.T) {
	sayBinEnv(t)
	o := DefaultOptions()
	o.Text = "Put the cans on."
	var out, errBuf bytes.Buffer
	if code := Run(context.Background(), o, nil, &out, &errBuf); code != ExitOK {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if out.String() != "ttfa_ms=12\n" {
		t.Fatalf("stdout %q", out.String())
	}
	if errBuf.Len() != 0 {
		t.Fatalf("stderr %q", errBuf.String())
	}
}

func TestRunOutWritesAndKeeps(t *testing.T) {
	fake := sayBinEnv(t)
	o := DefaultOptions()
	o.Text = "Put the cans on."
	o.Out = filepath.Join(t.TempDir(), "out", "take.wav")
	var out, errBuf bytes.Buffer
	if code := Run(context.Background(), o, nil, &out, &errBuf); code != ExitOK {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if out.String() != o.Out+"\n" {
		t.Fatalf("stdout %q, want %q", out.String(), o.Out+"\n")
	}
	if err := audio.HeaderOK(o.Out); err != nil {
		t.Fatalf("written wav: %v", err)
	}
	if _, err := os.Stat(fake.spoken); err != nil {
		t.Fatalf("say bin wav should survive: %v", err)
	}
}

func TestRunOnceFakeWorkerWritesOut(t *testing.T) {
	bin := buildFakeWorker(t)
	root := t.TempDir()
	t.Setenv("CANS_HOME", t.TempDir())
	t.Setenv("CANS_ROOT", root)
	t.Setenv("CANS_NOPLAY", "1")
	t.Setenv("CANS_SAY_BIN", "")
	t.Setenv("CANS_WORKER_BIN", bin)
	t.Setenv("CANS_WORKER_MODELS", t.TempDir())
	writeRef(t, root)
	o := DefaultOptions()
	o.Out = filepath.Join(t.TempDir(), "out", "take.wav")
	var out, errBuf bytes.Buffer
	if code := runOnce(context.Background(), o, "Put the cans on.", &out, &errBuf); code != ExitOK {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if out.String() != o.Out+"\n" {
		t.Fatalf("stdout %q", out.String())
	}
	if err := audio.HeaderOK(o.Out); err != nil {
		t.Fatalf("written wav: %v", err)
	}
}

// waitFor polls cond until it holds, failing the test after five seconds. It is
// a bounded poll, not a sleep: the test never proceeds on hope.
func waitFor(t *testing.T, cond func() bool, msg string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatal(msg)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// failingReader fails the test if the say flow reads stdin it should not.
type failingReader struct{ t *testing.T }

func (r failingReader) Read([]byte) (int, error) {
	r.t.Fatal("stdin read when it should not be")
	return 0, nil
}

func writeRef(t *testing.T, root string) string {
	t.Helper()
	ref := filepath.Join(root, "voices", "veronica", "ref.wav")
	if err := os.MkdirAll(filepath.Dir(ref), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ref, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	return ref
}

// fakeSay is the CANS_SAY_BIN seam: the wav the script reports, and the file
// it writes its argv to.
type fakeSay struct {
	spoken string
	argv   string
}

// sayBinEnv points CANS_SAY_BIN at a script that records its argv and reports
// a fixed wav.
func sayBinEnv(t *testing.T) fakeSay {
	t.Helper()
	root := t.TempDir()
	t.Setenv("CANS_HOME", t.TempDir())
	t.Setenv("CANS_ROOT", root)
	t.Setenv("CANS_NOPLAY", "1")
	writeRef(t, root)
	dir := t.TempDir()
	f := fakeSay{spoken: filepath.Join(dir, "spoken.wav"), argv: filepath.Join(dir, "argv.txt")}
	if err := os.WriteFile(f.spoken, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(dir, "fake-say")
	body := "#!/bin/sh\nprintf '%s\\n' \"$@\" > " + f.argv +
		"\necho '{\"wav\":\"" + f.spoken + "\",\"ttfa_ms\":12,\"sample_rate\":24000}'\n"
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CANS_SAY_BIN", script)
	return f
}

func buildFakeWorker(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(dir, "..", "tts", "testdata", "fakeworker")
	bin := filepath.Join(t.TempDir(), "fake-worker")
	cmd := exec.Command("go", "build", "-o", bin, src)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build fake worker: %s\n%s", err, out)
	}
	return bin
}

func TestRunStdinIsOneUtterance(t *testing.T) {
	fake := sayBinEnv(t)
	o := DefaultOptions()
	o.Stdin = true
	var out, errBuf bytes.Buffer
	in := strings.NewReader("Put the cans on.\n")
	if code := Run(context.Background(), o, in, &out, &errBuf); code != ExitOK {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	argv, err := os.ReadFile(fake.argv)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(argv), "--text\nPut the cans on.\n") {
		t.Fatalf("say bin argv %q", string(argv))
	}
}

func TestRunJSONRecord(t *testing.T) {
	fake := sayBinEnv(t)
	o := DefaultOptions()
	o.Text = "Put the cans on."
	o.JSON = true
	var out, errBuf bytes.Buffer
	if code := Run(context.Background(), o, nil, &out, &errBuf); code != ExitOK {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(lines) != 1 {
		t.Fatalf("stdout %q, want one record", out.String())
	}
	if strings.Contains(lines[0], `"wav"`) {
		t.Fatalf("wav in record without -o: %s", lines[0])
	}
	var r tts.Result
	if err := json.Unmarshal([]byte(lines[0]), &r); err != nil {
		t.Fatalf("stdout %q: %v", lines[0], err)
	}
	if r.Wav != "" || r.TTFAMs != 12 || r.SampleRate != 24000 {
		t.Fatalf("record %+v", r)
	}
	if _, err := os.Stat(fake.spoken); err == nil {
		t.Fatal("temp wav should be gone after --json without -o")
	}
}

func TestRunNowaitBusy(t *testing.T) {
	fakeWorkerEnv(t)
	lk, err := mouth.Acquire(context.Background(), mouth.Path(), 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer lk.Release()
	o := DefaultOptions()
	o.Text = "Put the cans on."
	o.Wait = 0
	var out, errBuf bytes.Buffer
	if code := Run(context.Background(), o, nil, &out, &errBuf); code != ExitBusy {
		t.Fatalf("code %d stderr %q, want %d", code, errBuf.String(), ExitBusy)
	}
	if !strings.Contains(errBuf.String(), "mouth busy") {
		t.Fatalf("stderr %q", errBuf.String())
	}
	if strings.Contains(errBuf.String(), "waiting for the mouth") {
		t.Fatalf("nowait waited: %q", errBuf.String())
	}
	if out.Len() != 0 {
		t.Fatalf("stdout %q", out.String())
	}
}

func TestRunWaitBusyPrintsWaiting(t *testing.T) {
	fakeWorkerEnv(t)
	lk, err := mouth.Acquire(context.Background(), mouth.Path(), 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer lk.Release()
	o := DefaultOptions()
	o.Text = "Put the cans on."
	o.Wait = 200 * time.Millisecond
	var out, errBuf bytes.Buffer
	start := time.Now()
	if code := Run(context.Background(), o, nil, &out, &errBuf); code != ExitBusy {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if elapsed := time.Since(start); elapsed < 200*time.Millisecond {
		t.Fatalf("elapsed %v", elapsed)
	}
	if n := strings.Count(errBuf.String(), "waiting for the mouth…"); n != 1 {
		t.Fatalf("waiting lines %d in %q", n, errBuf.String())
	}
	if out.Len() != 0 {
		t.Fatalf("stdout %q", out.String())
	}
}

func TestRunInterruptedWaiting(t *testing.T) {
	fakeWorkerEnv(t)
	lk, err := mouth.Acquire(context.Background(), mouth.Path(), 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer lk.Release()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	o := DefaultOptions()
	o.Text = "Put the cans on."
	o.Wait = -1
	var out, errBuf bytes.Buffer
	code := Run(ctx, o, nil, &out, &errBuf)
	if code != ExitInterrupted {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if !strings.Contains(errBuf.String(), "say: interrupted") {
		t.Fatalf("stderr %q", errBuf.String())
	}
}

// TestRunOnceCancelledMidSynthesis is D014 for one-shot: the worker is held
// inside synthesis, so cancel lands mid-utterance rather than at the lock.
func TestRunOnceCancelledMidSynthesis(t *testing.T) {
	fakeWorkerEnv(t)
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	marker := filepath.Join(tmp, "blocked")
	t.Setenv("CANS_FAKE_BLOCK_FILE", marker)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	o := DefaultOptions()
	var out, errBuf bytes.Buffer

	done := make(chan int, 1)
	go func() { done <- runOnce(ctx, o, "block", &out, &errBuf) }()
	waitFor(t, func() bool {
		_, err := os.Stat(marker)
		return err == nil
	}, "worker never reached synthesis")
	cancel()
	code := <-done
	if code != ExitInterrupted {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if !strings.Contains(errBuf.String(), "say: interrupted") {
		t.Fatalf("stderr %q", errBuf.String())
	}
	if out.Len() != 0 {
		t.Fatalf("stdout %q", out.String())
	}
	left, err := filepath.Glob(filepath.Join(tmp, "cans-*.wav"))
	if err != nil {
		t.Fatal(err)
	}
	if len(left) != 0 {
		t.Fatalf("temp wavs left: %v", left)
	}
	lk, err := mouth.Acquire(context.Background(), mouth.Path(), 0, nil)
	if err != nil {
		t.Fatalf("lock after cancel: %v", err)
	}
	lk.Release()
}

func TestRunAfterReleaseSucceeds(t *testing.T) {
	fakeWorkerEnv(t)
	lk, err := mouth.Acquire(context.Background(), mouth.Path(), 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := lk.Release(); err != nil {
		t.Fatal(err)
	}
	o := DefaultOptions()
	o.Text = "Put the cans on."
	o.Wait = 0
	var out, errBuf bytes.Buffer
	if code := Run(context.Background(), o, nil, &out, &errBuf); code != ExitOK {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
}

func fakeWorkerEnv(t *testing.T) {
	t.Helper()
	bin := buildFakeWorker(t)
	home := t.TempDir()
	root := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_ROOT", root)
	t.Setenv("CANS_NOPLAY", "1")
	t.Setenv("CANS_SAY_BIN", "")
	t.Setenv("CANS_WORKER_BIN", bin)
	t.Setenv("CANS_WORKER_MODELS", t.TempDir())
	t.Setenv("CANS_NATIVE_URL", "http://127.0.0.1/cans-test")
	writeRef(t, root)
	if err := os.WriteFile(filepath.Join(root, "character.toml"), []byte("name = \"veronica\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dylib := filepath.Join(filepath.Dir(bin), "libqwen3tts.0.dylib")
	if err := os.WriteFile(dylib, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRunJSONWithOutCarriesTheOutPath(t *testing.T) {
	sayBinEnv(t)
	o := DefaultOptions()
	o.Text = "Put the cans on."
	o.JSON = true
	o.Out = filepath.Join(t.TempDir(), "out", "take.wav")
	var out, errBuf bytes.Buffer
	if code := Run(context.Background(), o, nil, &out, &errBuf); code != ExitOK {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	var r tts.Result
	if err := json.Unmarshal(out.Bytes(), &r); err != nil {
		t.Fatalf("stdout %q: %v", out.String(), err)
	}
	if r.Wav != o.Out {
		t.Fatalf("record wav %q, want %q", r.Wav, o.Out)
	}
	if err := audio.HeaderOK(o.Out); err != nil {
		t.Fatalf("written wav: %v", err)
	}
}
