package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestStdinDevNullIsNotTTY(t *testing.T) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	old := stdin
	stdin = f
	defer func() { stdin = old }()
	if stdinIsTTY() {
		t.Fatal("/dev/null must not count as a TTY")
	}
}

func TestSayDevNullIsEmptyText(t *testing.T) {
	f, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var buf bytes.Buffer
	oldErr, oldIn := stderr, stdin
	stderr, stdin = &buf, f
	defer func() { stderr, stdin = oldErr, oldIn }()
	if code := run([]string{"say"}); code != 2 {
		t.Fatalf("code %d", code)
	}
	if !strings.Contains(buf.String(), "say: empty text") {
		t.Fatalf("%q", buf.String())
	}
}
