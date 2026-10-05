package play

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/veronica-agent/cans/internal/audio"
)

// File plays a wav. CANS_NOPLAY=1 skips after validating the file.
func File(path string) error {
	return FileContext(context.Background(), path)
}

// FileContext plays a wav until it finishes or ctx is canceled.
func FileContext(ctx context.Context, path string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := audio.HeaderOK(path); err != nil {
		return fmt.Errorf("play: %w", err)
	}
	if os.Getenv("CANS_NOPLAY") == "1" {
		return nil
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.CommandContext(ctx, "afplay", path)
	default:
		cmd = exec.CommandContext(ctx, "ffplay", "-nodisp", "-autoexit", "-loglevel", "error", path)
	}
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}
