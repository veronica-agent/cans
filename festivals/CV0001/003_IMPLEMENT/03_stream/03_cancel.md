---
fest_type: task
fest_id: 03_cancel.md
fest_name: cancel
fest_parent: 03_stream
fest_order: 3
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:56.825156-06:00
fest_updated: 2026-08-21T07:02:25.162165-06:00
fest_tracking: true
---


# Task: cancel

## Objective

Ctrl-C during a stream stops cleanly: finished wavs stay, the worker exits, the lock is released, stderr says where it stopped, exit 130 (D008).

## Requirements

- [x] `cmd/cans/main.go`: for `say`, `ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM); defer stop()`. The booth keeps its own Ctrl-C handling (bubbletea).
- [x] `runStream` checks `ctx.Err()` before reading each next line. When the in-flight `SayTo` returns `ctx.Err()` (the worker client's `readLine` selects on `ctx.Done()`), do **not** count it as a failed line.
- [x] On cancellation: stop reading stdin; `sess.Close()` runs (deferred — sends `shutdown`, waits for the worker, releases the lock); stderr `interrupted after line N` where N is the last **completed** stdin line; return `ExitInterrupted` (130). Completed wavs are left in place.
- [x] One-shot: cancelled while waiting for the lock (`mouth.Acquire` returns `ctx.Err()`) → 130 with `say: interrupted` on stderr. Cancelled mid-synthesis → 130, temp wav removed if it exists.
- [x] Cancel lands **between** requests at the worker: the worker has no mid-synth abort, so it finishes the utterance it is on before `shutdown` takes effect. The README states this (04_tape). Nothing in this task pretends otherwise.

## Implementation

1. `internal/say/stream.go`: a `ctx.Err()` check at the top of the loop, and an `errors.Is(err, context.Canceled)` branch in the per-line error path that breaks instead of counting.
2. `cmd/cans/main.go`: the `NotifyContext` pair; pass `ctx` into `say.Run`.
3. Test (the important one) in `stream_test.go`, fake worker:
   - `pr, pw := io.Pipe()`; `stdout` is a mutex-guarded buffer. Goroutine: write `"a\n"`; poll the buffer until the first record appears (a loop with a 5 s cap, 10 ms steps); `cancel()`; write `"b\n"` and keep the pipe open.
   - `Run` returns **130**; `001.wav` exists; `002.wav` does not; stderr contains `interrupted after line 1`.
   - After `Run` returns: `mouth.Acquire(ctx, mouth.Path(), 0, nil)` succeeds (lock released) and the wrapper counter file still has one line (no restart).
   - Close the pipe writer in `t.Cleanup`.
4. One-shot cancel test: hold the lock, `Run` with `Wait: -1` and a ctx cancelled after 50 ms → 130.

## Done when

- [x] Tests green; `CANS_NOPLAY=1 go test ./...` green; `gofmt -l .` empty; `go vet` clean
- [x] Manual, recorded in the testing gate: `seq 1 20 | sed 's/^/Line /' | ./bin/cans say --stream -o '/tmp/cans-out/%03d.wav'`, Ctrl-C after two files → `interrupted after line 2`, `echo $?` → 130, `ls /tmp/cans-out` shows 001 and 002, `pgrep -f qwen3-tts-worker` → nothing, and the next `./bin/cans say "x"` starts at once