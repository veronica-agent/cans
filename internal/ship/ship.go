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

// Complete reports whether root has the payload cans needs to speak:
// the Veronica clip and character.toml.
func Complete(root string) bool {
	if strings.TrimSpace(root) == "" {
		return false
	}
	if audio.HeaderOK(DefaultWav(root)) != nil {
		return false
	}
	st, err := os.Stat(filepath.Join(root, "character.toml"))
	return err == nil && !st.IsDir()
}

// DefaultWav is the shipped Veronica clip under root.
func DefaultWav(root string) string {
	return filepath.Join(root, "voices", "veronica", "ref.wav")
}
