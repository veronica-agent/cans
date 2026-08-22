---
fest_type: task
fest_id: 02_session_lock.md
fest_name: session_lock
fest_parent: 02_lock
fest_order: 2
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:56.60691-06:00
fest_updated: 2026-08-21T06:33:42.282307-06:00
fest_tracking: true
---


# Task: session_lock

## Objective

Every `tts.Session` is lock-guarded by construction: `OpenWith` acquires the mouth lock **before** `StartWorker`, the `Session` carries it, and `Close` releases it **after** the worker has exited.

## Requirements

- [x] `tts.Options{ Wait time.Duration; OnWait func() }`, `func DefaultOptions() Options` (`Wait: -1`, `OnWait: defaultOnWait`), and `func OpenWith(ctx context.Context, o Options) (*Session, error)`. `Open(ctx)` = `OpenWith(ctx, DefaultOptions())`. `defaultOnWait` writes `waiting for the mouth…` to `os.Stderr` (P0-12).
- [x] Order inside `OpenWith`: `ctx.Err()` → `mouth.Acquire(ctx, mouth.Path(), o.Wait, o.OnWait)` → stat the worker binary → `StartWorker`. If anything after `Acquire` fails, `Release` before returning the error. `ErrBusy` is returned unwrapped (callers map it to 75).
- [x] `Session{c *Client; lock *mouth.Lock}`. `Close`: `err := s.c.Close()` (shutdown + wait), then `s.lock.Release()`, return the worker error if any (D001: lock lifetime == session lifetime).
- [x] `SayToWith(ctx, text, cur, out string, o Options)`; `SayTo` = `SayToWith(…, DefaultOptions())`. The `CANS_SAY_BIN` path takes no lock (no worker).
- [x] `internal/tts/worker.go` is not touched.

## Implementation

1. `internal/tts/session.go`: add `Options`, `DefaultOptions()`, `OpenWith`; `Open` becomes a wrapper. Keep `OpenWith` under 50 lines by moving the worker start into `startSession(ctx, lock *mouth.Lock) (*Session, error)`.
2. `internal/tts/synth.go`: `SayToWith`; existing funcs become wrappers.
3. Tests in `internal/tts/session_test.go` (`CANS_HOME` = temp dir; `CANS_NOPLAY=1`; fake worker built as in `synth_test.go`):
   - **ordering proof** (error case first): open `A := OpenWith(ctx, Options{Wait: -1})` on the fake worker; then set `CANS_WORKER_BIN` to a path that does not exist and call `OpenWith(ctx, Options{Wait: 0})` — it must return `mouth.ErrBusy`, **not** `native mouth missing`. That proves the lock is taken before the worker is looked at.
   - `A.Close()`, then `OpenWith(Options{Wait: 0})` on the fake worker succeeds, and `Close` works.
   - `OnWait` fires once while A is held and a second `OpenWith(Options{Wait: 300ms})` waits and returns `ErrBusy`.
   - `$CANS_HOME/mouth.lock` exists after the first `Open`.

## Done when

- [x] `CANS_NOPLAY=1 go test ./internal/tts/ -v` green, existing tests unchanged
- [x] `./bin/cans say "Put the cans on."` still behaves exactly as before and `ls ~/.cans/mouth.lock` exists afterwards
- [x] `gofmt -l .` empty; `go vet` clean; `session.go` < 200 lines