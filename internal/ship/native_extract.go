package ship

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractTarGz(archivePath, destRoot string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	var prefix string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := writeTarEntry(destRoot, &prefix, hdr, tr); err != nil {
			return err
		}
	}
}

func writeTarEntry(destRoot string, prefix *string, hdr *tar.Header, tr *tar.Reader) error {
	if strings.HasPrefix(hdr.Name, "/") || filepath.IsAbs(hdr.Name) {
		return fmt.Errorf("refusing absolute path %q", hdr.Name)
	}
	if appleDouble(hdr.Name) {
		return nil
	}
	if *prefix == "" {
		if parts := strings.SplitN(hdr.Name, "/", 2); len(parts) > 1 {
			*prefix = parts[0] + "/"
		}
	}
	rel := strings.TrimPrefix(hdr.Name, *prefix)
	if rel == "" || rel == "." || appleDouble(rel) {
		return nil
	}
	target, err := safeJoin(destRoot, rel)
	if err != nil {
		return err
	}
	switch hdr.Typeflag {
	case tar.TypeDir:
		return os.MkdirAll(target, 0o755)
	case tar.TypeReg, tar.TypeRegA:
		return writeTarFile(target, hdr.Mode, tr)
	case tar.TypeSymlink:
		return writeTarSymlink(target, hdr.Linkname)
	default:
		return nil
	}
}

func appleDouble(name string) bool {
	base := filepath.Base(name)
	return strings.HasPrefix(base, "._") || base == ".DS_Store"
}

func writeTarSymlink(target, link string) error {
	if filepath.IsAbs(link) {
		return fmt.Errorf("refusing absolute symlink %q", link)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	_ = os.Remove(target)
	return os.Symlink(link, target)
}

func fixDylibNames(binDir string) {
	var real string
	for _, c := range []string{"libqwen3tts.0.1.0.dylib", "libqwen3tts.0.dylib", "libqwen3tts.dylib"} {
		path := filepath.Join(binDir, c)
		info, err := os.Stat(path)
		if err == nil && info.Mode().IsRegular() {
			real = c
			break
		}
	}
	if real == "" {
		return
	}
	for _, name := range []string{"libqwen3tts.0.dylib", "libqwen3tts.dylib"} {
		dst := filepath.Join(binDir, name)
		if _, err := os.Lstat(dst); err == nil {
			continue
		}
		_ = os.Symlink(real, dst)
	}
}

func writeTarFile(target string, mode int64, r io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	out, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(mode)|0o644)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, r)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}

func safeJoin(dir, rel string) (string, error) {
	if filepath.IsAbs(rel) {
		return "", fmt.Errorf("unsafe absolute path %q", rel)
	}
	target := filepath.Join(dir, rel)
	within, err := filepath.Rel(dir, target)
	if err != nil || within == ".." || strings.HasPrefix(within, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("unsafe path %q", rel)
	}
	return target, nil
}
