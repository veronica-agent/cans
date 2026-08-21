package say

import (
	"bytes"
	"context"
	"testing"
)

func TestCheckOutErrors(t *testing.T) {
	tests := []struct {
		name   string
		out    string
		stream bool
		want   string
	}{
		{"stream %s", "%s", true, "say: -o needs one %d in --stream"},
		{"stream two verbs", "%03d-%d", true, "say: -o needs one %d in --stream"},
		{"stream no verb", "out.wav", true, "say: -o needs one %d in --stream"},
		{"one-shot %03d", "%03d", false, "say: -o template needs --stream"},
		{"one-shot %q", "%q", false, "say: -o template needs --stream"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkOut(tt.out, tt.stream)
			if err == nil || err.Error() != tt.want {
				t.Fatalf("err %v, want %q", err, tt.want)
			}
		})
	}
}

func TestOutPath(t *testing.T) {
	tests := []struct {
		out  string
		idx  int
		want string
	}{
		{"%03d", 1, "001"},
		{"%d", 1, "1"},
		{"out/%02d.wav", 1, "out/01.wav"},
		{"foo%%bar.wav", 0, "foo%bar.wav"},
		{"x%%y%03d", 1, "x%y001"},
	}
	for _, tt := range tests {
		t.Run(tt.out, func(t *testing.T) {
			if got := outPath(tt.out, tt.idx); got != tt.want {
				t.Fatalf("outPath(%q,%d)=%q want %q", tt.out, tt.idx, got, tt.want)
			}
		})
	}
}

func TestStreamOutWithoutVerbIsUsage(t *testing.T) {
	o := DefaultOptions()
	o.Stream = true
	o.Out = "out.wav"
	var out, errBuf bytes.Buffer
	code := Run(context.Background(), o, failingReader{t}, &out, &errBuf)
	if code != ExitUsage {
		t.Fatalf("code %d", code)
	}
	if errBuf.String() != "say: -o needs one %d in --stream\n" {
		t.Fatalf("stderr %q", errBuf.String())
	}
}
