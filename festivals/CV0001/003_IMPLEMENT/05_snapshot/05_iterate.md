---
fest_type: gate
fest_id: 05_iterate.md
fest_name: Review Results and Iterate
fest_parent: 05_snapshot
fest_order: 5
fest_status: completed
fest_autonomy: medium
fest_gate_id: iterate
fest_gate_type: iterate
fest_managed: true
fest_created: 2026-08-21T05:04:57.562955-06:00
fest_updated: 2026-08-21T18:13:58.775332-06:00
fest_tracking: true
fest_version: "1.0"
---


# Gate: Review Results and Iterate

Address all findings from testing and code review. Iterate until the sequence meets quality standards.

## Findings to Address

### From Testing

- [x] **No defects.** `03_testing` was green in every block — ten packages on the fake worker, `gofmt`/`go vet` silent, `go.mod`/`go.sum` undiffed, the real-mouth one-shot at `ttfa_ms=5669` / exit 0. Nothing to iterate on.
- [x] **Stale temp file** `/var/folders/…/T/cans-say.F26P2oxvoX`, 355 B, dated Aug 19 — noted in the gate as pre-dating this branch. Confirmed not ours (the run removed its own temp wav); left alone, no fix warranted.

### From Code Review

All three Criticals were **one root cause**: the recorded rsync did not implement D009 as amended at 18:02, and the copy on disk was made before the amendment.

- [x] **Critical 1 — `D009_public_snapshot.md` stale in the copy.** The source was amended at 18:02, the last rsync ran at 18:00, so the shipped tree stated the old four-item exclusion list — the document a stranger reads to understand the scrub, misstating the policy that produced the tree. Fixed by re-running the rsync; `diff -r` source→copy with the exclusions applied is now empty, so nothing else was stale either.
- [x] **Critical 2 — two hidden reviewer scratch files shipped.** `003_IMPLEMENT/01_out/.review-01_out.md` and `02_lock/.review-02_lock.md` are excluded by amended D009 and were the only reason `fest validate` scored 90. Fixed by `--exclude '.review-*'` **plus `--delete-excluded`** — the exclusion alone was not enough, because rsync protects an already-present excluded file at the destination from `--delete`; the first re-sync left both files in place and the score at 90. Now 100/100, and the copy is 86 files rather than 88.
- [x] **Critical 3 — `01_snapshot.md` recorded the wrong command.** That file is explicitly the command `004_REVIEW` re-runs before the PR, so leaving it would have reintroduced Criticals 1 and 2 at the pre-PR re-sync. Both copies of the command in the file (the requirement line and the recorded block) now carry the six exclusions and `--delete-excluded`, with a note on why the flag is load-bearing. The command is now idempotent from any destination state and self-corrects a copy made under an older exclusion list.

Suggestions taken:

- [x] **Dangling `CONTEXT.md` / `input_specs/` pointers** (~50 and ~90 citations across 28 files). Took the reviewer's own recommendation — one line near the top of `TODO.md` naming both as campaign-private and pointing at `002_PLAN/decisions/`, rather than 140 edits. No substance is lost: D001–D014 are all in the tree.
- [x] **`IMPLEMENTATION_PLAN.md` §05_snapshot row 01** restated the pre-amendment four-exclusion rsync as "per D009 exclusions", contradicting D009 two directories away. Now cites `002_PLAN/decisions/D009_public_snapshot.md` instead of restating it, so it cannot drift again.
- [x] **`measurements.md` §Baseline** — the probe (`scratchpad/measure/main.go`) and the `drop-sidecar` build worktree are campaign-side and not reproducible from this repo. Both rows now say so, and point at §Stream as the numbers a stranger can reproduce with the shipped binary.

Suggestions deliberately **not** taken (recorded, not silently dropped):

- [x] **`festivals/CA0001/` phrase-lock hits** — reviewer tightened my count to 25 lines / 11 files, 21 lines / 7 files of them the tracked `001_INGEST/input_specs/` pack. Committed, public, predating this branch, outside the `cans-v2` diff, and a call on another festival's shipped output. Still untouched; it is `004_REVIEW`'s / the operator's decision to scrub or accept on record.
- [x] **`internal/tts/worker.go` at 201 lines vs the 196 `FESTIVAL_RULES.md` pins** — inherited from `03_stream`, committed, no Go file was touched in this sequence. For `004_REVIEW` to accept or amend the rule.
- [x] **`fest.yaml` campaign-side `status_history` paths** — reviewer marked it machine metadata needing no action; `fest validate` passes. Agreed, no change.

## Iteration

For each finding:

1. Fix the issue
2. Re-run affected tests
3. Verify linting passes

### Verification after the fix

```
$ rg -i '<pattern 2 — CONTEXT.md §Professional grep>' <source, minus the six excluded paths>
(no output)                       # writing Results re-introduces hits; re-grepped after every edit

$ rsync -a --delete --delete-excluded \
    --exclude CONTEXT.md --exclude '001_INGEST/input_specs' \
    --exclude .fest --exclude .workflow --exclude .festival-checksums.json \
    --exclude '.review-*' \
    festivals/active/cans-v2-CV0001/ projects/worktrees/cans/cans-v2/festivals/CV0001/
exit=0

$ diff -r -x CONTEXT.md -x input_specs -x .fest -x .workflow \
       -x .festival-checksums.json -x '.review-*' <source> <copy>
exit=0                            # identical — Critical 1 cleared, nothing else was stale

$ find festivals/CV0001 \( -name '.review-*' -o -name '.workflow' \
       -o -name '.festival-checksums.json' -o -name CONTEXT.md -o -name input_specs \)
(no output)                       # Critical 2 cleared

$ fest validate festivals/CV0001
Score 100/100
VALIDATION PASSED                 # was 90/100

$ find festivals/CV0001 -type f | wc -l
      86                          # was 88

$ rg -i '<pattern 1>' README.md | rg -v '<allowed footer line>'   (no output)
$ rg -i '<pattern 2>' README.md docs/ tapes/ festivals/CV0001/    (no output)
$ rg -c 'fest.build' README.md                                    1

$ git check-ignore -v festivals/CV0001/fest.yaml                  (no output — not ignored)
$ git status --short
?? festivals/CV0001/              # the new tree and nothing else
```

No Go code was touched in this sequence, so tests and linting are unaffected by the iteration —
`03_testing`'s green suite still stands. Re-confirmed cheaply: `gofmt -l .` silent.

## Definition of Done

- [x] All critical findings fixed — all three, one root cause, verified above
- [x] All tests pass after changes — no Go touched; `03_testing`'s ten-package suite stands, `gofmt -l .` re-run silent
- [x] Linting passes
- [x] Code review findings addressed — 3 Criticals fixed, 3 Suggestions taken, 3 declined on the record with reasons
- [x] Ready to commit