---
fest_type: task
fest_id: 02_internal_say.md
fest_name: internal_say
fest_parent: 01_out
fest_order: 2
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:56.482119-06:00
fest_updated: 2026-08-21T05:37:13.353773-06:00
fest_tracking: true
---


# Task: internal_say

## Objective

Create `internal/say.Run` and move the `case "say"` flow into it; add `tts.SayTo` so a caller can name the output path; `-o` writes there and never deletes, `--play` plays after writing. `cans say "x"` stays identical.

## Requirements

- [x] `func Run(ctx context.Context, o Options, stdin io.Reader, stdout, stderr io.Writer) int` in `internal/say/say.go`. Exit constants in `internal/say/exit.go`: `ExitOK = 0`, `ExitFail = 1`, `ExitUsage = 2`, `ExitBusy = 75`, `ExitInterrupted = 130` (75 and 130 are declared now, wired in 02_lock / 03_stream).
- [x] `tts.SayTo(ctx, text, cur, out string) (Result, error)` and `(*Session).SayTo(ctx, text, cur, out string)`. `out == ""` → the temp path exactly as today. `out != ""` → `os.MkdirAll(filepath.Dir(out), 0o755)` (D012) then `audio.WritePCM16(out, …)`; `Result.Wav == out`. `Say` / `SayWith` / `(*Session).Say` become one-line wrappers that pass `""`.
- [x] `CANS_SAY_BIN` path: when `out != ""`, copy the script's wav to `out` (small `copyFile` helper in `synth_bin.go`) and return `Wav: out`. The seam keeps working for tests and tapes.
- [x] In `Run`, one-shot flow in this order: `doctor.Prepare(ctx, stderr)` (returns nil early when `CANS_SAY_BIN` is set — unchanged); empty `o.Text` → `say: missing text` on stderr, `ExitUsage` (stdin arrives in task 03); `keep.Load()`; `tts.SayTo(ctx, text, cur, o.Out)`; error → stderr, `ExitFail`.
- [x] No `-o`: print `ttfa_ms=N` to stdout, `play.File`, `tts.RemoveTemp`, play error → `ExitFail`. Exactly today's behavior and order (P0-17).
- [x] `-o`: print the path to stdout (JSON comes in task 03); if `o.Play`, `play.File(out)`; **never** `RemoveTemp`.
- [x] `cmd/cans/main.go`: add `stdin io.Reader = os.Stdin` next to `stdout` / `stderr`; `case "say"` becomes: `o, err := parseSay(args[1:])` → on error print to stderr and return 2; `return say.Run(context.Background(), o, stdin, stdout, stderr)`.

## Implementation

1. `internal/tts/session.go`: rename the body of `Say` into `SayTo`, branch on `out`. Keep `audio.Clean` and the sample-rate default exactly where they are.
2. `internal/tts/synth.go`: `SayTo` mirrors `SayWith` (`CANS_SAY_BIN` → `sayBin` then copy if `out != ""`; else `Open`, `defer Close`, `sess.SayTo`). `SayWith` = `SayTo(ctx, text, cur, "")`.
3. `internal/say/say.go`: `Run` dispatches to `runOnce` (stream comes later). Keep `Run` under 30 lines; `runOnce` under 50.
4. `cmd/cans/main.go`: delete the inlined flow; wire `parseSay` + `say.Run`. `main.go` shrinks.
5. Tests, error cases first, in `internal/say/say_test.go`:
   - `Out` under an unwritable path (a regular file used as a directory) → `ExitFail`, stderr non-empty, stdout empty.
   - `CANS_SAY_BIN` fake script (copy the shape from `internal/tts/synth_test.go` `TestSayMockBin`): `Out == ""` → stdout is `ttfa_ms=12\n`; `Out = t.TempDir()/out/take.wav` → file exists with a valid header (`audio.HeaderOK`), stdout is the path + newline, original fake wav untouched.
   - Fake worker (`CANS_WORKER_BIN`): build `internal/tts/testdata/fakeworker` into a temp dir with `go build` the way `synth_test.go` does (copy that ~10-line helper; `CANS_WORKER_MODELS` any temp dir; `CANS_HOME` temp; `CANS_NOPLAY=1`). `Out` set → wav written at `Out`, stdout is the path.
   - `cmd/cans/main_test.go`: `run([]string{"say"})` → 2; `run([]string{"say", "--bogus"})` → 2 with stderr mentioning the flag.

## Done when

- [x] `CANS_NOPLAY=1 go test ./...` green (existing tts tests unchanged)
- [x] `just build quick && ./bin/cans say "Put the cans on."` prints `ttfa_ms=N`, plays, leaves no `cans-*.wav` in `$TMPDIR` — same as `1e8cea2`
- [x] `./bin/cans say -o /tmp/cans-out/take.wav "Put the cans on."` prints `/tmp/cans-out/take.wav`, the file plays with `afplay`, nothing on stderr
- [x] `gofmt -l .` empty; `go vet` clean; every touched file < 500 lines, functions < 50