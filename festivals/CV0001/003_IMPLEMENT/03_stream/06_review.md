---
fest_type: gate
fest_id: 06_review.md
fest_name: Code Review
fest_parent: 03_stream
fest_order: 6
fest_status: completed
fest_autonomy: low
fest_gate_id: review
fest_gate_type: review
fest_managed: true
fest_created: 2026-08-21T05:04:57.503707-06:00
fest_updated: 2026-08-21T14:23:48.325905-06:00
fest_tracking: true
fest_version: "1.0"
---


# Gate: Code Review

Read cold by a different agent: `git diff HEAD` plus the four untracked files
(`internal/say/stream.go`, `stream_test.go`, `template.go`, `template_test.go`),
against `01_stream_loop.md`, `02_out_template.md`, `03_cancel.md`, D005–D008 and
P0-5…P0-9 / P0-16 / P0-19 / P0-20. No code was edited. The real mouth was not run.

## Verification run

```
gofmt -l .                                          → empty
go vet ./...                                        → clean
CANS_NOPLAY=1 go test -count=1 ./...                → 10/10 packages ok (internal/say 11.757s)
CANS_NOPLAY=1 go test -race -count=2 ./internal/say ./internal/tts → ok (13.3s / 6.7s)
git diff origin/main -- go.mod go.sum               → empty
internal/tts/worker.go                              → 196 lines, no diff
```

Longest new function: `scanStream` at 41 lines (`stream.go:50`). Largest touched file:
`say_test.go` at 361 lines. No goroutine, buffer or queue was added around the worker —
`scanStream` sends the next `SayTo` only after the previous returns.

## Review Checklist

### Code Quality

- [x] Code is readable and well-organized
- [x] Functions are focused (single responsibility)
- [x] Naming is clear and consistent
- [x] No unnecessary complexity or duplication

### Standards Compliance

- [x] Linting passes without warnings
- [x] Formatting is consistent
- [x] Project conventions are followed

### Error Handling & Security

- [x] Errors are handled appropriately
- [x] No secrets in code
- [x] Input validation present where needed
- [x] No obvious security issues

### Alignment

- [x] Changes align with sequence goal
- [x] No scope creep beyond what was requested

## Findings

**Critical Issues:** (must fix)

- internal/tts/session.go:67 (reached from internal/say/stream.go:42) — Ctrl-C SIGKILLs the worker instead of shutting it down: `StartWorker` builds the process with `exec.CommandContext(ctx, …)` (worker.go:46) and `cmd/cans/main.go:54` now hands `say` the first cancellable ctx, so Go's exec watchdog calls `Process.Kill()` the moment SIGINT lands and the deferred `sess.Close()` (stream.go:46) writes `{"type":"shutdown"}` into a dead pipe — `cmd.Wait()` returns `signal: killed` (confirmed with a standalone `CommandContext` repro) — which contradicts D008 ("Not: killing the worker mid-utterance"), `03_cancel.md` req 5 ("finishes the utterance it is on before `shutdown` takes effect") and P0-19, and 04_tape's README is contracted to repeat that claim. Fix: start the worker on `context.WithoutCancel(ctx)` in `startSession` **and** set `cmd.WaitDelay` (unread PCM will otherwise fill the pipe and wedge `Wait`), or amend D008 and the README line to say cancel kills the worker.

**Suggestions:** (should consider)

