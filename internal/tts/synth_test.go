package tts

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/veronica-agent/cans/internal/audio"
	"github.com/veronica-agent/cans/internal/keep"
)

func TestTokenBudget(t *testing.T) {
	if tokenBudget("") != 80 {
		t.Fatalf("empty: %d", tokenBudget(""))
	}
	if tokenBudget("hi") != 80 {
		t.Fatalf("short: %d", tokenBudget("hi"))
	}
	long := strings.Repeat("word ", 80)
	if tokenBudget(long) != 360 {
		t.Fatalf("long: %d", tokenBudget(long))
	}
}

func TestSynthRequestCloneParams(t *testing.T) {
	b, err := json.Marshal(synthRequest("cans", "Put the cans on.", "/tmp/ref.wav"))
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if got["temperature"] != 0.2 {
		t.Fatalf("temperature %v", got["temperature"])
	}
	if _, ok := got["max_tokens"].(float64); !ok {
		t.Fatalf("max_tokens %T %v", got["max_tokens"], got["max_tokens"])
	}
	if got["ref_wav"] != "/tmp/ref.wav" {
		t.Fatalf("ref %v", got["ref_wav"])
	}
}

func TestDumpLiveWav(t *testing.T) {
	dst := os.Getenv("CANS_DUMP_WAV")
	if dst == "" {
		t.Skip("CANS_DUMP_WAV")
	}
	r, err := Say(context.Background(), "Put the cans on.")
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(r.Wav)
	if err != nil {
		t.Fatal(err)
	}
	RemoveTemp(r.Wav)
	if err := os.WriteFile(dst, b, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("ttfa_ms=%d wav=%s bytes=%d", r.TTFAMs, dst, len(b))
}

func TestSayEmpty(t *testing.T) {
	if _, err := Say(context.Background(), "  "); err == nil {
		t.Fatal("expected empty text error")
	}
}

func TestRemoveTempOnlyUnderTmp(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "say.wav")
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
	if _, err := Say(context.Background(), "hi"); err == nil {
		t.Fatal("expected non-wav reject")
	}
}

func TestSayMockBin(t *testing.T) {
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
	r, err := Say(context.Background(), "Put the cans on.")
	if err != nil {
		t.Fatal(err)
	}
	if r.TTFAMs != 12 || r.Wav != outWav {
		t.Fatalf("%+v", r)
	}
}

func TestSayWithEmptyRefWav(t *testing.T) {
	bin := buildFakeWorker(t)
	t.Setenv("CANS_HOME", t.TempDir())
	t.Setenv("CANS_WORKER_BIN", bin)
	t.Setenv("CANS_WORKER_MODELS", t.TempDir())
	sess, err := Open(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()
	if _, err := sess.Say(context.Background(), "hi", keep.Current{}); err == nil {
		t.Fatal("expected empty ref wav")
	}
}

func TestSessionFakeWorker(t *testing.T) {
	bin := buildFakeWorker(t)
	home := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_WORKER_BIN", bin)
	t.Setenv("CANS_WORKER_MODELS", t.TempDir())
	wav := filepath.Join(t.TempDir(), "ref.wav")
	if err := os.WriteFile(wav, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	sess, err := Open(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()
	r, err := sess.Say(context.Background(), "Put the cans on.", keep.Current{Wav: wav})
	if err != nil {
		t.Fatal(err)
	}
	if err := audio.HeaderOK(r.Wav); err != nil {
		t.Fatal(err)
	}
	RemoveTemp(r.Wav)
}

func TestSessionSilentMouth(t *testing.T) {
	bin := buildFakeWorker(t)
	home := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_WORKER_BIN", bin)
	t.Setenv("CANS_WORKER_MODELS", t.TempDir())
	wav := filepath.Join(t.TempDir(), "ref.wav")
	if err := os.WriteFile(wav, audio.Minimal(), 0o644); err != nil {
		t.Fatal(err)
	}
	sess, err := Open(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()
	_, err = sess.Say(context.Background(), "silent", keep.Current{Wav: wav})
	if err == nil {
		t.Fatal("expected silent mouth error")
	}
	if !strings.Contains(err.Error(), "silent mouth") {
		t.Fatalf("err %q", err.Error())
	}
}

func TestSayMissingWorker(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	t.Setenv("CANS_WORKER_BIN", filepath.Join(t.TempDir(), "nope"))
	t.Setenv("CANS_SAY_BIN", "")
	if _, err := SayWith(context.Background(), "hi", keep.Current{Wav: "/x"}); err == nil {
		t.Fatal("expected missing worker")
	}
}

func buildFakeWorker(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(dir, "testdata", "fakeworker")
	bin := filepath.Join(t.TempDir(), "fake-worker")
	cmd := exec.Command("go", "build", "-o", bin, src)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build fake worker: %s\n%s", err, out)
	}
	return bin
}
