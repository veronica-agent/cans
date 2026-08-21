package ship

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/veronica-agent/cans/internal/audio"
)

// payload is every file the embed ships, relative to the checkout root.
var payload = []string{
	"character.toml",
	"voices/veronica/meta.json",
	"voices/veronica/ref.wav",
}

func TestCompleteFalseOnEmpty(t *testing.T) {
	if Complete("") || Complete(t.TempDir()) {
		t.Fatal("expected incomplete")
	}
}

func TestCompleteNeedsEveryPiece(t *testing.T) {
	for _, missing := range []string{"character.toml", "voices/veronica/ref.wav"} {
		t.Run(missing, func(t *testing.T) {
			dest := t.TempDir()
			if err := Materialize(dest); err != nil {
				t.Fatal(err)
			}
			if err := os.Remove(filepath.Join(dest, filepath.FromSlash(missing))); err != nil {
				t.Fatal(err)
			}
			if Complete(dest) {
				t.Fatalf("complete without %s", missing)
			}
		})
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
	char, err := os.ReadFile(filepath.Join(dest, "character.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(char, []byte(`name = "Veronica"`)) {
		t.Fatalf("character.toml missing name: %s", char)
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
	if !currentPayload(Shipped()) {
		t.Fatal("expected stamp")
	}
	if Root() == Shipped() {
		t.Fatal("CANS_ROOT must still win even when incomplete")
	}
}

func TestEnsureRefreshesStaleShipped(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_ROOT", t.TempDir())
	if err := Ensure(); err != nil {
		t.Fatal(err)
	}
	char := filepath.Join(Shipped(), "character.toml")
	if err := os.WriteFile(char, []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(Shipped(), stampFile), []byte("old\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Ensure(); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(char)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(got, []byte("stale")) {
		t.Fatal("expected rematerialize")
	}
	if !currentPayload(Shipped()) {
		t.Fatal("stamp should match embed")
	}
}

func TestEnsureSkipsCompleteRoot(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CANS_HOME", home)
	root := t.TempDir()
	if err := Materialize(root); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CANS_ROOT", root)
	if err := Ensure(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(Shipped()); !os.IsNotExist(err) {
		t.Fatalf("shipped should not exist: %v", err)
	}
}

func TestRootHonorsCANS_ROOT(t *testing.T) {
	root := t.TempDir()
	t.Setenv("CANS_ROOT", root)
	if got := Root(); got != root {
		t.Fatalf("got %s", got)
	}
}

// TestEmbedIsExactlyPayload fails if anything sneaks into internal/ship/fs
// that is not on the payload list — the Python sidecar must not come back.
func TestEmbedIsExactlyPayload(t *testing.T) {
	var got []string
	err := fs.WalkDir(bundled, embedRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(embedRoot, path)
		if err != nil {
			return err
		}
		got = append(got, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)
	want := append([]string(nil), payload...)
	sort.Strings(want)
	if len(got) != len(want) {
		t.Fatalf("embed has %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("embed has %v, want %v", got, want)
		}
	}
}

func TestEmbedMatchesCheckout(t *testing.T) {
	root, ok := checkoutRoot()
	if !ok {
		t.Skip("not in a checkout")
	}
	dest := t.TempDir()
	if err := Materialize(dest); err != nil {
		t.Fatal(err)
	}
	for _, rel := range payload {
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

// checkoutRoot walks up to the directory holding go.mod and character.toml.
// A CANS_ROOT payload never has go.mod, so the pair is unambiguous.
func checkoutRoot() (string, bool) {
	wd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	dir := wd
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			if _, err := os.Stat(filepath.Join(dir, "character.toml")); err == nil {
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
