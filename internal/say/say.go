package say

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/veronica-agent/cans/internal/doctor"
	"github.com/veronica-agent/cans/internal/keep"
	"github.com/veronica-agent/cans/internal/play"
	"github.com/veronica-agent/cans/internal/tts"
)

// Run speaks one `cans say` and returns the process exit code.
// stdout carries data only: ttfa_ms=, a wav path, or a JSON record.
func Run(ctx context.Context, o Options, stdin io.Reader, stdout, stderr io.Writer) int {
	// Usage errors must not fetch the mouth.
	text, code := resolveText(o, stdin, stderr)
	if code != ExitOK {
		return code
	}
	if err := doctor.Prepare(ctx, stderr); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitFail
	}
	return runOnce(ctx, o, text, stdout, stderr)
}

// runOnce speaks one utterance: to o.Out when named, else a temp wav that is
// played then deleted.
func runOnce(ctx context.Context, o Options, text string, stdout, stderr io.Writer) int {
	cur, err := keep.Load()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return ExitFail
	}
	r, err := tts.SayTo(ctx, text, cur, o.Out)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return ExitFail
	}
	if err := emit(stdout, o, r); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitFail
	}
	if o.Out == "" {
		playErr := play.File(r.Wav)
		tts.RemoveTemp(r.Wav)
		if playErr != nil {
			fmt.Fprintln(stderr, playErr)
			return ExitFail
		}
		return ExitOK
	}
	if o.Play {
		if err := play.File(r.Wav); err != nil {
			fmt.Fprintln(stderr, err)
			return ExitFail
		}
	}
	return ExitOK
}

// emit writes the one record for an utterance: a JSON line, the wav path, or
// the v1 ttfa_ms line. stdout carries nothing else.
func emit(stdout io.Writer, o Options, r tts.Result) error {
	switch {
	case o.JSON:
		return json.NewEncoder(stdout).Encode(r)
	case o.Out != "":
		_, err := fmt.Fprintln(stdout, r.Wav)
		return err
	default:
		_, err := fmt.Fprintf(stdout, "ttfa_ms=%d\n", r.TTFAMs)
		return err
	}
}
