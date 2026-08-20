package ship

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func downloadArchive(ctx context.Context, root, url, wantSHA string) (string, error) {
	tmp, err := os.CreateTemp(root, ".cans-native-*.tar.gz.part")
	if err != nil {
		return "", err
	}
	name := tmp.Name()
	ok := false
	defer func() {
		_ = tmp.Close()
		if !ok {
			_ = os.Remove(name)
		}
	}()
	body, err := openSource(ctx, url)
	if err != nil {
		return "", err
	}
	defer body.Close()
	h := sha256.New()
	w := io.MultiWriter(tmp, h)
	buf := make([]byte, 256<<10)
	for {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		n, readErr := body.Read(buf)
		if n > 0 {
			if _, err := w.Write(buf[:n]); err != nil {
				return "", err
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return "", readErr
		}
	}
	sum := hex.EncodeToString(h.Sum(nil))
	if wantSHA != "" && !strings.EqualFold(wantSHA, sum) {
		return "", fmt.Errorf("sha256 mismatch (got %s)", sum)
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	ok = true
	return name, nil
}

func openSource(ctx context.Context, url string) (io.ReadCloser, error) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("download HTTP %d", resp.StatusCode)
		}
		return resp.Body, nil
	}
	path := strings.TrimPrefix(url, "file://")
	return os.Open(path)
}
