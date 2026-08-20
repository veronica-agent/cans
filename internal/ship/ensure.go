package ship

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:fs
var bundled embed.FS

const embedRoot = "fs"

// Ensure extracts the embedded payload into ~/.cans/shipped when Root is incomplete.
func Ensure() error {
	if Complete(Root()) {
		return nil
	}
	if err := Materialize(Shipped()); err != nil {
		return err
	}
	if !Complete(Shipped()) {
		return fmt.Errorf("ship: extracted payload incomplete at %s", Shipped())
	}
	return nil
}

// Materialize writes the embedded payload into dest.
func Materialize(dest string) error {
	if dest == "" {
		return fmt.Errorf("ship: empty dest")
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	return fs.WalkDir(bundled, embedRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(embedRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dest, filepath.FromSlash(rel))
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := bundled.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}
