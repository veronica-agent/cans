---
# Template metadata (for fest CLI discovery)
id: QUALITY_GATE_REVIEW
aliases:
  - code-review
  - qg-review
description: Standard quality gate task for code review

# Fest document metadata (becomes document frontmatter)
fest_type: gate
fest_id: <no value>
fest_name: Code Review
fest_parent: <no value>
fest_order: <no value>
fest_gate_type: review
fest_autonomy: low
fest_status: pending
fest_tracking: true
fest_created: 2026-08-21T04:32:49-06:00
---

# Gate: Code Review

Review all code changes in this sequence for quality, correctness, and standards compliance.

## Review Checklist

### Code Quality

- [ ] Code is readable and well-organized
- [ ] Functions are focused (single responsibility)
- [ ] Naming is clear and consistent
- [ ] No unnecessary complexity or duplication

### Standards Compliance

- [ ] Linting passes without warnings
- [ ] Formatting is consistent
- [ ] Project conventions are followed

### Error Handling & Security

- [ ] Errors are handled appropriately
- [ ] No secrets in code
- [ ] Input validation present where needed
- [ ] No obvious security issues

### Alignment

- [ ] Changes align with sequence goal
- [ ] No scope creep beyond what was requested

## Findings

Document any issues that must be addressed before commit.

**Critical Issues:** (must fix)

**Suggestions:** (should consider)

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
