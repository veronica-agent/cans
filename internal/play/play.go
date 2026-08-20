package play

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/veronica-agent/cans/internal/audio"
)

// File plays a wav. CANS_NOPLAY=1 skips after validating the file.
func File(path string) error {
	if err := audio.HeaderOK(path); err != nil {
		return fmt.Errorf("play: %w", err)
	}
	if os.Getenv("CANS_NOPLAY") == "1" {
		return nil
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("afplay", path)
	default:
		cmd = exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "error", path)
	}
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
