package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/veronica-agent/cans/internal/audio"
	"github.com/veronica-agent/cans/internal/mouth"
)

func TestParseKeepBothOrders(t *testing.T) {
	wav, text, err := parseKeep([]string{"take.wav", "-text", "hello there"})
	if err != nil || wav != "take.wav" || text != "hello there" {
		t.Fatalf("%q %q %v", wav, text, err)
	}
	wav, text, err = parseKeep([]string{"-text", "hello there", "take.wav"})
	if err != nil || wav != "take.wav" || text != "hello there" {
		t.Fatalf("%q %q %v", wav, text, err)
	}
	wav, text, err = parseKeep([]string{"--text=hello there", "take.wav"})
	if err != nil || wav != "take.wav" || text != "hello there" {
		t.Fatalf("%q %q %v", wav, text, err)
	}
}

func TestKeepMissing(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	var buf bytes.Buffer
	old := stderr
	stderr = &buf
	defer func() { stderr = old }()
	if code := run([]string{"keep", "/nope.wav", "-text", "hi"}); code != 2 {
		t.Fatalf("code %d", code)
	}
	msg := buf.String()
	if strings.Contains(msg, "need a wav path") {
		t.Fatalf("flag-order bug still: %s", msg)
	}
	if !strings.Contains(msg, "nope.wav") && !strings.Contains(msg, "not a wav") && !strings.Contains(msg, "no such") {
		t.Fatalf("expected missing-file error, got %q", msg)
	}
}

func TestKeepWavThenText(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	src := filepath.Join(t.TempDir(), "take.wav")
	if err := os.WriteFile(src, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"keep", src, "-text", "hello there"}); code != 0 {
		t.Fatalf("wav then -text: code %d", code)
	}
	if code := run([]string{"keep", "-text", "hello there", src}); code != 0 {
		t.Fatalf("-text then wav: code %d", code)
	}
}

func TestKeepRequiresText(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	src := filepath.Join(t.TempDir(), "take.wav")
	if err := os.WriteFile(src, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	if code := run([]string{"keep", src}); code != 2 {
		t.Fatalf("code %d", code)
	}
}

func TestSayHelp(t *testing.T) {
	if code := run([]string{"help"}); code != 0 {
		t.Fatalf("code %d", code)
	}
}

func TestVersion(t *testing.T) {
	var buf bytes.Buffer
	old := stdout
	stdout = &buf
	defer func() { stdout = old }()
	if code := run([]string{"version"}); code != 0 {
		t.Fatalf("code %d", code)
	}
	if !strings.Contains(buf.String(), "cans") {
		t.Fatalf("%q", buf.String())
	}
}

func TestSayMock(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_ROOT", root)
	t.Setenv("CANS_NOPLAY", "1")
	wav := filepath.Join(root, "voices", "veronica", "ref.wav")
	if err := os.MkdirAll(filepath.Dir(wav), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(wav, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	outWav := filepath.Join(t.TempDir(), "out.wav")
	if err := os.WriteFile(outWav, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(t.TempDir(), "fake-say")
	body := "#!/bin/sh\necho '{\"wav\":\"" + outWav + "\",\"ttfa_ms\":9,\"sample_rate\":24000}'\n"
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CANS_SAY_BIN", script)
	if code := run([]string{"say", "Put the cans on."}); code != 0 {
		t.Fatalf("code %d", code)
	}
}

func TestSayNoTextIsUsage(t *testing.T) {
	var buf bytes.Buffer
	oldErr, oldIn := stderr, stdin
	stderr, stdin = &buf, strings.NewReader("")
	defer func() { stderr, stdin = oldErr, oldIn }()
	if code := run([]string{"say"}); code != 2 {
		t.Fatalf("code %d", code)
	}
	if !strings.Contains(buf.String(), "say: empty text") {
		t.Fatalf("%q", buf.String())
	}
}

func TestSayUnknownFlag(t *testing.T) {
	var buf bytes.Buffer
	old := stderr
	stderr = &buf
	defer func() { stderr = old }()
	if code := run([]string{"say", "--bogus"}); code != 2 {
		t.Fatalf("code %d", code)
	}
	if !strings.Contains(buf.String(), "--bogus") {
		t.Fatalf("%q", buf.String())
	}
}

func TestSayNowaitBusy(t *testing.T) {
	setupFakeMouth(t)
	lk, err := mouth.Acquire(context.Background(), mouth.Path(), 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer lk.Release()
	var buf bytes.Buffer
	old := stderr
	stderr = &buf
	defer func() { stderr = old }()
	if code := run([]string{"say", "--nowait", "Put the cans on."}); code != 75 {
		t.Fatalf("code %d stderr %q", code, buf.String())
	}
	if !strings.Contains(buf.String(), "mouth busy") {
		t.Fatalf("%q", buf.String())
	}
}

func setupFakeMouth(t *testing.T) {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(dir, "..", "..", "internal", "tts", "testdata", "fakeworker")
	bin := filepath.Join(t.TempDir(), "fake-worker")
	out, err := exec.Command("go", "build", "-o", bin, src).CombinedOutput()
	if err != nil {
		t.Fatalf("build fake worker: %s\n%s", err, out)
	}
	home := t.TempDir()
	root := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_ROOT", root)
	t.Setenv("CANS_NOPLAY", "1")
	t.Setenv("CANS_SAY_BIN", "")
	t.Setenv("CANS_WORKER_BIN", bin)
	t.Setenv("CANS_WORKER_MODELS", t.TempDir())
	t.Setenv("CANS_NATIVE_URL", "http://127.0.0.1/cans-test")
	ref := filepath.Join(root, "voices", "veronica", "ref.wav")
	if err := os.MkdirAll(filepath.Dir(ref), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(ref, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "character.toml"), []byte("name = \"veronica\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(filepath.Dir(bin), "libqwen3tts.0.dylib"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
}
