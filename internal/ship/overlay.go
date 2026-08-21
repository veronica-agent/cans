package ship

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func applyWorkerOverlay() error {
	if len(workerOverlay) == 0 {
		return nil
	}
	if os.Getenv("CANS_NATIVE_URL") != "" {
		return nil
	}
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		return nil
	}
	dst := WorkerBin()
	if cur, err := os.ReadFile(dst); err == nil && bytes.Equal(cur, workerOverlay) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	tmp := dst + ".overlay"
	if err := os.WriteFile(tmp, workerOverlay, 0o755); err != nil {
		return fmt.Errorf("worker overlay: %w", err)
	}
	if err := os.Rename(tmp, dst); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("worker overlay: %w", err)
	}
	_ = os.Chmod(dst, 0o755)
	return nil
}
