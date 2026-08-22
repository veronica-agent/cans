---
fest_type: task
fest_id: 01_flock.md
fest_name: flock
fest_parent: 02_lock
fest_order: 1
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:56.606176-06:00
fest_updated: 2026-08-21T06:30:55.366718-06:00
fest_tracking: true
---


# Task: flock

## Objective

`internal/mouth`: the mouth lock. One exclusive `flock` on a file, acquired with wait-forever / try-once / bounded semantics under a `context.Context`, released explicitly, dropped by the kernel if the holder dies.

## Requirements

- [x] API (D011):
  ```go
  package mouth
  var ErrBusy = errors.New("mouth busy")
  type Lock struct{ f *os.File }
  // Path is CANS_HOME/mouth.lock.
  func Path() string
  // Acquire takes an exclusive flock on path. wait < 0 waits forever; wait == 0 tries once;
  // wait > 0 gives up after wait. onWait, if non-nil, runs once the first time the lock is found held.
  func Acquire(ctx context.Context, path string, wait time.Duration, onWait func()) (*Lock, error)
  func (l *Lock) Release() error
  ```
- [x] `syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)` polled every 100 ms (D003). `EWOULDBLOCK` means held. Any other errno is wrapped: `fmt.Errorf("mouth: lock %s: %w", path, err)`.
- [x] `ctx.Done()` while waiting → close the fd, return `ctx.Err()`. Deadline passed, or `wait == 0` and held → close the fd, return `ErrBusy` (unwrapped, so callers use `errors.Is`).
- [x] The file is created `0644` if missing (parent dir `MkdirAll`), and **never deleted**. `Release` = `LOCK_UN` then `Close`; `Release` on a nil `*Lock` is a no-op returning nil.
- [x] No new dependencies; `syscall` only.

## Implementation

1. `internal/mouth/lock.go`, under 120 lines. `Acquire` checks `ctx.Err()` first. Keep the poll loop as its own function (`tryUntil`) so `Acquire` stays under 50 lines.
2. `Path()` uses `ship.Home()` (import `internal/ship`; no cycle — `ship` does not import `mouth`).
3. Tests in `internal/mouth/lock_test.go` (error cases first; `t.TempDir()` for the path; no `time.Sleep` in assertions except the bounded-wait one):
   - held, `wait == 0` → `ErrBusy` (two `Acquire` calls in one process on the same path **do** conflict under `flock` — separate open file descriptions).
   - held, `wait = 150ms` → `ErrBusy`, elapsed ≥ 150 ms.
   - held, `wait < 0`, ctx cancelled after 50 ms → `ctx.Err()`.
   - `onWait` fires exactly once across several polls (an `atomic.Int32`, `wait = 350ms`).
   - release, then `Acquire(wait == 0)` succeeds; the lock file still exists.
   - cross-process `kill -9`: the re-exec helper pattern. `TestHelperHoldLock` returns immediately unless `MOUTH_HELPER=1`; when set it `Acquire`s `MOUTH_LOCK`, prints `held\n`, and sleeps. The parent runs `exec.Command(os.Args[0], "-test.run=TestHelperHoldLock")` with that env, waits for `held`, asserts `Acquire(wait == 0)` → `ErrBusy`, then `cmd.Process.Kill()`, `cmd.Wait()`, and asserts `Acquire(wait == 0)` succeeds.

## Done when

- [x] `go test ./internal/mouth/ -v` green with every case above
- [x] `gofmt -l .` empty; `go vet` clean; `lock.go` < 120 lines; no function > 50 lines