package ship

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/veronica-agent/cans/internal/audio"
)

// Version is set at release via -X ldflags. Dev builds stay "dev".
var Version = "dev"

// Home is ~/.cans (or CANS_HOME).
func Home() string {
	if h := os.Getenv("CANS_HOME"); h != "" {
		return filepath.Clean(h)
	}
	u, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(u) == "" {
		return ".cans"
	}
	return filepath.Join(u, ".cans")
}

// Shipped is the extracted payload directory.
func Shipped() string { return filepath.Join(Home(), "shipped") }

// Venv is the leftover uv tree. Product speech does not use it.
func Venv() string { return filepath.Join(Home(), "venv") }

// Native is the qwen3-tts-worker install root.
func Native() string { return filepath.Join(Home(), "native") }

// WorkerBin is qwen3-tts-worker. CANS_WORKER_BIN overrides.
func WorkerBin() string {
	if b := os.Getenv("CANS_WORKER_BIN"); b != "" {
		return b
	}
	return filepath.Join(Native(), "bin", "qwen3-tts-worker")
}

// WorkerModels is the GGUF + presets directory. CANS_WORKER_MODELS overrides.
func WorkerModels() string {
	if d := os.Getenv("CANS_WORKER_MODELS"); d != "" {
		return d
	}
	return filepath.Join(Native(), "models")
}

// HFHome is the Hugging Face cache for the 0.6B clone model.
func HFHome() string { return filepath.Join(Home(), "hf") }

// Complete reports whether root has the payload cans needs to speak.
func Complete(root string) bool {
	if strings.TrimSpace(root) == "" {
		return false
	}
	if audio.HeaderOK(filepath.Join(root, "voices", "veronica", "ref.wav")) != nil {
		return false
	}
	for _, rel := range []string{"sidecar/say.py", "pyproject.toml", "character.toml"} {
		st, err := os.Stat(filepath.Join(root, filepath.FromSlash(rel)))
		if err != nil || st.IsDir() {
			return false
		}
	}
	return true
}

// Env overlays UV_PROJECT_ENVIRONMENT and HF_HOME on a copy of base.
func Env(base []string) []string {
	out := make([]string, 0, len(base)+2)
	for _, e := range base {
		if strings.HasPrefix(e, "UV_PROJECT_ENVIRONMENT=") || strings.HasPrefix(e, "HF_HOME=") {
			continue
		}
		out = append(out, e)
	}
	return append(out,
		"UV_PROJECT_ENVIRONMENT="+Venv(),
		"HF_HOME="+HFHome(),
	)
}

// VenvReady is true after a successful uv sync.
func VenvReady() bool {
	_, err := os.Stat(filepath.Join(Venv(), "pyvenv.cfg"))
	return err == nil
}

// Sidecar is the clone script under root.
func Sidecar(root string) string {
	return filepath.Join(root, "sidecar", "say.py")
}

// DefaultWav is the shipped Veronica clip under root.
func DefaultWav(root string) string {
	return filepath.Join(root, "voices", "veronica", "ref.wav")
}
