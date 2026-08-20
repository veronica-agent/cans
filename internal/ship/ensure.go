package ship

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed all:fs
var bundled embed.FS

const embedRoot = "fs"
const stampFile = ".stamp"

// Ensure extracts the embedded payload into ~/.cans/shipped when there is no
// live checkout (or CANS_ROOT / brew share). A stamp of the embed refreshes
// ~/.cans/shipped after brew --HEAD or a new binary.
func Ensure() error {
	if livePayload() {
		return nil
	}
	dest := Shipped()
	if currentPayload(dest) {
		return nil
	}
	if err := Materialize(dest); err != nil {
		return err
	}
	if err := writeStamp(dest); err != nil {
		return err
	}
	if !Complete(dest) {
		return fmt.Errorf("ship: extracted payload incomplete at %s", dest)
	}
	return nil
}

func livePayload() bool {
	if r := os.Getenv("CANS_ROOT"); r != "" {
		return Complete(filepath.Clean(r))
	}
	shipped := Shipped()
	for _, start := range searchStarts() {
		dir := start
		for i := 0; i < 8; i++ {
			if Complete(dir) && dir != shipped {
				return true
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return false
}

func currentPayload(dest string) bool {
	if !Complete(dest) {
		return false
	}
	got, err := os.ReadFile(filepath.Join(dest, stampFile))
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(got)) == fingerprint()
}

func writeStamp(dest string) error {
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dest, stampFile), []byte(fingerprint()+"\n"), 0o644)
}

func fingerprint() string {
	var files []string
	err := fs.WalkDir(bundled, embedRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return "unreadable"
	}
	sort.Strings(files)
	h := sha256.New()
	for _, path := range files {
		data, err := bundled.ReadFile(path)
		if err != nil {
			return "unreadable"
		}
		_, _ = fmt.Fprintf(h, "%s %d\n", path, len(data))
		_, _ = h.Write(data)
	}
	return hex.EncodeToString(h.Sum(nil))
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
