package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKeepMissing(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	if code := run([]string{"keep", "/nope.wav"}); code != 2 {
		t.Fatalf("code %d", code)
	}
}

func TestSayHelp(t *testing.T) {
	if code := run([]string{"help"}); code != 0 {
		t.Fatalf("code %d", code)
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
	if err := os.WriteFile(wav, []byte("RIFF"), 0o644); err != nil {
		t.Fatal(err)
	}
	outWav := filepath.Join(t.TempDir(), "out.wav")
	if err := os.WriteFile(outWav, []byte("RIFF"), 0o644); err != nil {
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
