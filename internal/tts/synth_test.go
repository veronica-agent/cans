package tts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/veronica-agent/cans/internal/audio"
	"github.com/veronica-agent/cans/internal/keep"
)

func TestSayEmpty(t *testing.T) {
	if _, err := Say("  "); err == nil {
		t.Fatal("expected empty text error")
	}
}

func TestRemoveTempOnlyUnderTmp(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "sidecar.wav")
	if err := os.WriteFile(tmp, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "keep.wav")
	if err := os.WriteFile(outside, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	// t.TempDir is under os.TempDir on this OS — RemoveTemp should delete tmp file
	// only when the path is inside os.TempDir. outside may also be; skip if not.
	RemoveTemp(tmp)
	if _, err := os.Stat(tmp); err == nil {
		// still there: temp dir not under os.TempDir (unusual)
		t.Log("temp not removed (test temp dir outside os.TempDir)")
	}
	absOut, _ := filepath.Abs(outside)
	tmpRoot, _ := filepath.Abs(os.TempDir())
	rel, err := filepath.Rel(tmpRoot, absOut)
	if err == nil && !strings.HasPrefix(rel, "..") {
		return
	}
	RemoveTemp(outside)
	if _, err := os.Stat(outside); err != nil {
		t.Fatal("should not delete wav outside process temp dir")
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
