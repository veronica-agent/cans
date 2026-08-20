package ship

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

const (
	nativeURL = "https://github.com/Obedience-Corp/qwen3-tts-native/releases/download/v0.1.0/qwen3-tts-native-darwin-arm64.tar.gz"
	nativeSHA = "67b6e7514168d146cf4d892d6594559796bcea4f130c9c837beb8eac1c416e2f"
)

// NativeReady is true when the worker binary is on disk.
func NativeReady() bool {
	_, err := os.Stat(WorkerBin())
	return err == nil
}

func dylibReady() bool {
	_, err := os.Stat(filepath.Join(filepath.Dir(WorkerBin()), "libqwen3tts.0.dylib"))
	return err == nil
}

// EnsureNative downloads and unpacks qwen3-tts-native if the worker is missing.
func EnsureNative(ctx context.Context, log io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if os.Getenv("CANS_SAY_BIN") != "" {
		return nil
	}
	fixDylibNames(filepath.Join(Native(), "bin"))
	if NativeReady() && dylibReady() {
		return nil
	}
	url := os.Getenv("CANS_NATIVE_URL")
	sha := os.Getenv("CANS_NATIVE_SHA256")
	if url == "" {
		if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
			return fmt.Errorf("native mouth: Apple Silicon macOS only")
		}
		url = nativeURL
		sha = nativeSHA
	}
	if log == nil {
		log = io.Discard
	}
	fmt.Fprintln(log, "fetching native mouth (~1.6 GB, once)…")
	root := Native()
	if err := os.MkdirAll(root, 0o755); err != nil {
		return err
	}
	tmp, err := downloadArchive(ctx, root, url, sha)
	if err != nil {
		return fmt.Errorf("native mouth: %w", err)
	}
	defer os.Remove(tmp)
	if err := extractTarGz(tmp, root); err != nil {
		return fmt.Errorf("native mouth: %w", err)
	}
	if !NativeReady() {
		return fmt.Errorf("native mouth: unpack did not produce %s", WorkerBin())
	}
	_ = os.Chmod(WorkerBin(), 0o755)
	fixDylibNames(filepath.Join(Native(), "bin"))
	return nil
}
