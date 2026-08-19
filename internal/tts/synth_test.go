package tts

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/veronica-agent/cans/internal/audio"
	"github.com/veronica-agent/cans/internal/keep"
)

func TestSayEmpty(t *testing.T) {
	if _, err := Say("  "); err == nil {
		t.Fatal("expected empty text error")
	}
}

func TestLastJSONLineIgnoresJunk(t *testing.T) {
	raw := []byte("Initialized encoder\n{\"noise\":true}\n{\"wav\":\"/tmp/x.wav\",\"ttfa_ms\":12,\"sample_rate\":24000}\n")
	got := lastJSONLine(raw)
	if string(got) != `{"wav":"/tmp/x.wav","ttfa_ms":12,"sample_rate":24000}` {
		t.Fatalf("%q", got)
	}
}

func TestSayRejectsNonWAV(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_ROOT", root)
	def := filepath.Join(root, "voices", "veronica", "ref.wav")
	if err := os.MkdirAll(filepath.Dir(def), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(def, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	hosts := filepath.Join(t.TempDir(), "hosts")
	if err := os.WriteFile(hosts, []byte("root:x:0:0"), 0o644); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(t.TempDir(), "fake-say")
	body := "#!/bin/sh\necho '{\"wav\":\"" + hosts + "\",\"ttfa_ms\":1,\"sample_rate\":24000}'\n"
	if err := os.WriteFile(script, []byte(body), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CANS_SAY_BIN", script)
	if _, err := Say("hi"); err == nil {
		t.Fatal("expected non-wav reject")
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
	if err := os.WriteFile(wav, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	outWav := filepath.Join(t.TempDir(), "out.wav")
	if err := os.WriteFile(outWav, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(t.TempDir(), "fake-say")
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

func TestSayWithEmptyRefText(t *testing.T) {
	if _, err := SayWith("hi", keep.Current{Wav: "/x", RefText: ""}); err == nil {
		t.Fatal("expected empty ref text")
	}
}
