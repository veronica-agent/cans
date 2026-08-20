package ship

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Sync installs sidecar deps into Venv with the locked pyproject.
func Sync(ctx context.Context, out io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := Ensure(); err != nil {
		return err
	}
	if _, err := exec.LookPath("uv"); err != nil {
		return fmt.Errorf("need uv on PATH — brew install uv")
	}
	root := Root()
	if !Complete(root) {
		root = Shipped()
	}
	if !Complete(root) {
		return fmt.Errorf("ship: no payload to sync")
	}
	if out == nil {
		out = io.Discard
	}
	cmd := exec.CommandContext(ctx, "uv", "sync", "--project", root, "--locked", "--python", "3.12")
	cmd.Env = Env(os.Environ())
	cmd.Dir = root
	cmd.Stdout = out
	cmd.Stderr = out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("uv sync: %w", err)
	}
	return nil
}

// LookUV is the uv binary on PATH, or empty.
func LookUV() string {
	p, err := exec.LookPath("uv")
	if err != nil {
		return ""
	}
	return p
}

// ImportMLX runs a locked uv check that mlx_audio imports.
func ImportMLX(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	root := Root()
	if !Complete(root) {
		root = Shipped()
	}
	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "uv", "run", "--project", root, "--locked", "python", "-c", "import mlx_audio")
	cmd.Env = Env(os.Environ())
	cmd.Dir = root
	cmd.Stdout = io.Discard
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("mlx-audio: %s", msg)
	}
	return nil
}
