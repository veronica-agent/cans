package say

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestResolveTextErrors(t *testing.T) {
	tests := []struct {
		name  string
		opts  func(*Options)
		stdin string
		tty   bool
		want  string
	}{
		{"empty argv on a terminal", func(o *Options) { o.StdinTTY = true }, "", true, "say: missing text\n"},
		{"empty pipe", nil, "", false, "say: empty text\n"},
		{"whitespace pipe", nil, "  \n\t \n", false, "say: empty text\n"},
		{"bare dash on an empty pipe", func(o *Options) { o.Stdin = true }, "", false, "say: empty text\n"},
		{"bare dash on a terminal reads it", func(o *Options) { o.Stdin, o.StdinTTY = true, true }, "", true, "say: empty text\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := DefaultOptions()
			if tt.opts != nil {
				tt.opts(&o)
			}
			var errBuf bytes.Buffer
			var in io.Reader = strings.NewReader(tt.stdin)
			if tt.tty && !o.Stdin {
				in = failingReader{t}
			}
			got, code := resolveText(o, in, &errBuf)
			if code != ExitUsage {
				t.Fatalf("code %d, want %d", code, ExitUsage)
			}
			if got != "" {
				t.Fatalf("text %q, want empty", got)
			}
			if errBuf.String() != tt.want {
				t.Fatalf("stderr %q, want %q", errBuf.String(), tt.want)
			}
		})
	}
}

func TestResolveText(t *testing.T) {
	tests := []struct {
		name  string
		opts  func(*Options)
		stdin string
		want  string
	}{
		{"argv text wins", func(o *Options) { o.Text = "Put the cans on." }, "", "Put the cans on."},
		{"pipe is one utterance", nil, "Put the cans on.\n", "Put the cans on."},
		{"bare dash reads the pipe", func(o *Options) { o.Stdin = true }, "Put the cans on.\n", "Put the cans on."},
		{"whole pipe, newlines and all", nil, "line one\nline two\n", "line one\nline two"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := DefaultOptions()
			if tt.opts != nil {
				tt.opts(&o)
			}
			var errBuf bytes.Buffer
			got, code := resolveText(o, strings.NewReader(tt.stdin), &errBuf)
			if code != ExitOK {
				t.Fatalf("code %d stderr %q", code, errBuf.String())
			}
			if got != tt.want {
				t.Fatalf("text %q, want %q", got, tt.want)
			}
			if errBuf.Len() != 0 {
				t.Fatalf("stderr %q", errBuf.String())
			}
		})
	}
}

func TestResolveTextDoesNotReadStdinForArgvText(t *testing.T) {
	o := DefaultOptions()
	o.Text = "Put the cans on."
	var errBuf bytes.Buffer
	got, code := resolveText(o, failingReader{t}, &errBuf)
	if code != ExitOK || got != o.Text {
		t.Fatalf("%q %d", got, code)
	}
}

func TestResolveTextStdinTooLarge(t *testing.T) {
	o := DefaultOptions()
	var errBuf bytes.Buffer
	got, code := resolveText(o, io.LimitReader(fillReader{}, stdinLimit+1), &errBuf)
	if code != ExitFail {
		t.Fatalf("code %d, want %d", code, ExitFail)
	}
	if got != "" {
		t.Fatalf("text %q, want empty", got)
	}
	if errBuf.String() != "say: stdin too large\n" {
		t.Fatalf("stderr %q", errBuf.String())
	}
}

func TestResolveTextStdinAtCap(t *testing.T) {
	o := DefaultOptions()
	var errBuf bytes.Buffer
	got, code := resolveText(o, io.LimitReader(fillReader{}, stdinLimit), &errBuf)
	if code != ExitOK {
		t.Fatalf("code %d stderr %q", code, errBuf.String())
	}
	if len(got) != stdinLimit {
		t.Fatalf("len %d, want %d", len(got), stdinLimit)
	}
}

// fillReader fills every Read with 'x' and never EOFs.
type fillReader struct{}

func (fillReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 'x'
	}
	return len(p), nil
}
