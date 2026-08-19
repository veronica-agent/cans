package keep

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Current is the pinned throat.
type Current struct {
	Wav     string `json:"wav"`
	RefText string `json:"ref_text"`
}

func homeDir() string {
	if h := os.Getenv("CANS_HOME"); h != "" {
		return h
	}
	u, err := os.UserHomeDir()
	if err != nil {
		return ".cans"
	}
	return filepath.Join(u, ".cans")
}

func rootDir() string {
	if r := os.Getenv("CANS_ROOT"); r != "" {
		return r
	}
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}

func currentPath() string {
	return filepath.Join(homeDir(), "current.json")
}

// Default is shipped Veronica.
func Default() Current {
	return Current{
		Wav:     filepath.Join(rootDir(), "voices", "veronica", "ref.wav"),
		RefText: "Just like that, feel the rhythm of my voice.",
	}
}

// Load returns the kept throat, or Default if none.
func Load() (Current, error) {
	p := currentPath()
	f, err := os.Open(p)
	if err != nil {
		if os.IsNotExist(err) {
			d := Default()
			if _, statErr := os.Stat(d.Wav); statErr != nil {
				return Current{}, fmt.Errorf("default ref missing: %s", d.Wav)
			}
			return d, nil
		}
		return Current{}, err
	}
	defer f.Close()
	var c Current
	if err := json.NewDecoder(f).Decode(&c); err != nil && err != io.EOF {
		return Current{}, err
	}
	if c.Wav == "" {
		return Default(), nil
	}
	if _, err := os.Stat(c.Wav); err != nil {
		return Current{}, fmt.Errorf("kept wav missing: %s", c.Wav)
	}
	return c, nil
}

// Pin copies wav into CANS_HOME and writes current.json.
func Pin(wavPath, refText string) (Current, error) {
	src, err := filepath.Abs(wavPath)
	if err != nil {
		return Current{}, err
	}
	if _, err := os.Stat(src); err != nil {
		return Current{}, fmt.Errorf("keep: wav not found: %s", wavPath)
	}
	dir := filepath.Join(homeDir(), "current")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Current{}, err
	}
	dst := filepath.Join(dir, "ref.wav")
	in, err := os.Open(src)
	if err != nil {
		return Current{}, err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return Current{}, err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return Current{}, err
	}
	if err := out.Close(); err != nil {
		return Current{}, err
	}
	c := Current{Wav: dst, RefText: refText}
	body, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return Current{}, err
	}
	if err := os.MkdirAll(homeDir(), 0o755); err != nil {
		return Current{}, err
	}
	if err := os.WriteFile(currentPath(), body, 0o644); err != nil {
		return Current{}, err
	}
	return c, nil
}
