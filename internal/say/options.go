// Package say runs cans say: one-shot, file output, stdin, stream.
package say

import "time"

// WaitForever is the default Options.Wait: block until the mouth is free.
const WaitForever time.Duration = -1

// Options is one parsed `cans say` command line.
type Options struct {
	// Text is the positional arguments joined with single spaces.
	Text string
	// Stdin is a bare `-`: read the utterance from stdin.
	Stdin bool
	// StdinTTY is set by main, never by the parser.
	StdinTTY bool
	// Out is the -o/--out path; empty means a temp wav.
	Out    string
	JSON   bool
	Stream bool
	Play   bool
	// Wait bounds the mouth lock: WaitForever, 0 (--nowait), or --wait d.
	Wait time.Duration
}

// DefaultOptions is a say with no flags: wait forever for the mouth.
func DefaultOptions() Options {
	return Options{Wait: WaitForever}
}
