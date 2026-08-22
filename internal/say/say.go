package say

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/veronica-agent/cans/internal/doctor"
	"github.com/veronica-agent/cans/internal/keep"
	"github.com/veronica-agent/cans/internal/mouth"
	"github.com/veronica-agent/cans/internal/play"
	"github.com/veronica-agent/cans/internal/tts"
)

// Run speaks one `cans say` and returns the process exit code.
// stdout carries data only: ttfa_ms=, a wav path, or a JSON record.
func Run(ctx context.Context, o Options, stdin io.Reader, stdout, stderr io.Writer) int {
	if err := checkOut(o.Out, o.Stream); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitUsage
	}
	if o.Stream {
		return runStream(ctx, o, stdin, stdout, stderr)
	}
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
	r, err := tts.SayToWith(ctx, text, cur, outPath(o.Out, 0), lockOpts(o, stderr))
	if err != nil {
		return exitFor(err, stderr)
	}
	if err := emit(stdout, o, r); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitFail
	}
	if err := playTail(o, r.Wav); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitFail
	}
	return ExitOK
}

// playTail is the tail every spoken utterance shares. With no -o the wav is a
// temp file: it is played, then removed whether or not playing worked. With -o
// the wav is the caller's and is only played on --play.
func playTail(o Options, wav string) error {
	if o.Out == "" {
		err := play.File(wav)
		tts.RemoveTemp(wav)
		return err
	}
	if o.Play {
		return play.File(wav)
	}
	return nil
}

// jsonShot is a one-shot --json record. wav is omitted when there is no -o:
// that file is a temp that playTail deletes, and a script must not trust it.
type jsonShot struct {
	Wav        string `json:"wav,omitempty"`
	TTFAMs     int    `json:"ttfa_ms"`
	SampleRate int    `json:"sample_rate"`
}

// emit writes the one record for an utterance: a JSON line, the wav path, or
// the v1 ttfa_ms line. stdout carries nothing else.
func emit(stdout io.Writer, o Options, r tts.Result) error {
	switch {
	case o.JSON:
		rec := jsonShot{TTFAMs: r.TTFAMs, SampleRate: r.SampleRate}
		if o.Out != "" {
			rec.Wav = r.Wav
		}
		return json.NewEncoder(stdout).Encode(rec)
	case o.Out != "":
		_, err := fmt.Fprintln(stdout, r.Wav)
		return err
	default:
		_, err := fmt.Fprintf(stdout, "ttfa_ms=%d\n", r.TTFAMs)
		return err
	}
}

func lockOpts(o Options, stderr io.Writer) tts.Options {
	return tts.Options{
		Wait:   o.Wait,
		OnWait: func() { fmt.Fprintln(stderr, "waiting for the mouth…") },
	}
}

func exitFor(err error, stderr io.Writer) int {
	if errors.Is(err, mouth.ErrBusy) {
		fmt.Fprintln(stderr, "say: mouth busy")
		return ExitBusy
	}
	if interrupted(err) {
		fmt.Fprintln(stderr, "say: interrupted")
		return ExitInterrupted
	}
	fmt.Fprintln(stderr, err)
	return ExitFail
}

func interrupted(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