- internal/say/stream_test.go:178 — `pw.Write([]byte("b\n"))` runs on the test goroutine after `cancel()`; if `Run` already returned at the top-of-loop `ctx.Err()` check (stream.go:57) nothing reads the pipe and the write blocks forever, wedging the package for the full test timeout — narrow window, but the risk register lists stream-test flakiness as open. Fix: do that write in a goroutine (or `pw.CloseWithError`) before receiving on `done`.
- internal/say/stream.go:107 — D007 / P0-8 (stream with no `-o`: play each line, `ttfa_ms=` per line, temp removed) is implemented but untested; every stream test sets `o.Out`. Fix: add a stream case with `o.Out == ""` asserting two `ttfa_ms=` lines and no surviving temp wav.
- internal/say/say_test.go:279 — `03_cancel.md` req 4 is only half covered: cancelled-while-waiting-for-the-lock is tested, cancelled-mid-synthesis (130, temp wav removed) is not. Fix: add a one-shot test that cancels during `SayTo` on the fake worker.
- internal/say/stream.go:83 — `sc.Err()` is printed bare (`bufio.Scanner: token too long` for a >1 MiB line), not wrapped with the failing operation the way `input.go:29` does. Fix: `fmt.Fprintf(stderr, "say: stdin: %v\n", err)`.
- internal/say/stream.go:126 — a cancel before any line completes prints `interrupted after line 0`, and a line that failed does not advance `last`, so the message can name an earlier line than the last one processed. Fix: special-case 0 and decide whether a reported failure counts as completed.
- cmd/cans/say_args.go:87 — `cans say "x" --stream` silently drops the argv text, while `-` plus text is already a usage error (say_args.go:94). Fix: reject `o.Stream && o.Text != ""` in `validateSay`.
- cmd/cans/main.go:54 — after the first SIGINT the `NotifyContext` ctx is already done, so a second Ctrl-C is swallowed and the user cannot interrupt a wedged `Close`. Fix: call `stop()` once cancellation is observed, restoring the default disposition.
- internal/say/stream.go:103 and say.go:75 — `emitStream` + the play/`RemoveTemp` tail duplicate `emit` + `runOnce`'s tail; `speakLine` also carries 8 parameters. Fix: fold the shared tail into one helper if it grows again.

## cans-v2 review points

1. **PASS** — `ctx` is first on `runStream` / `scanStream` / `speakLine` / `SayTo`; `ctx.Err()` is checked at the top of the loop (stream.go:57) and again after `Scan` (63); `SayTo`'s `ctx.Err()` reaches `speakLine` unwrapped enough for `errors.Is` (session.go:99 uses `%w`, worker_pcm.go:17/21 return it bare) and is not counted as a failed line (stream.go:72, 96).
2. **PASS** — stdout carries only `ttfa_ms=`, wav paths, or JSONL (stream.go:130–146), through a `bufio.Writer` flushed after every record; every error, the wait line and `interrupted after line N` go to stderr.
3. **PASS** — no new flags. `--stream` was already in `parseSay` from 01_out; the set is still `-o/--out`, `--json`, `--stream`, `--play`, `--nowait`, `--wait`, `-`.
4. **PASS** — `OpenWith` (lock → `StartWorker`) is unchanged and called once per stream (stream.go:42); `defer sess.Close()` releases after `Client.Close`; `TestStreamCancel:192` reacquires the lock after 130, and `TestStreamOneWorker` proves exactly one worker start over five lines via an `exec` wrapper counter.
5. **PASS** — `checkOut` runs at `say.go:20`, before `doctor.Prepare` and before the lock; messages are exactly `say: -o needs one %d in --stream` and `say: -o template needs --stream`, both exit 2; `%%` is a literal and `%s` / two verbs / no verb / a verb in one-shot are all rejected (template_test.go:9).
6. **PASS** — `go.mod`/`go.sum` unchanged; every file under 500 lines and every function under 50; `internal/tts/worker.go` untouched at 196.
7. **PASS** — all tests use the fake worker with `CANS_NOPLAY=1`; error cases lead both new test files; the cancel test polls a mutex-guarded buffer with a 5 s cap rather than sleeping and hoping; the fakeworker change is the single `"fail"` branch (4 lines) and nothing else.
8. **PASS** — no README, tape or help text was added in this sequence; the new user-visible strings (`line N: …`, `interrupted after line N`, the two `-o` messages) are boring and technical, and the professional-surface grep over `internal/say`, `cmd/cans` and `internal/tts/testdata` is empty. Test fixtures are `a`/`b`/`fail`/`Put the cans on.`
9. **PASS** — `cans say "x"` is unchanged: `checkOut("" , false)` returns nil, `outPath("", 0)` returns `""`, and `runOnce` is otherwise byte-identical. The only new one-shot behavior is the `context.Canceled` branch in `exitFor` (say.go:100), unreachable before this sequence because `main` passed `context.Background()`. Note that SIGTERM is now a graceful 130 instead of default signal death.