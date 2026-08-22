package say

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/veronica-agent/cans/internal/doctor"
	"github.com/veronica-agent/cans/internal/keep"
	"github.com/veronica-agent/cans/internal/tts"
)

var errLineFailed = errors.New("line failed")

type okRecord struct {
	Line       int    `json:"line"`
	Wav        string `json:"wav,omitempty"`
	TTFAMs     int    `json:"ttfa_ms"`
	SampleRate int    `json:"sample_rate"`
}

type errRecord struct {
	Line  int    `json:"line"`
	Error string `json:"error"`
}

// streamer is the fixed half of a stream: the worker, the throat and the
// writers, all of which outlive every line.
type streamer struct {
	sess   *tts.Session
	cur    keep.Current
	o      Options
	stdout *bufio.Writer
	stderr io.Writer
}

func runStream(ctx context.Context, o Options, stdin io.Reader, stdout, stderr io.Writer) int {
	if err := doctor.Prepare(ctx, stderr); err != nil {
		fmt.Fprintln(stderr, err)
		return ExitFail
	}
	cur, err := keep.Load()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return ExitFail
	}
	sess, err := tts.OpenWith(ctx, lockOpts(o, stderr))
	if err != nil {
		return exitFor(err, stderr)
	}
	defer sess.Close()
	outw := bufio.NewWriter(stdout)
	defer outw.Flush()
	s := &streamer{sess: sess, cur: cur, o: o, stdout: outw, stderr: stderr}
	return s.scan(ctx, stdin)
}

// scan speaks stdin one line at a time. last is the last line fully processed,
// spoken or reported as failed; idx counts only the lines that were spoken, so
// an -o template stays dense (D005).
func (s *streamer) scan(ctx context.Context, stdin io.Reader) int {
	sc := bufio.NewScanner(stdin)
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	lineNo, idx, failed, last := 0, 0, 0, 0
	for {
		if ctx.Err() != nil {
			return s.interruptedAfter(last)
		}
		if !sc.Scan() {
			break
		}
		if ctx.Err() != nil {
			return s.interruptedAfter(last)
		}
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		err := s.speak(ctx, line, lineNo, idx+1)
		if interrupted(err) {
			return s.interruptedAfter(last)
		}
		last = lineNo
		if err != nil {
			failed++
			continue
		}
		idx++
	}
	// Ctrl-C while blocked in Scan: the producer got the same signal and
	// closed the pipe, so Scan returned false — still an interrupt, not EOF.
	if ctx.Err() != nil {
		return s.interruptedAfter(last)
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintf(s.stderr, "say: stdin: %v\n", err)
		return ExitFail
	}
	if failed > 0 {
		return ExitFail
	}
	return ExitOK
}

// speak says one stdin line. It returns the cancellation error unwrapped so the
// caller can tell an interrupt from a line that failed.
func (s *streamer) speak(ctx context.Context, line string, lineNo, idx int) error {
	r, err := s.sess.SayTo(ctx, line, s.cur, outPath(s.o.Out, idx))
	if err != nil {
		if interrupted(err) {
			return err
		}
		fmt.Fprintf(s.stderr, "line %d: %v\n", lineNo, err)
		_ = s.emit(errRecord{Line: lineNo, Error: err.Error()})
		return errLineFailed
	}
	rec := okRecord{Line: lineNo, TTFAMs: r.TTFAMs, SampleRate: r.SampleRate}
	if s.o.Out != "" {
		rec.Wav = r.Wav
	}
	if err := s.emit(rec); err != nil {
		fmt.Fprintln(s.stderr, err)
		return errLineFailed
	}
	if err := playTail(s.o, r.Wav); err != nil {
		fmt.Fprintln(s.stderr, err)
		return errLineFailed
	}
	return nil
}

// interruptedAfter reports the Ctrl-C (D014): the in-flight line is dropped, so
// the last line named is the last one fully processed.
func (s *streamer) interruptedAfter(line int) int {
	if line == 0 {
		fmt.Fprintln(s.stderr, "interrupted before the first line")
	} else {
		fmt.Fprintf(s.stderr, "interrupted after line %d\n", line)
	}
	return ExitInterrupted
}

// emit writes the one record for a line and flushes, so a reader downstream
// sees it before the next line is spoken.
func (s *streamer) emit(rec any) error {
	var err error
	switch r := rec.(type) {
	case errRecord:
		if s.o.JSON {
			err = json.NewEncoder(s.stdout).Encode(r)
		}
	case okRecord:
		switch {
		case s.o.JSON:
			err = json.NewEncoder(s.stdout).Encode(r)
		case s.o.Out != "":
			_, err = fmt.Fprintln(s.stdout, r.Wav)
		default:
			_, err = fmt.Fprintf(s.stdout, "ttfa_ms=%d\n", r.TTFAMs)
		}
	}
	if err != nil {
		return err
	}
	return s.stdout.Flush()
}
