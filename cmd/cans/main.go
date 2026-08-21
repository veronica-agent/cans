package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/veronica-agent/cans/internal/booth"
	"github.com/veronica-agent/cans/internal/doctor"
	"github.com/veronica-agent/cans/internal/keep"
	"github.com/veronica-agent/cans/internal/say"
	"github.com/veronica-agent/cans/internal/ship"
)

var (
	stdin  io.Reader = os.Stdin
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr
)

const usage = `cans — put the cans on.

Apple Silicon. The mouth is a native Qwen3-TTS worker cloning a wav.

  cans                         booth (throat frozen for the session)
  cans say <text>              speak one line
  cans keep <wav> -text WORDS  freeze this throat (both orders work)
  cans doctor                  set up the mouth, check the machine
  cans version                 print version

exit 75 when another cans holds the mouth and --nowait was set
`

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		return runBooth()
	}
	switch args[0] {
	case "say":
		return runSay(args[1:])
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

// runSay speaks one `cans say`, cancellable by SIGINT or SIGTERM.
func runSay(args []string) int {
	o, err := parseSay(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	o.StdinTTY = stdinIsTTY()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	// D014: the first signal cancels; stop() then restores the default
	// disposition, so a second Ctrl-C ends the process at once. stop() cancels
	// ctx as well, so this goroutine always ends with runSay.
	go func() {
		<-ctx.Done()
		stop()
	}()
	return say.Run(ctx, o, stdin, stdout, stderr)
}

// runBooth prepares the mouth and opens the TUI on the frozen throat.
func runBooth() int {
	if err := doctor.Prepare(context.Background(), stderr); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	throat, err := keep.Load()
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := booth.Run(context.Background(), keep.Quote(), throat); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
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
