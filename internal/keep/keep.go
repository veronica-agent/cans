package keep

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/veronica-agent/cans/internal/audio"
)

// Current is the pinned throat.
type Current struct {
	Wav     string `json:"wav"`
	RefText string `json:"ref_text"`
}

func homeDir() string {
	if h := os.Getenv("CANS_HOME"); h != "" {
		return filepath.Clean(h)
	}
	u, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(u) == "" {
		return ".cans"
	}
	return filepath.Join(u, ".cans")
}

func currentDir() string {
	return filepath.Join(homeDir(), "current")
}

func currentPath() string {
	return filepath.Join(homeDir(), "current.json")
}

// Root is the repo root that holds voices/veronica.
func Root() string {
	if r := os.Getenv("CANS_ROOT"); r != "" {
		return r
	}
	var starts []string
	if wd, err := os.Getwd(); err == nil {
		starts = append(starts, wd)
	}
	if exe, err := os.Executable(); err == nil {
		starts = append(starts, filepath.Dir(exe))
	}
	for _, start := range starts {
		dir := start
		for i := 0; i < 8; i++ {
			cand := filepath.Join(dir, "voices", "veronica", "ref.wav")
			if audio.HeaderOK(cand) == nil {
				return dir
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "."
}

// Default is shipped Veronica.
func Default() Current {
	return Current{
		Wav:     filepath.Join(Root(), "voices", "veronica", "ref.wav"),
		RefText: "Just like that, feel the rhythm of my voice.",
	}
}

func allowed(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if err := audio.HeaderOK(abs); err != nil {
		return err
	}
	def, err := filepath.Abs(Default().Wav)
	if err == nil && abs == def {
		return nil
	}
	cur, err := filepath.Abs(currentDir())
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(cur, abs)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("keep: wav not in CANS_HOME/current or shipped default")
	}
	return nil
}

// Load returns the kept throat, or Default if none is pinned.
func Load() (Current, error) {
	p := currentPath()
	f, err := os.Open(p)
	if err != nil {
		if os.IsNotExist(err) {
			d := Default()
			if err := audio.HeaderOK(d.Wav); err != nil {
				return Current{}, fmt.Errorf("default ref missing: %s", d.Wav)
			}
			return d, nil
		}
		return Current{}, err
	}
	defer f.Close()
	var c Current
	if err := json.NewDecoder(f).Decode(&c); err != nil && err != io.EOF {
		return Current{}, fmt.Errorf("keep: corrupt current.json: %w", err)
	}
	if strings.TrimSpace(c.Wav) == "" || strings.TrimSpace(c.RefText) == "" {
		return Current{}, fmt.Errorf("keep: corrupt current.json")
	}
	if !filepath.IsAbs(c.Wav) {
		return Current{}, fmt.Errorf("keep: wav path must be absolute")
	}
	if err := allowed(c.Wav); err != nil {
		return Current{}, err
	}
	return c, nil
}

// Pin copies wav into CANS_HOME/current and writes current.json.
func Pin(wavPath, refText string) (Current, error) {
	refText = strings.TrimSpace(refText)
	if refText == "" {
		return Current{}, fmt.Errorf("keep: -text is required (the words spoken in the wav)")
	}
	src, err := filepath.Abs(wavPath)
	if err != nil {
		return Current{}, err
	}
	if err := audio.HeaderOK(src); err != nil {
		return Current{}, fmt.Errorf("keep: %w", err)
	}
	dir := currentDir()
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
	if err := os.WriteFile(currentPath(), body, 0o644); err != nil {
		return Current{}, err
	}
	return c, nil
}

// Quote from character.toml, or the pin line.
func Quote() string {
	p := filepath.Join(Root(), "character.toml")
	b, err := os.ReadFile(p)
	if err != nil {
		return "Put the cans on."
	}
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "quote") {
			_, rest, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			q := strings.Trim(strings.TrimSpace(rest), `"`)
			if q != "" {
				return q
			}
		}
	}
	return "Put the cans on."
}
