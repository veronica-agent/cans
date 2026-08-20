package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/veronica-agent/cans/internal/booth"
	"github.com/veronica-agent/cans/internal/doctor"
	"github.com/veronica-agent/cans/internal/keep"
	"github.com/veronica-agent/cans/internal/play"
	"github.com/veronica-agent/cans/internal/ship"
	"github.com/veronica-agent/cans/internal/tts"
)

var (
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr
)

const usage = `cans — put the cans on.

Apple Silicon + MLX. The mouth is Qwen3-TTS 0.6B Base cloning a wav.

  cans                         booth (throat frozen for the session)
  cans say <text>              speak one line
  cans keep <wav> -text WORDS  freeze this throat (both orders work)
  cans doctor                  install sidecar, check the machine
  cans version                 print version
`

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		if err := doctor.Prepare(context.Background(), stderr); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		throat, err := keep.Load()
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		if err := booth.Run(keep.Quote(), throat); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		return 0
	}
	switch args[0] {
	case "say":
		text := strings.TrimSpace(strings.Join(args[1:], " "))
		if text == "" {
			fmt.Fprintln(stderr, "say: missing text")
			return 2
		}
		if err := doctor.Prepare(context.Background(), stderr); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		r, err := tts.Say(text)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintf(stdout, "ttfa_ms=%d\n", r.TTFAMs)
		playErr := play.File(r.Wav)
		tts.RemoveTemp(r.Wav)
		if playErr != nil {
			fmt.Fprintln(stderr, playErr)
			return 1
		}
		return 0
	case "doctor":
		if err := doctor.Run(context.Background(), stdout, stderr); err != nil {
			return 1
		}
		return 0
	case "version":
		fmt.Fprintln(stdout, "cans "+ship.Version)
		return 0
	case "keep":
		wav, text, err := parseKeep(args[1:])
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		c, err := keep.Pin(wav, text)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		fmt.Fprintf(stdout, "kept %s\n", c.Wav)
		return 0
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usage)
		return 0
	default:
		fmt.Fprint(stderr, usage)
		return 2
	}
}

// parseKeep accepts both `keep take.wav -text words` and `keep -text words take.wav`.
func parseKeep(args []string) (wav, text string, err error) {
	var positional []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-text" || a == "--text":
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("keep: -text needs a value")
			}
			i++
			text = args[i]
		case strings.HasPrefix(a, "-text="):
			text = strings.TrimPrefix(a, "-text=")
		case strings.HasPrefix(a, "--text="):
			text = strings.TrimPrefix(a, "--text=")
		case a == "-h" || a == "--help":
			return "", "", fmt.Errorf("%s", strings.TrimSpace(usage))
		case strings.HasPrefix(a, "-"):
			return "", "", fmt.Errorf("keep: unknown flag %s", a)
		default:
			positional = append(positional, a)
		}
	}
	if len(positional) != 1 {
		return "", "", fmt.Errorf("keep: need a wav path")
	}
	if strings.TrimSpace(text) == "" {
		return "", "", fmt.Errorf("keep: -text is required (the words spoken in the wav)")
	}
	return positional[0], text, nil
}
