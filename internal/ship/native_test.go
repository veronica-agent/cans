package ship

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestNativeReadyFalse(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	t.Setenv("CANS_WORKER_BIN", "")
	if NativeReady() {
		t.Fatal("expected missing worker")
	}
}

func TestEnsureNativeFromFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_WORKER_BIN", "")
	archive, sum := fakeNativeTar(t)
	t.Setenv("CANS_NATIVE_URL", archive)
	t.Setenv("CANS_NATIVE_SHA256", sum)
	if err := EnsureNative(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if !NativeReady() {
		t.Fatal("worker missing after ensure")
	}
	if err := EnsureNative(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
}

func TestEnsureNativeSHAMismatch(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CANS_HOME", home)
	t.Setenv("CANS_WORKER_BIN", "")
	archive, _ := fakeNativeTar(t)
	t.Setenv("CANS_NATIVE_URL", archive)
	t.Setenv("CANS_NATIVE_SHA256", "deadbeef")
	if err := EnsureNative(context.Background(), nil); err == nil {
		t.Fatal("expected sha mismatch")
	}
}

func TestEnsureNativeCanceled(t *testing.T) {
	t.Setenv("CANS_HOME", t.TempDir())
	t.Setenv("CANS_WORKER_BIN", "")
	archive, sum := fakeNativeTar(t)
	t.Setenv("CANS_NATIVE_URL", archive)
	t.Setenv("CANS_NATIVE_SHA256", sum)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := EnsureNative(ctx, nil); err == nil {
		t.Fatal("expected cancel")
	}
}

func fakeNativeTar(t *testing.T) (string, string) {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	files := map[string][]byte{
		"pkg/bin/qwen3-tts-worker":           []byte("#!/bin/sh\n"),
		"pkg/models/qwen3-tts-0.6b-f16.gguf": []byte("gguf"),
	}
	for name, body := range files {
		hdr := &tar.Header{Name: name, Mode: 0o755, Size: int64(len(body))}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(body); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "native.tar.gz")
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(buf.Bytes())
	return path, hex.EncodeToString(sum[:])
}
