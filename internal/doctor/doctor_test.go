package doctor

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/veronica-agent/cans/internal/audio"
	"github.com/veronica-agent/cans/internal/ship"
)

func TestMachineCheck(t *testing.T) {
	c := machineCheck()
	want := runtime.GOOS == "darwin" && runtime.GOARCH == "arm64"
	if c.OK != want {
		t.Fatalf("ok=%v want %v (%s)", c.OK, want, c.Detail)
	}
}

func TestDiagnoseMissingPayload(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	t.Setenv("CANS_ROOT", t.TempDir())
	checks := Diagnose(context.Background())
	by := map[string]Check{}
	for _, c := range checks {
		by[c.Name] = c
	}
	if by["payload"].OK {
		t.Fatal("payload should fail")
	}
	if by["throat"].OK {
		t.Fatal("throat should fail")
	}
}

func TestDiagnoseAfterMaterialize(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_ROOT", root)
	if err := ship.Materialize(root); err != nil {
		t.Fatal(err)
	}
	checks := Diagnose(context.Background())
	by := map[string]Check{}
	for _, c := range checks {
		by[c.Name] = c
	}
	if !by["payload"].OK {
		t.Fatalf("payload: %+v", by["payload"])
	}
	if !by["throat"].OK {
		t.Fatalf("throat: %+v", by["throat"])
	}
}

func TestPrepareSkipsSyncWhenSayBinSet(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	t.Setenv("CANS_ROOT", t.TempDir())
	t.Setenv("CANS_SAY_BIN", "/bin/true")
	if err := Prepare(context.Background(), os.Stderr); err != nil {
		t.Fatal(err)
	}
	if ship.VenvReady() {
		t.Fatal("should not have synced")
	}
}

func TestThroatUsesPinnedWav(t *testing.T) {
	home := t.TempDir()
	root := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_ROOT", root)
	if err := ship.Materialize(root); err != nil {
		t.Fatal(err)
	}
	wav := filepath.Join(home, "current", "ref.wav")
	if err := os.MkdirAll(filepath.Dir(wav), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(wav, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	body := `{"wav":"` + wav + `","ref_text":"hello"}`
	if err := os.WriteFile(filepath.Join(home, "current.json"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	c := throatCheck()
	if !c.OK || c.Detail != wav {
		t.Fatalf("%+v", c)
	}
}
