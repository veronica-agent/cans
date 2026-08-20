package ship

import (
	"os"
	"path/filepath"
)

// Root is the payload directory: voices/, sidecar/, pyproject.toml.
//
// Order: CANS_ROOT, a checkout or brew share/cans next to the binary,
// then ~/.cans/shipped, then cwd.
func Root() string {
	if r := os.Getenv("CANS_ROOT"); r != "" {
		return filepath.Clean(r)
	}
	for _, start := range searchStarts() {
		dir := start
		for i := 0; i < 8; i++ {
			if Complete(dir) {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	if Complete(Shipped()) {
		return Shipped()
	}
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "."
}

func searchStarts() []string {
	var starts []string
	if wd, err := os.Getwd(); err == nil {
		starts = append(starts, wd)
	}
	exe, err := os.Executable()
	if err != nil {
		return starts
	}
	dir := filepath.Dir(exe)
	starts = append(starts, dir)
	share := filepath.Clean(filepath.Join(dir, "..", "share", "cans"))
	starts = append(starts, share)
	return starts
}
