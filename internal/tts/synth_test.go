package tts

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/veronica-agent/cans/internal/keep"
)

func TestSayEmpty(t *testing.T) {
	if _, err := Say("  "); err == nil {
		t.Fatal("expected empty text error")
	}
}

func TestSayMockSidecar(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_ROOT", root)
	wav := filepath.Join(root, "voices", "veronica", "ref.wav")
	if err := os.MkdirAll(filepath.Dir(wav), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(wav, []byte("RIFF"), 0o644); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(t.TempDir(), "fake-say")
	outWav := filepath.Join(t.TempDir(), "out.wav")
	if err := os.WriteFile(outWav, []byte("RIFF"), 0o644); err != nil {
		t.Fatal(err)
	}
	body := "#!/bin/sh\necho '{\"wav\":\"" + outWav + "\",\"ttfa_ms\":12,\"sample_rate\":24000}'\n"
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CANS_SAY_BIN", script)
	if _, err := keep.Load(); err != nil {
		t.Fatal(err)
	}
	r, err := Say("Put the cans on.")
	if err != nil {
		t.Fatal(err)
	}
	if r.TTFAMs != 12 || r.Wav != outWav {
		t.Fatalf("%+v", r)
	}
}
