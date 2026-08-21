package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/veronica-agent/cans/internal/say"
)

// sayValueFlags are the say flags that take a value, so `--flag=v` splits.
var sayValueFlags = map[string]bool{"-o": true, "--out": true, "--wait": true}

// parseSay accepts flags and text in either order, the way parseKeep does:
// `say "line" -o out.wav` and `say -o out.wav "line"` parse identically.
func parseSay(args []string) (say.Options, error) {
	o := say.DefaultOptions()
	var positional []string
	var sawWait, sawNoWait bool
	args = splitSayEquals(args)
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch a {
		case "-o", "--out":
			i++
			if i >= len(args) {
				return say.Options{}, fmt.Errorf("say: %s needs a path", a)
			}
			o.Out = args[i]
		case "--wait":
			i++
			if i >= len(args) {
				return say.Options{}, fmt.Errorf("say: --wait needs a duration")
			}
			d, err := parseWait(args[i])
			if err != nil {
				return say.Options{}, err
			}
			o.Wait, sawWait = d, true
		case "--json":
			o.JSON = true
		case "--stream":
			o.Stream = true
		case "--play":
			o.Play = true
		case "--nowait":
			o.Wait, sawNoWait = 0, true
		case "-":
			o.Stdin = true
		case "-h", "--help":
			return say.Options{}, fmt.Errorf("%s", strings.TrimSpace(usage))
		default:
			if strings.HasPrefix(a, "-") {
				return say.Options{}, fmt.Errorf("say: unknown flag %s", a)
			}
			positional = append(positional, a)
		}
	}
	o.Text = strings.TrimSpace(strings.Join(positional, " "))
	return o, validateSay(o, sawWait, sawNoWait)
}

// splitSayEquals rewrites `--flag=value` into `--flag value` for value flags.
func splitSayEquals(args []string) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		if name, val, ok := strings.Cut(a, "="); ok && sayValueFlags[name] {
			out = append(out, name, val)
			continue
		}
		out = append(out, a)
	}
	return out
}

func parseWait(v string) (time.Duration, error) {
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("say: --wait %q is not a duration", v)
	}
	if d <= 0 {
		return 0, fmt.Errorf("say: --wait must be positive (use --nowait for none)")
	}
	return d, nil
}

func validateSay(o say.Options, sawWait, sawNoWait bool) error {
	if sawWait && sawNoWait {
		return fmt.Errorf("say: --nowait and --wait together")
	}
	if o.Play && o.Out == "" {
		return fmt.Errorf("say: --play needs -o")
	}
	if o.Stdin && o.Text != "" {
		return fmt.Errorf("say: - and text together")
	}
	return nil
}
