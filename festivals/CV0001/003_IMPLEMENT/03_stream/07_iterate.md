---
fest_type: gate
fest_id: 07_iterate.md
fest_name: Review Results and Iterate
fest_parent: 03_stream
fest_order: 7
fest_status: completed
fest_autonomy: medium
fest_gate_id: iterate
fest_gate_type: iterate
fest_managed: true
fest_created: 2026-08-21T05:04:57.508847-06:00
fest_updated: 2026-08-21T15:38:55.178785-06:00
fest_tracking: true
fest_version: "1.0"
---


# Gate: Review Results and Iterate

Address all findings from testing and code review. Iterate until the sequence meets quality standards.

## Findings to Address

### From Testing

- [x] Real-mouth Ctrl-C re-check on the post-iterate binary (load 10.5): `seq 1 20 | … --stream -o '%03d.wav'`, SIGINT after `002.wav` → exit **130**, stderr exactly `interrupted after line 2`, `ls` = 001 002 (003 in flight, dropped per D014), `pgrep` empty at once, next one-shot exit 0 in 14 s with no `waiting for the mouth…`. One-shot gate: `ttfa_ms=5652`, unchanged from baseline.

### From the orchestrator's read of the iterate diff

- [x] `internal/say/stream.go` `scan`: Ctrl-C while blocked in `sc.Scan()` (idle between lines) let the producer's death close the pipe and the loop report a normal exit. Added a post-loop `ctx.Err()` check → `interruptedAfter(last)`, exit 130. `go test -race ./internal/say` green.

### From Code Review

**Critical**

- [x] `session.go:67` / `worker.go:46` — Ctrl-C SIGKILLed the worker. Resolved the
  way D014 amends D008: `StartWorker` now sets `cmd.Cancel` (SIGTERM) and
  `cmd.WaitDelay = 2s`, so cancel terminates the process cleanly and SIGKILLs
  only a worker that ignores SIGTERM for two seconds. `Session` keeps the
  session context's `Done` channel so `Close` swallows the refused `shutdown`
  write and the `signal: terminated` exit on the cancel path — 130 with no
  stderr noise past the `interrupted …` line.

**Suggestions**

- [x] `stream_test.go:178` — the post-cancel `pw.Write` moved to a goroutine and
  the pipe is closed with `io.EOF` after `<-done`, so it can never wedge the
  package.
- [x] `stream.go:107` — D007 / P0-8 now covered by `TestStreamNoOutPlaysEachLine`:
  `o.Out == ""`, two `ttfa_ms=` lines on stdout, `TMPDIR` scanned for surviving
  `cans-*.wav` (none).
- [x] `say_test.go:279` — `TestRunOnceCancelledMidSynthesis` cancels *inside*
  `SayTo`: the fake worker's new `block` branch reports it reached synthesis by
  touching `CANS_FAKE_BLOCK_FILE`, then waits to be signalled. 130, no temp wav,
  and the mouth lock is reacquirable afterwards. `TestStreamCancelBeforeAnyLine`
  is the same shape for the stream.
- [x] `stream.go:83` — `sc.Err()` is now `say: stdin: %v`, covered by
  `TestStreamLongLineReportsStdin` (a >1 MiB line).
- [x] `stream.go:126` — `interruptedAfter` special-cases 0 as
  `interrupted before the first line`, and a reported failure now advances
  `last`, so N is the last line *fully processed*, spoken or failed (D014).
  Pinned by `TestStreamCancelBeforeAnyLine` and `TestStreamCancelAfterFailedLine`.
- [x] `say_args.go:87` — `validateSay` rejects `--stream` with argv text:
  `say: --stream reads stdin; drop the text`, with a case in the parse-error table.
- [x] `main.go:54` — the say path moved into `runSay`, which starts a goroutine
  that calls `stop()` once `ctx.Done()` fires, restoring the default disposition
  so a second Ctrl-C ends the process at once. `stop()` also cancels `ctx`, so
  the goroutine always ends with `runSay`. `main.go` is 148 lines.
- [x] `stream.go:103` / `say.go:75` — the play/`RemoveTemp` tail is one helper,
  `playTail(o, wav)`, used by both `runOnce` and the stream. `speakLine`'s 9
  parameters became the `streamer` struct (worker, throat, options, writers) with
  `speak(ctx, line, lineNo, idx)`. `emitStream` is `(*streamer).emit`.

## Verification

```
gofmt -l .                                                        → empty
go vet ./...                                                      → clean
go build -o /dev/null ./cmd/cans                                  → ok
CANS_NOPLAY=1 go test ./...                                       → 10/10 ok
CANS_NOPLAY=1 go test -race -count=2 ./internal/say ./internal/tts → ok (14.2s / 5.6s)
CANS_NOPLAY=1 go test -race -count=2 ./cmd/...                    → ok (1.9s)
```

`internal/tts/worker.go` grew from 196 to 201 lines (the four `Cancel`/`WaitDelay`
lines plus the `syscall` import). Every non-table function is under 50 lines; the
largest touched file is `internal/say/say_test.go` at 419. No new module
dependency; the professional-surface grep over `internal/say`, `internal/tts` and
`cmd/cans` is empty. The real mouth was not run and `./bin/cans` was not rebuilt.

## Definition of Done

- [x] All critical findings fixed
- [x] All tests pass after changes
- [x] Linting passes
- [x] Code review findings addressed
- [x] Ready to commit