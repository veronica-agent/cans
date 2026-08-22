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

- see BAR.md

<!-- Add more review items as needed -->

## Review Criteria

Criteria each item must meet:

- [ ] see BAR.md

<!-- Add more criteria as needed -->

## Stakeholder Sign-off

| Stakeholder | Role | Status | Date |
|-------------|------|--------|------|
| see BAR.md | see BAR.md | [ ] Approved | |

<!-- Add rows for each required sign-off -->

## Approval Gates

Gates that must pass before review completion:

- [ ] see BAR.md

<!-- Add more gates as needed -->

## Go/No-Go Decision

**Decision:** [ ] GO / [ ] NO-GO

**Conditions for GO:**
- [ ] All review criteria passed
- [ ] All stakeholder sign-offs received
- [ ] All approval gates satisfied

**If NO-GO, actions required:**
- Document blockers
- Return to relevant implementation tasks
- Schedule re-review

## Notes

see BAR.md

---

*Review phases validate completed work. All sign-offs required before marking complete.*

## The bar (commands and recorded output live in `BAR.md`)

| # | Check | Pass |
|---|-------|------|
| 1 | `cat lines.txt \| cans say --stream -o 'out/%03d.wav'` writes one wav per line; `pgrep -f qwen3-tts-worker` shows **one** worker throughout; one GGUF load | [ ] |
| 2 | `xargs -P 8` over 24 lines completes with one worker resident at every sample, no pageouts | [ ] |
| 3 | Ctrl-C mid-stream: completed wavs remain, no orphaned worker, next `cans say` runs immediately, exit 130 | [ ] |
| 4 | `kill -9` on a running cans leaves the next run unblocked | [ ] |
| 5 | Booth session + background `cans say --nowait` → 75; background without `--nowait` waits with the stderr line | [ ] |
| 6 | `cans say "x"` matches `1e8cea2`; `CANS_NOPLAY=1 go test ./...` green; stream path runs on the fake worker | [ ] |
| 7 | 50-line stream beats the 50-call loop by the margin recorded in `002_PLAN/inputs/measurements.md` | [ ] |
| 8 | Professional-surface grep (`CONTEXT.md §Professional grep`) empty over README, docs/, tapes/, festivals/; exactly one Festival footer | [ ] |
| 9 | `just test unit` green; fresh `CANS_HOME` doctor green with the binary outside the checkout | [ ] |
| 10 | `git log origin/main..cans-v2 --format='%an <%ae>'` is only Veronica; no `Co-authored-by` | [ ] |
| 11 | `fest validate` green on this festival and on `festivals/CV0001/` in the repo | [ ] |
| 12 | Snapshot re-synced after this review's statuses are set; committed | [ ] |
| 13 | PR opened from `cans-v2` to `main` on `veronica-agent/cans` under `veronica-agent`, body from `BAR.md`, CI green | [ ] |

## Sign-off

| Role | Who | Date | Verdict |
|------|-----|------|---------|
| Orchestrator | | | |
| Operator | | | (after the fact — may veto) |
