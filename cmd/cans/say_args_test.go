package main

import (
	"strings"
	"testing"
	"time"

	"github.com/veronica-agent/cans/internal/say"
)

func TestParseSayErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"unknown flag", []string{"--bogus"}, "say: unknown flag --bogus"},
		{"unknown short flag", []string{"-x", "Put the cans on."}, "say: unknown flag -x"},
		{"help short", []string{"-h"}, "cans say [-o out.wav]"},
		{"help long", []string{"--help"}, "cans say [-o out.wav]"},
		{"wait unparsable", []string{"--wait", "bogus"}, `say: --wait "bogus" is not a duration`},
		{"wait zero", []string{"--wait", "0s"}, "say: --wait must be positive"},
		{"wait negative", []string{"--wait=-2s"}, "say: --wait must be positive"},
		{"nowait and wait", []string{"--nowait", "--wait", "1s"}, "say: --nowait and --wait together"},
		{"play without out", []string{"--play", "Put the cans on."}, "say: --play needs -o"},
		{"stdin with text", []string{"-", "Put the cans on."}, "say: - and text together"},
		{"stream with text", []string{"--stream", "Put the cans on."}, "say: --stream reads stdin; drop the text"},
		{"out without value", []string{"-o"}, "say: -o needs a path"},
		{"wait without value", []string{"--wait"}, "say: --wait needs a duration"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseSay(tt.args)
			if err == nil {
				t.Fatalf("expected an error for %q", tt.args)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error %q does not contain %q", err, tt.want)
			}
		})
	}
}

func TestParseSayGrammar(t *testing.T) {
	base := say.DefaultOptions()
	with := func(f func(*say.Options)) say.Options {
		o := base
		f(&o)
		return o
	}
	tests := []struct {
		name string
		args []string
		want say.Options
	}{
		{"no args", nil, base},
		{"text only", []string{"Put the cans on."},
			with(func(o *say.Options) { o.Text = "Put the cans on." })},
		{"multiple positionals joined", []string{"Put", "the", "cans", "on."},
			with(func(o *say.Options) { o.Text = "Put the cans on." })},
		{"out before text", []string{"-o", "out.wav", "Put the cans on."},
			with(func(o *say.Options) { o.Text, o.Out = "Put the cans on.", "out.wav" })},
		{"out after text", []string{"Put the cans on.", "-o", "out.wav"},
			with(func(o *say.Options) { o.Text, o.Out = "Put the cans on.", "out.wav" })},
		{"out equals", []string{"-o=out.wav", "Put the cans on."},
			with(func(o *say.Options) { o.Text, o.Out = "Put the cans on.", "out.wav" })},
		{"long out", []string{"--out", "out.wav", "Put the cans on."},
			with(func(o *say.Options) { o.Text, o.Out = "Put the cans on.", "out.wav" })},
		{"long out equals", []string{"--out=out.wav", "Put the cans on."},
			with(func(o *say.Options) { o.Text, o.Out = "Put the cans on.", "out.wav" })},
		{"json and stream", []string{"--json", "--stream"},
			with(func(o *say.Options) { o.JSON, o.Stream = true, true })},
		{"stdin alone", []string{"-"},
			with(func(o *say.Options) { o.Stdin = true })},
		{"play with out", []string{"--play", "-o", "out.wav", "Put the cans on."},
			with(func(o *say.Options) { o.Text, o.Out, o.Play = "Put the cans on.", "out.wav", true })},
		{"wait duration", []string{"--wait", "30s", "Put the cans on."},
			with(func(o *say.Options) { o.Text, o.Wait = "Put the cans on.", 30*time.Second })},
		{"wait equals", []string{"--wait=30s"},
			with(func(o *say.Options) { o.Wait = 30 * time.Second })},
		{"nowait", []string{"--nowait", "Put the cans on."},
			with(func(o *say.Options) { o.Text, o.Wait = "Put the cans on.", 0 })},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSay(tt.args)
			if err != nil {
				t.Fatalf("parseSay(%q): %v", tt.args, err)
			}
			if got != tt.want {
				t.Fatalf("parseSay(%q) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}
}

func TestParseSayBothOrdersMatch(t *testing.T) {
	first, err := parseSay([]string{"Put the cans on.", "-o", "out.wav"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := parseSay([]string{"-o", "out.wav", "Put the cans on."})
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatalf("%+v != %+v", first, second)
	}
}

func TestParseSayDefaultWaitsForever(t *testing.T) {
	o, err := parseSay([]string{"Put the cans on."})
	if err != nil {
		t.Fatal(err)
	}
	if o.Wait != say.WaitForever {
		t.Fatalf("wait %v, want %v", o.Wait, say.WaitForever)
	}
}
