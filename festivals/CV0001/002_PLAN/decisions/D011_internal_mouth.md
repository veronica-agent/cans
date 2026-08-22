# D011 — internal/mouth is the lock

**Decision:** Package `internal/mouth`: `Acquire(ctx, path, wait time.Duration, onWait func()) (*Lock, error)`, `ErrBusy`, `(*Lock).Release()`. `wait < 0` waits forever, `wait == 0` is `--nowait`, `wait > 0` is `--wait`. `onWait` fires once, the first time the lock is found held. `tts.OpenWith` is the only production caller.

**Why:** One small package with one job, testable in-process (two opens of the same file in one process do conflict under `flock`) and across processes (a re-exec'd test helper holds the lock, gets `kill -9`'d, and the next acquire succeeds).

**Not:** A lock in `ship` (it is not payload) or in `tts` (which is already the worker client).
