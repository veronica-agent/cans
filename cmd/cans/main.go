package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/veronica-agent/cans/internal/booth"
	"github.com/veronica-agent/cans/internal/keep"
	"github.com/veronica-agent/cans/internal/play"
	"github.com/veronica-agent/cans/internal/tts"
)

const usage = `cans — put the cans on.

  cans              booth
  cans say <text>   speak one line
  cans keep <wav>   freeze this throat
`

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		if err := booth.Run("Put the cans on."); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	}
	switch args[0] {
	case "say":
		text := strings.TrimSpace(strings.Join(args[1:], " "))
		if text == "" {
			fmt.Fprintln(os.Stderr, "say: missing text")
			return 2
		}
		r, err := tts.Say(text)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		fmt.Printf("ttfa_ms=%d\n", r.TTFAMs)
		if err := play.File(r.Wav); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		return 0
	case "keep":
		fs := flag.NewFlagSet("keep", flag.ContinueOnError)
		refText := fs.String("text", "", "transcript of the keep clip")
		fs.SetOutput(os.Stderr)
		if err := fs.Parse(args[1:]); err != nil {
			return 2
		}
		if fs.NArg() != 1 {
			fmt.Fprintln(os.Stderr, "keep: need a wav path")
			return 2
		}
		c, err := keep.Pin(fs.Arg(0), *refText)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
		fmt.Printf("kept %s\n", c.Wav)
		return 0
	case "-h", "--help", "help":
		fmt.Fprint(os.Stdout, usage)
		return 0
	default:
		fmt.Fprint(os.Stderr, usage)
		return 2
	}
}
