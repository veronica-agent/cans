package play

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/veronica-agent/cans/internal/audio"
)

func TestFileRejectsGarbage(t *testing.T) {
	t.Setenv("CANS_NOPLAY", "1")
	p := filepath.Join(t.TempDir(), "hosts")
	if err := os.WriteFile(p, []byte("not audio"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := File(p); err == nil {
		t.Fatal("expected reject")
	}
}

func TestFileAcceptsWAV(t *testing.T) {
	t.Setenv("CANS_NOPLAY", "1")
	p := filepath.Join(t.TempDir(), "a.wav")
	if err := os.WriteFile(p, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := File(p); err != nil {
		t.Fatal(err)
	}
}
