---
fest_type: gate
fest_id: 05_review.md
fest_name: Code Review
fest_parent: 02_lock
fest_order: 5
fest_status: completed
fest_autonomy: low
fest_gate_id: review
fest_gate_type: review
fest_managed: true
fest_created: 2026-08-21T05:04:57.435924-06:00
fest_updated: 2026-08-21T06:55:20.007767-06:00
fest_tracking: true
fest_version: "1.0"
---


# Gate: Code Review

Review all code changes in this sequence for quality, correctness, and standards compliance.

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

Cold review by a different agent (`02_lock/.review-02_lock.md`). All nine festival points PASS. `worker.go` still 196.

**Critical Issues:** (must fix)

None.

**Suggestions:** (should consider)

1. `internal/mouth/lock_test.go:96` — `TestHelperHoldLock` dropped the `*Lock` (`_ = l`) then `select {}`. GC closing the fd would release the flock while the helper still lived. `defer runtime.KeepAlive(l)`.
2. `internal/mouth/lock.go:75` — try flock before treating a positive wait as expired, so `--wait 1ns` on a free mouth still tries once. `onWait` only when the caller is about to block.

## cans-v2 review points

The reviewer is a **different agent** than the implementer and reads `git diff origin/main` cold. Check, and write a finding for each miss:

- `context.Context` is the first parameter on anything that does I/O; `ctx.Err()` checked before long work; cancellation reaches the worker
- stdout carries only `ttfa_ms=`, wav paths, or JSONL; everything else is on stderr
- every flag is one of `-o/--out`, `--json`, `--stream`, `--play`, `--nowait`, `--wait`, `-` — nothing else exists
- the lock is acquired **before** `StartWorker` and released **after** `Client.Close` returns; the lock file is never deleted
- errors are wrapped with the failing operation (`fmt.Errorf("say: %w", err)`)
- no new `go.mod` requires; files < 500 lines; functions < 50 lines; `internal/tts/worker.go` unchanged in length
- tests run on the fake worker with `CANS_NOPLAY=1`; error cases first; no sleeps in assertions
- any README / tape / fixture / help text added is boring and technical and passes the professional-surface grep (`CONTEXT.md §Professional grep`)
- `cans say "x"` is byte-identical in behavior to `1e8cea2`