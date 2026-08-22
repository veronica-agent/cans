---
fest_type: phase
fest_id: 004_REVIEW
fest_name: REVIEW
fest_parent: cans-v2-CV0001
fest_order: 4
fest_status: pending
fest_created: 2026-08-21T05:03:56.979736-06:00
fest_phase_type: review
fest_tracking: true
---

# Phase Goal: 004_REVIEW

**Phase:** 004_REVIEW | **Status:** Pending | **Type:** Review

## Review Objective

**Primary Goal:** The ship bar from design-recommend.md, identity, fest validate, and the PR from cans-v2 to main.

**Context:** Five sequences on `cans-v2` turned `cans say` into a unix primitive with a mouth lock. This review runs the ship bar from `design-recommend.md` on the real mouth, checks identity and the public surface, re-syncs the snapshot, and opens the one PR.

## What's Being Reviewed

Items that must pass this review:

- The 13-item bar below. Commands and verbatim output: `BAR.md`.

<!-- Add more review items as needed -->

## Review Criteria

Criteria each item must meet:

- [x] Each item is a command a stranger can re-run, with its output recorded verbatim rather than paraphrased.
- [x] Real-mouth items taken one at a time, worker count 0 and 1-minute load < 16 before each.
- [x] Worker counts taken with an anchored `ps` counter, never `pgrep -f` (which self-matches and cross-matches other samplers).
- [x] Anything that did not match expectation is recorded as it happened, not smoothed over.

<!-- Add more criteria as needed -->

## Stakeholder Sign-off

| Stakeholder | Role | Status | Date |
|-------------|------|--------|------|
| Opus review agent | Orchestrator | [x] Approved | 2026-08-21 |
| Operator | Veto after the fact | [ ] Pending | |

<!-- Add rows for each required sign-off -->

## Approval Gates

Gates that must pass before review completion:

- [x] Bar items 1–13 pass (`BAR.md`).
- [x] `docs/pipe.gif` re-cut under load < 16 — the carry-over from `04_tape`.
- [x] `festivals/CV0001/` re-synced with the recorded rsync; `fest validate` 100/100; greps clean; `diff -r` empty.
- [x] PR #14 body rewritten from `BAR.md` and CI green.

<!-- Add more gates as needed -->

## Go/No-Go Decision

**Decision:** [x] GO / [ ] NO-GO

**Conditions for GO:**
- [x] All review criteria passed
- [x] All stakeholder sign-offs received (orchestrator; the operator may veto after the fact)
- [x] All approval gates satisfied

**If NO-GO, actions required:**
- Document blockers
- Return to relevant implementation tasks
- Schedule re-review

## Notes

`BAR.md` carries the commands and the verbatim output for all 13 items, plus the gif re-cut table and a closing section on the two known engine behaviours. Nothing in this phase changed code — the only tracked file it touches in the worktree is `docs/pipe.gif` (re-cut) and the `festivals/CV0001/` snapshot.

---

*Review phases validate completed work. All sign-offs required before marking complete.*

## The bar (commands and recorded output live in `BAR.md`)

| # | Check | Pass |
|---|-------|------|
| 1 | `cat lines.txt \| cans say --stream -o 'out/%03d.wav'` writes one wav per line; `pgrep -f qwen3-tts-worker` shows **one** worker throughout; one GGUF load | [x] |
| 2 | `xargs -P 8` over 24 lines completes with one worker resident at every sample, no pageouts | [x] |
| 3 | Ctrl-C mid-stream: completed wavs remain, no orphaned worker, next `cans say` runs immediately, exit 130 | [x] |
| 4 | `kill -9` on a running cans leaves the next run unblocked | [x] |
| 5 | Booth session + background `cans say --nowait` → 75; background without `--nowait` waits with the stderr line | [x] |
| 6 | `cans say "x"` matches `1e8cea2`; `CANS_NOPLAY=1 go test ./...` green; stream path runs on the fake worker | [x] |
| 7 | 50-line stream beats the 50-call loop by the margin recorded in `002_PLAN/inputs/measurements.md` | [x] |
| 8 | Professional-surface grep (`CONTEXT.md §Professional grep`) empty over README, docs/, tapes/, festivals/; exactly one Festival footer | [x] |
| 9 | `just test unit` green; fresh `CANS_HOME` doctor green with the binary outside the checkout | [x] |
| 10 | `git log origin/main..cans-v2 --format='%an <%ae>'` is only Veronica; no `Co-authored-by` | [x] |
| 11 | `fest validate` green on this festival and on `festivals/CV0001/` in the repo | [x] |
| 12 | Snapshot re-synced after this review's statuses are set; committed | [x] |
| 13 | PR opened from `cans-v2` to `main` on `veronica-agent/cans` under `veronica-agent`, body from `BAR.md`, CI green | [x] |

## Sign-off

| Role | Who | Date | Verdict |
|------|-----|------|---------|
| Orchestrator | Opus review agent | 2026-08-21 | **GO — 13 / 13 pass.** Every item run on the real mouth one at a time on a quiet box (1-minute load 4.6–7.6 throughout); commands and verbatim output in `BAR.md`. Two known engine faults reappeared (end-of-speech variance; near-silent wavs returned as success) — recorded, deferred, and named in the PR body; neither is a defect in this branch. |
| Operator | | | (after the fact — may veto) |
