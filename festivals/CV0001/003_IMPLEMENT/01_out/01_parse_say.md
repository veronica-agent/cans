---
fest_type: task
fest_id: 01_parse_say.md
fest_name: parse_say
fest_parent: 01_out
fest_order: 1
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:56.481543-06:00
fest_updated: 2026-08-21T05:19:29.19295-06:00
fest_tracking: true
---


# Task: parse_say

## Objective

`parseSay(args []string) (say.Options, error)` in `cmd/cans/say_args.go` — the whole `say` grammar, settled once, with the `Options` struct it fills in `internal/say/options.go`.

## Requirements

- [x] Flags and text interleave, both orders (D004): `cans say "$line" -o out.wav` and `cans say -o out.wav "$line"` parse identically.
- [x] Flags, and only these: `-o <path>` / `--out <path>` / `-o=<path>` / `--out=<path>`; `--json`; `--stream`; `--play`; `--nowait`; `--wait <duration>` / `--wait=<duration>`; a bare `-` (stdin). `-h` / `--help` returns an error whose text is the `usage` const, like `parseKeep` does.
- [x] Any other argument starting with `-` is `say: unknown flag <arg>` (the caller maps it to exit 2).
- [x] Positionals are joined with a single space into `Options.Text`.
- [x] `--wait` parses with `time.ParseDuration`; `<= 0` or unparsable is an error. `--nowait` and `--wait` together is an error. `--play` without `-o` is an error (`say: --play needs -o`). `-` together with text is an error (`say: - and text together`).
- [x] `Options` (in `internal/say/options.go`): `Text string`, `Stdin bool` (bare `-`), `StdinTTY bool` (set by `main`, not by the parser), `Out string`, `JSON bool`, `Stream bool`, `Play bool`, `Wait time.Duration` (default `-1` = wait forever; `--nowait` → `0`; `--wait d` → `d`). Add `func DefaultOptions() Options`.

## Implementation

1. Read `cmd/cans/main.go` `parseKeep` (line ~109). Copy its shape: a `for i` loop over `args`, `switch` on the arg, `i++` to consume a value, `strings.HasPrefix(a, "--wait=")` style for `=` forms, positionals collected in a slice.
2. Create `internal/say/options.go` with the struct and `DefaultOptions()`. Package doc comment: `// Package say runs cans say: one-shot, file output, stdin, stream.`
3. Create `cmd/cans/say_args.go` with `parseSay`. Start from `say.DefaultOptions()`. Validate at the end (the `--play`/`-o`, `--nowait`/`--wait`, `-`/text rules).
4. Create `cmd/cans/say_args_test.go`, table-driven. Error cases first: unknown flag; `-h`; `--wait bogus`; `--wait 0s`; `--nowait --wait 1s`; `--play` alone; `-` with text. Then: text only; `-o` before text; `-o` after text; `-o=path`; `--out path`; `--json --stream`; `-` alone; `--wait 30s`; multiple positionals joined.
5. Do not wire `main.go` yet — task 02 does. `go vet` is fine with an unused unexported function.

## Done when

- [x] `go test ./cmd/cans/ -run TestParseSay -v` green, every case above present
- [x] `gofmt -l .` empty; `go vet ./...` clean
- [x] `say_args.go` < 120 lines; no function > 50 lines