---
fest_type: task
fest_id: 03_stdin_json_exit.md
fest_name: stdin_json_exit
fest_parent: 01_out
fest_order: 3
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:56.482892-06:00
fest_updated: 2026-08-21T06:07:27.029893-06:00
fest_tracking: true
---


# Task: stdin_json_exit

## Objective

stdin as one utterance, the TTY rule, `--json` records, and the exit-code / stream discipline — stdout is data, stderr is prose.

## Requirements

- [x] Text resolution in `Run`: if `o.Stdin` (bare `-`) **or** `o.Text == ""`: when `o.StdinTTY && !o.Stdin` → `say: missing text`, `ExitUsage` (P0-4, the usage error it is today, never a hang). Otherwise read stdin (`io.ReadAll(io.LimitReader(stdin, 4<<20))`), `strings.TrimSpace`; empty → `say: empty text`, `ExitUsage`. The whole input is one utterance (P0-3).
- [x] `main` sets `o.StdinTTY` from a real isatty ioctl (`TIOCGETA`), not `Stat` `ModeCharDevice` — `/dev/null` is a char device but not a TTY, and `cans say < /dev/null` must be `say: empty text`. Tests set the field directly.
- [x] `--json`: one record on stdout, `json.NewEncoder(stdout).Encode(r)` where `r` is the `tts.Result` (tags already `wav`, `ttfa_ms`, `sample_rate`) — D006. Without `-o` the wav is a temp file removed after playback; say so in the README later (04_tape), not in code.
- [x] With `--json` and no `-o`: still play and `RemoveTemp` (v1 semantics; the record is for `ttfa_ms`). With `--json -o`: the record's `wav` is the `-o` path.
- [x] Exit mapping: every usage problem → `ExitUsage` (2); every runtime problem → `ExitFail` (1); messages prefixed `say:` on stderr; stdout never carries a message (P0-15, P0-16).

## Implementation

1. `internal/say/input.go`: `func resolveText(o Options, stdin io.Reader) (string, int)` returning text and an exit code (0 on success). Keep the TTY rule and the empty rule here; table-test it.
2. `internal/say/say.go`: call `resolveText` first; emit with a small `emit(stdout, o, r)` helper: JSON → encoder; `-o` → path; else `ttfa_ms=`.
3. `cmd/cans/main.go`: set `o.StdinTTY` after `parseSay` (a 4-line helper `stdinIsTTY() bool`).
4. Tests (error first) in `internal/say/input_test.go` and `say_test.go`:
   - `Text == ""`, `StdinTTY = true` → 2, stderr `say: missing text`, stdin **not read** (use a reader whose `Read` fails the test).
   - piped empty / whitespace stdin → 2, `say: empty text`.
   - `Stdin = true` with `strings.NewReader("Put the cans on.\n")` → the fake say bin receives that text (have the fake script write `"$@"` to a file and assert on it).
   - `--json` with fake say bin → stdout parses as JSON with `ttfa_ms == 12` and `sample_rate == 24000`; nothing else on stdout.
   - `--json -o` → record `wav` equals the out path and the file exists.

## Done when

- [x] `echo "Put the cans on." | ./bin/cans say -o /tmp/cans-out/a.wav --json` prints exactly one JSON line; `./bin/cans say < /dev/null` exits 2 with `say: empty text`; `./bin/cans say` in a terminal exits 2 with `say: missing text`
- [x] `CANS_NOPLAY=1 go test ./...` green; `gofmt -l .` empty; `go vet` clean