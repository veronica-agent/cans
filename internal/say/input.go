package say

import (
	"fmt"
	"io"
	"strings"
)

// stdinLimit caps one stdin utterance. Over the cap is an error, not a clip.
const stdinLimit = 4 << 20

// resolveText picks the utterance: the argv text, or the whole of stdin when a
// bare `-` or an empty argv asked for it. The int is an exit code, 0 on
// success; the reason for anything else is already on stderr.
func resolveText(o Options, stdin io.Reader, stderr io.Writer) (string, int) {
	if !o.Stdin && o.Text != "" {
		return o.Text, ExitOK
	}
	if o.StdinTTY && !o.Stdin {
		fmt.Fprintln(stderr, "say: missing text")
		return "", ExitUsage
	}
	if stdin == nil {
		fmt.Fprintln(stderr, "say: empty text")
		return "", ExitUsage
	}
	b, err := io.ReadAll(io.LimitReader(stdin, stdinLimit+1))
	if err != nil {
		fmt.Fprintf(stderr, "say: stdin: %v\n", err)
		return "", ExitFail
	}
	if len(b) > stdinLimit {
		fmt.Fprintln(stderr, "say: stdin too large")
		return "", ExitFail
	}
	text := strings.TrimSpace(string(b))
	if text == "" {
		fmt.Fprintln(stderr, "say: empty text")
		return "", ExitUsage
	}
	return text, ExitOK
}
