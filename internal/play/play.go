package play

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// File plays a wav. CANS_NOPLAY=1 skips (tests).
func File(path string) error {
	if os.Getenv("CANS_NOPLAY") == "1" {
		return nil
	}
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("play: %w", err)
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("afplay", path)
	default:
		cmd = exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", path)
	}
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
