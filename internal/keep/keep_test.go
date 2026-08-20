package keep

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/veronica-agent/cans/internal/audio"
)

func writeWAV(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadDefaultMissing(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	t.Setenv("CANS_ROOT", t.TempDir())
	if _, err := Load(); err == nil {
		t.Fatal("expected missing default ref")
	}
}

func TestPinRequiresText(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	src := filepath.Join(t.TempDir(), "take.wav")
	writeWAV(t, src)
	if _, err := Pin(src, "  "); err == nil {
		t.Fatal("expected -text required")
	}
}

func TestPinAndLoad(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_ROOT", root)
	writeWAV(t, filepath.Join(root, "voices", "veronica", "ref.wav"))
	src := filepath.Join(t.TempDir(), "take.wav")
	writeWAV(t, src)
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

func TestLoadCorruptEmpty(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_ROOT", t.TempDir())
	if err := os.WriteFile(filepath.Join(home, "current.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(); err == nil {
		t.Fatal("expected corrupt state")
	}
}

func TestLoadRejectsOutsidePath(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_ROOT", root)
	writeWAV(t, filepath.Join(root, "voices", "veronica", "ref.wav"))
	hosts := filepath.Join(root, "hosts.wav")
	writeWAV(t, hosts)
	body := `{"wav":"` + hosts + `","ref_text":"nope"}`
	if err := os.WriteFile(filepath.Join(home, "current.json"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(); err == nil {
		t.Fatal("expected rejected path")
	}
}

func TestPinMissing(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	if _, err := Pin("/no/such.wav", "hi"); err == nil {
		t.Fatal("expected error")
	}
}

func TestQuote(t *testing.T) {
	root := t.TempDir()
	t.Setenv("CANS_ROOT", root)
	if err := os.WriteFile(filepath.Join(root, "character.toml"), []byte("quote = \"Wine's already poured.\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if Quote() != "Wine's already poured." {
		t.Fatalf("%q", Quote())
	}
}
