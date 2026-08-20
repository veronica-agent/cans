package ship

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/veronica-agent/cans/internal/audio"
)

func TestCompleteFalseOnEmpty(t *testing.T) {
	if Complete("") || Complete(t.TempDir()) {
		t.Fatal("expected incomplete")
	}
}

func TestMaterializeThenComplete(t *testing.T) {
	dest := t.TempDir()
	if err := Materialize(dest); err != nil {
		t.Fatal(err)
	}
	if !Complete(dest) {
		t.Fatal("expected complete payload")
	}
	if err := audio.HeaderOK(DefaultWav(dest)); err != nil {
		t.Fatal(err)
	}
	py, err := os.ReadFile(Sidecar(dest))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(py, []byte("Qwen3-TTS")) {
		t.Fatalf("sidecar missing clone model: %s", py[:min(80, len(py))])
	}
}

func TestEnsureUsesShippedWhenRootIncomplete(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_ROOT", t.TempDir())
	if Complete(Root()) {
		t.Fatal("CANS_ROOT should be incomplete")
	}
	if err := Ensure(); err != nil {
		t.Fatal(err)
	}
	if !Complete(Shipped()) {
		t.Fatal("expected shipped payload")
	}
	if Root() == Shipped() {
		t.Fatal("CANS_ROOT must still win even when incomplete")
	}
}

func TestRootHonorsCANS_ROOT(t *testing.T) {
	root := t.TempDir()
	t.Setenv("CANS_ROOT", root)
	if got := Root(); got != root {
		t.Fatalf("got %s", got)
	}
}

func TestEnvOverridesHFAndVenv(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	got := Env([]string{"PATH=/bin", "HF_HOME=/nope", "UV_PROJECT_ENVIRONMENT=/nope"})
	wantHF := "HF_HOME=" + HFHome()
	wantUV := "UV_PROJECT_ENVIRONMENT=" + Venv()
	var sawHF, sawUV, sawOld bool
	for _, e := range got {
		switch e {
		case wantHF:
			sawHF = true
		case wantUV:
			sawUV = true
		case "HF_HOME=/nope", "UV_PROJECT_ENVIRONMENT=/nope":
			sawOld = true
		}
	}
	if !sawHF || !sawUV || sawOld {
		t.Fatalf("%v", got)
	}
}

func TestEmbedMatchesCheckout(t *testing.T) {
	root, ok := checkoutRoot()
	if !ok {
		t.Skip("not in a checkout")
	}
	pairs := []string{
		"character.toml",
		"pyproject.toml",
		"uv.lock",
		"sidecar/say.py",
		"voices/veronica/ref.wav",
		"voices/veronica/meta.json",
	}
	dest := t.TempDir()
	if err := Materialize(dest); err != nil {
		t.Fatal(err)
	}
	for _, rel := range pairs {
		want, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(filepath.Join(dest, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(want, got) {
			t.Fatalf("embed drifted from checkout: %s (run just dist stamp)", rel)
		}
	}
}

func checkoutRoot() (string, bool) {
	wd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	dir := wd
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			if _, err := os.Stat(filepath.Join(dir, "sidecar", "say.py")); err == nil {
				return dir, true
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}
