---
fest_type: task
fest_id: 03_flags_booth.md
fest_name: flags_booth
fest_parent: 02_lock
fest_order: 3
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:56.607551-06:00
fest_updated: 2026-08-21T06:39:44.74907-06:00
fest_tracking: true
---


# Task: flags_booth

## Objective

`--nowait` and `--wait` reach the lock from `cans say`; `mouth busy` is exit 75; the booth takes the lock for its whole run (D001).

## Requirements

- [x] `say.Run` builds `tts.Options{Wait: o.Wait, OnWait: func() { fmt.Fprintln(stderr, "waiting for the mouth…") }}` and calls `tts.SayToWith`. `errors.Is(err, mouth.ErrBusy)` → stderr `say: mouth busy`, return `ExitBusy` (75). P0-13: `--nowait` is `Wait == 0`; `--wait 30s` is `Wait == 30s`; default `-1` waits forever.
- [x] `booth.Run`: replace `tts.Open(ctx)` with `tts.OpenWith(ctx, tts.Options{Wait: -1, OnWait: func() { fmt.Fprintln(os.Stderr, "waiting for the mouth…") }})` — **before** `tea.NewProgram`, so the line is visible in the terminal and the TUI opens only once the mouth is held. The existing `defer sess.Close()` releases at exit. `CANS_SAY_BIN` path unchanged (no session, no lock).
- [x] Exit 75 is documented in the `usage` const in one line: `exit 75 when another cans holds the mouth and --nowait was set`.

## Implementation

1. `internal/say/say.go`: a `lockOpts(o Options, stderr io.Writer) tts.Options` helper; map `ErrBusy` in one place (`exitFor(err error) int`).
2. `internal/booth/booth.go` `Run`: a three-line change. Do not touch the model.
3. Tests (error first) in `internal/say/say_test.go`:
   - hold the lock in-process: `lk, _ := mouth.Acquire(ctx, mouth.Path(), 0, nil)` (with `CANS_HOME` temp); `Run` with the fake worker and `Wait: 0` → 75, stderr contains `mouth busy`, stdout empty.
   - `Wait: 200ms` → 75 and elapsed ≥ 200 ms; stderr contains `waiting for the mouth…` exactly once.
   - `lk.Release()`, `Run` again with `Wait: 0` → 0.
   - `cmd/cans/main_test.go`: `run([]string{"say", "--nowait", "x"})` against a held lock (temp `CANS_HOME`, fake worker) → 75.

## Done when

- [x] Manual, recorded in the testing gate: terminal A `./bin/cans` (booth, leave it open); terminal B `./bin/cans say --nowait "Put the cans on."` → prints `say: mouth busy`, `echo $?` → 75; terminal B `./bin/cans say "Put the cans on."` → prints `waiting for the mouth…` and speaks after the booth is closed with Esc.
- [x] `CANS_NOPLAY=1 go test ./...` green; `gofmt -l .` empty; `go vet` clean