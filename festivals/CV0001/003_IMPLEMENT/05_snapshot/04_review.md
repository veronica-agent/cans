---
fest_type: gate
fest_id: 04_review.md
fest_name: Code Review
fest_parent: 05_snapshot
fest_order: 4
fest_status: completed
fest_autonomy: low
fest_gate_id: review
fest_gate_type: review
fest_managed: true
fest_created: 2026-08-21T05:04:57.561891-06:00
fest_updated: 2026-08-21T18:09:57.506987-06:00
fest_tracking: true
fest_version: "1.0"
---


# Gate: Code Review

Review all code changes in this sequence for quality, correctness, and standards compliance.

## Review Checklist

### Code Quality

- [x] Code is readable and well-organized — no Go code here; the tree itself reads cleanly
- [x] Functions are focused (single responsibility) — n/a
- [x] Naming is clear and consistent
- [x] No unnecessary complexity or duplication

### Standards Compliance

- [x] Linting passes without warnings — `gofmt -l .` re-run silent
- [x] Formatting is consistent
- [x] Project conventions are followed — one exception, see Critical 2

### Error Handling & Security

- [x] Errors are handled appropriately — n/a
- [x] No secrets in code — leak greps clean; only the public `veronica-agent` commit identity appears
- [x] Input validation present where needed — n/a
- [x] No obvious security issues

### Alignment

- [x] Changes align with sequence goal
- [x] No scope creep beyond what was requested — the source scrub was in scope and recorded

## Findings

Reviewed cold. `05_snapshot` adds **no Go code** — the surface is the untracked
`festivals/CV0001/` tree plus the recorded recheck.

**Green, verified independently:** all five exclusions absent; 88 files; the three greps print
nothing / nothing / `1`; grep 2 over the snapshot **source** minus the excluded paths is also
empty, so the copy is not cleaner than its source; no operator account name anywhere; the 14
links in `002_PLAN/decisions/INDEX.md` all resolve; `fest validate` 90/100 with exactly the two
expected warnings; `git status --short` is `?? festivals/CV0001/` alone; `go.mod`/`go.sum` undiffed;
HEAD still `5e8123d`. Re-ran the safe recheck blocks — `gofmt -l .` silent, `just --list vhs`
still carries `pipe`, largest `.go` 236 / 352 / 418 — all match the pasted output. `02_recheck.md` / `03_testing.md` are complete for every required block and
honest, unflattering entries included. Overview, plan, measurements and task files read as a plan
a stranger could follow with `fest` uninstalled.

**Critical Issues:** (must fix)

- `002_PLAN/decisions/D009_public_snapshot.md:3` — the shipped copy is **stale**. `diff -r`
  source-vs-snapshot with the exclusions applied is empty except this one file: the source D009
  was amended at 18:02 (adding `.workflow/` and `.review-*`), the last rsync ran at 18:00, so
  the copy still states the old four-item list — the tree misstates the policy that produced it,
  on the document a stranger reads to understand the scrub. Fix: re-run the rsync (idempotent).
- `003_IMPLEMENT/01_out/.review-01_out.md`, `003_IMPLEMENT/02_lock/.review-02_lock.md` — two
  hidden reviewer scratch files ship. Amended D009 excludes `.review-*`, so this is now a D009
  violation, and they are the only reason `fest validate` scores 90 and not 100. Read both in
  full: ordinary code notes, no phrase-lock hit, no path leak — a policy break, not a leak.
  Fix: add `--exclude '.review-*'` to the rsync.
- `003_IMPLEMENT/05_snapshot/01_snapshot.md:64` — the rsync recorded here is explicitly the
  command `004_REVIEW` re-runs before the PR, and it has no `--exclude '.review-*'`; left as is,
  both findings above reappear after the pre-PR re-sync. Fix: update it, then re-run it.

**Suggestions:** (should consider)

- `FESTIVAL_OVERVIEW.md:57`, `TODO.md:51` (28 files in all) — the tree cites `CONTEXT.md` ~50
  times and the `001_INGEST/input_specs/` design pack ~90 times; D009 excludes both, so those
  pointers dangle for a stranger. No substance is lost — D001–D014 are all present under
  `002_PLAN/decisions/`. Fix: one line in `TODO.md` saying both are campaign-private and the
  decisions live in `002_PLAN/decisions/`, rather than 140 edits.
- `002_PLAN/plan/IMPLEMENTATION_PLAN.md:59` — the `05_snapshot` row restates the pre-amendment
  four-exclusion rsync as "per D009 exclusions", contradicting D009 two directories away. It is
  a plan of record, so rewriting is a judgement call. Fix: cite D009 instead of restating it.
- `002_PLAN/inputs/measurements.md:9,28` — the baseline probe (`scratchpad/measure/main.go`) and
  its build worktree (`.../cans/drop-sidecar`) are not in the public repo, so those two rows are
  the only numbers a stranger cannot reproduce (§Stream is). Fix: note the probe was throwaway.
- `fest.yaml:13,17,20,61` — carries the campaign-side `status_history` paths and names the
  excluded `.festival-checksums.json`. Machine metadata, validate passes; no action needed.
- **For `004_REVIEW` / the operator — `festivals/CA0001/`.** Confirmed the implementer's finding,
  count tightened: grep 2 over `festivals/CA0001/` returns **25 matching lines across 11 files**
  (recorded as 24/11), of which **7 files / 21 lines** are the tracked `001_INGEST/input_specs/`
  pack — the exact directory D009 excludes here; the rest are `FESTIVAL_OVERVIEW.md`,
  `FESTIVAL_RULES.md`, `output_specs/constraints.md`, `003_IMPLEMENT/06_snapshot/01_readme.md`.
  Committed, public, predating this branch and outside the `cans-v2` diff — rightly untouched and
  not a blocker here. Fix: an operator call before the PR — scrub, or accept on record.
- `internal/tts/worker.go` is 201 lines against the 196 `FESTIVAL_RULES.md` pins — inherited from
  `03_stream`, committed, untouched here. `004_REVIEW` should accept it or amend the rule.

## cans-v2 review points

The reviewer is a **different agent** than the implementer and reads `git diff origin/main` cold. Check, and write a finding for each miss:

Points 1–7 and 9 are **n/a**: no Go, flag, or lock changes here; `go.mod`/`go.sum` undiffed.
Point 8 is the whole sequence and is green; point 6's length clause is inherited (last Suggestion).

- `context.Context` is the first parameter on anything that does I/O; `ctx.Err()` checked before long work; cancellation reaches the worker
- stdout carries only `ttfa_ms=`, wav paths, or JSONL; everything else is on stderr
- every flag is one of `-o/--out`, `--json`, `--stream`, `--play`, `--nowait`, `--wait`, `-` — nothing else exists
- the lock is acquired **before** `StartWorker` and released **after** `Client.Close` returns; the lock file is never deleted
- errors are wrapped with the failing operation (`fmt.Errorf("say: %w", err)`)
- no new `go.mod` requires; files < 500 lines; functions < 50 lines; `internal/tts/worker.go` unchanged in length
- tests run on the fake worker with `CANS_NOPLAY=1`; error cases first; no sleeps in assertions
- any README / tape / fixture / help text added is boring and technical and passes the professional-surface grep (`CONTEXT.md §Professional grep`)
- `cans say "x"` is byte-identical in behavior to `1e8cea2`