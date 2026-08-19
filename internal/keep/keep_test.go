package keep

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultMissing(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	t.Setenv("CANS_ROOT", t.TempDir())
	if _, err := Load(); err == nil {
		t.Fatal("expected missing default ref")
	}
}

func TestPinAndLoad(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_ROOT", root)
	src := filepath.Join(t.TempDir(), "take.wav")
	if err := os.WriteFile(src, []byte("RIFF"), 0o644); err != nil {
		t.Fatal(err)
	}
	c, err := Pin(src, "hello")
	if err != nil {
		t.Fatal(err)
	}
	got, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.Wav != c.Wav || got.RefText != "hello" {
		t.Fatalf("got %+v", got)
	}
}

func TestPinMissing(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	if _, err := Pin("/no/such.wav", ""); err == nil {
		t.Fatal("expected error")
	}
}
