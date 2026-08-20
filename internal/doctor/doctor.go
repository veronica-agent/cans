package doctor

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/tabwriter"

	"github.com/veronica-agent/cans/internal/audio"
	"github.com/veronica-agent/cans/internal/keep"
	"github.com/veronica-agent/cans/internal/ship"
)

// Check is one doctor line.
type Check struct {
	Name   string
	Detail string
	Hint   string
	OK     bool
}

// Run extracts the payload, syncs the sidecar, and prints a report.
func Run(ctx context.Context, stdout, stderr io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := ship.Ensure(); err != nil {
		fmt.Fprintln(stderr, err)
		return err
	}
	checks := Diagnose(ctx, true)
	printReport(stdout, checks)
	for _, c := range checks {
		if !c.OK {
			return fmt.Errorf("doctor: %s failed", c.Name)
		}
	}
	fmt.Fprintln(stdout, "put the cans on.")
	return nil
}

// Diagnose inspects the machine. mlx is checked only when the venv exists.
func Diagnose(ctx context.Context, mlx bool) []Check {
	_ = mlx
	return []Check{machineCheck(), workerCheck(), payloadCheck(), throatCheck(), playCheck()}
}

func machineCheck() Check {
	c := Check{Name: "machine", Detail: runtime.GOOS + "/" + runtime.GOARCH}
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		c.OK = true
		return c
	}
	c.Hint = "Apple Silicon macOS only"
	return c
}

func workerCheck() Check {
	p := ship.WorkerBin()
	if _, err := os.Stat(p); err != nil {
		return Check{Name: "worker", Hint: "native mouth missing — unpack qwen3-tts-native under ~/.cans/native"}
	}
	return Check{Name: "worker", Detail: p, OK: true}
}

func payloadCheck() Check {
	root := ship.Root()
	if ship.Complete(root) {
		return Check{Name: "payload", Detail: root, OK: true}
	}
	if ship.Complete(ship.Shipped()) {
		return Check{Name: "payload", Detail: ship.Shipped(), OK: true}
	}
	return Check{Name: "payload", Hint: "run cans doctor"}
}

func throatCheck() Check {
	c := Check{Name: "throat"}
	cur, err := keep.Load()
	if err != nil {
		c.Hint = err.Error()
		return c
	}
	if err := audio.HeaderOK(cur.Wav); err != nil {
		c.Hint = err.Error()
		return c
	}
	c.Detail = cur.Wav
	c.OK = true
	return c
}

func playCheck() Check {
	c := Check{Name: "play"}
	bin := "afplay"
	if runtime.GOOS != "darwin" {
		bin = "ffplay"
	}
	p, err := exec.LookPath(bin)
	if err != nil {
		c.Hint = bin + " not on PATH"
		return c
	}
	c.Detail = p
	c.OK = true
	return c
}

func printReport(w io.Writer, checks []Check) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, c := range checks {
		status := "FAIL"
		if c.OK {
			status = "ok"
		}
		detail := c.Detail
		if detail == "" {
			detail = c.Hint
		}
		fmt.Fprintf(tw, "  %s\t%s\t%s\n", c.Name, status, strings.TrimSpace(detail))
		if !c.OK && c.Hint != "" && c.Detail != "" {
			fmt.Fprintf(tw, "  \t\t%s\n", c.Hint)
		}
	}
	_ = tw.Flush()
}

// Prepare extracts the character payload and checks the native worker.
func Prepare(ctx context.Context, stderr io.Writer) error {
	if os.Getenv("CANS_SAY_BIN") != "" {
		return nil
	}
	if err := ship.Ensure(); err != nil {
		return err
	}
	if _, err := os.Stat(ship.WorkerBin()); err != nil {
		return fmt.Errorf("native mouth missing — cans doctor")
	}
	return nil
}
