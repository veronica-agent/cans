---
fest_type: task
fest_id: 01_snapshot.md
fest_name: snapshot
fest_parent: 05_snapshot
fest_order: 1
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:57.276998-06:00
fest_updated: 2026-08-21T17:55:14.197889-06:00
fest_tracking: true
---


# Task: snapshot

## Objective

Copy this festival into the repo as `festivals/CV0001/` — the second readable plan — with D009's exclusions, and prove the copy is clean.

## Requirements

- [x] Source: this festival's directory in the campaign (`festivals/active/cans-v2-CV0001` by the time this runs — confirm with `pwd` from `fest next`). Destination: `<worktree>/festivals/CV0001/`.
- [x] Command (record the exact one you ran): `rsync -a --delete --delete-excluded --exclude CONTEXT.md --exclude '001_INGEST/input_specs' --exclude .fest --exclude .workflow --exclude .festival-checksums.json --exclude '.review-*' <festival>/ <worktree>/festivals/CV0001/`. `--exclude '.review-*'` and `--delete-excluded` were added in `05_iterate`: D009 was amended in `004_REVIEW` to exclude the reviewers' hidden scratch notes, and plain `--exclude` will not remove a copy that is already there (rsync protects excluded files at the destination from `--delete`).
- [x] Run the professional-surface grep from `CONTEXT.md §Professional grep` over `festivals/CV0001/` — it must print nothing. If it does, the offending text is fixed in the **campaign festival** (the source), not in the copy, and the rsync is re-run.
- [x] `fest validate festivals/CV0001` passes from the worktree (the snapshot carries `fest.yaml`).
- [x] No README change: the tree is discoverable next to `festivals/CA0001/`, and the footer already points at Festival. Do not add chrome.
- [x] `git check-ignore -v festivals/CV0001/fest.yaml` prints nothing (the tree is not ignored).

## Implementation

1. Run the rsync; `find festivals/CV0001 -type f | wc -l` and eyeball the tree against `002_PLAN/plan/STRUCTURE.md`.
2. Run the grep and `fest validate`.
3. `004_REVIEW` re-runs the same rsync once its own statuses are set, so record the exact command in this task file for reuse.

## Results

### Step 0 — scrub the source before copying (not in the original task text)

The public tree must pass grep #2, whose pattern includes the operator's own name, and the festival's
recorded results carried absolute home paths. Scrubbed **in the campaign source**, per the
task's rule that fixes happen in the source and the rsync is re-run.

```bash
$ cd ~/Dev/AI/veronica-campaign/festivals/active/cans-v2-CV0001
$ find . -name '*.md' -not -path './.fest/*' -not -path './001_INGEST/input_specs/*' \
      -not -name 'CONTEXT.md' -print0 | xargs -0 sed -i '' 's|/Users/<user>|~|g'
# '<user>' above stands for the literal account name in the absolute home prefix on this Mac;
# the real sed had it spelled out. It is written this way here so this file passes grep #2.
```

Home-path hits in public-bound `.md` (everything but `CONTEXT.md` and `001_INGEST/input_specs/`):
**6 matching lines / 3 files before → 0 after** — `002_PLAN/inputs/measurements.md` 4,
`003_IMPLEMENT/03_stream/05_testing.md` 1, `003_IMPLEMENT/04_tape/03_testing.md` 1 (that last
line held two occurrences). Two further lines in the same files matched only the pattern below
and were handled by rewording.

Grep #2 over the source minus `CONTEXT.md`, `001_INGEST/input_specs/`, `.fest/`:
**13 hits / 7 files before → 7 hits / 6 files after the sed → 0 after rewording.**
Nothing was deleted; each hit was neutralised in place:

| File | Was | Now |
|------|-----|-----|
| `FESTIVAL_RULES.md` | the remote's display name spelled out, first word matching the pattern | `Display name stays the one the campaign identity lock already sets — do not change it.` |
| `001_INGEST/output_specs/requirements.md` | a `design-surface.md` section heading citation naming the operator | `design-surface.md §Two findings for the operator` |
| `003_IMPLEMENT/03_stream/05_testing.md` | `git describe`'s modified-tree suffix on the `ship.Version` string in the `just build quick` line | `ship.Version=v0.1.0-25-gc736f46-<uncommitted>` |
| `003_IMPLEMENT/04_tape/02_readme_scripting.md`, `03_testing.md`, `05_iterate.md` | the two grep patterns quoted verbatim in the recorded commands | `<pattern 1 — CONTEXT.md §Professional grep>` / `<pattern 2 — …>` — the patterns are campaign-private and must not ship |
| `003_IMPLEMENT/04_tape/03_testing.md` | the account-name owner column of an `ls -l docs/pipe.gif` line | `-rw-r--r--@ 1 user  staff  168370 …` |

Every result (counts, exit codes, timings) is unchanged; only the wording is.

### The rsync (exact command — `004_REVIEW` re-runs this one)

```bash
$ cd ~/Dev/AI/veronica-campaign
$ rsync -a --delete --delete-excluded \
    --exclude CONTEXT.md --exclude '001_INGEST/input_specs' \
    --exclude .fest --exclude .workflow --exclude .festival-checksums.json \
    --exclude '.review-*' \
    festivals/active/cans-v2-CV0001/ projects/worktrees/cans/cans-v2/festivals/CV0001/
exit=0

# Six exclusions, matching D009 as amended in `004_REVIEW`. `--delete-excluded` is load-bearing,
# not decoration: rsync protects excluded files that already exist at the destination from
# `--delete`, so adding `--exclude '.review-*'` on its own left both scratch files in the copy
# and `fest validate` still at 90. With it, the command is idempotent from any destination state
# and self-corrects a copy made under an older exclusion list. This is the command to re-run
# before the PR.
```

### Proof

```
$ find festivals/CV0001 -type f | wc -l
      86
# 88 before `05_iterate`; the two `.review-*` scratch files are now excluded.

$ find festivals/CV0001 -maxdepth 2 | sort        # matches 002_PLAN/plan/STRUCTURE.md
festivals/CV0001/001_INGEST/{GATES.md,output_specs,PHASE_GOAL.md,WORKFLOW.md}   # no input_specs
festivals/CV0001/002_PLAN/{decisions,GATES.md,inputs,PHASE_GOAL.md,plan,WORKFLOW.md}
festivals/CV0001/003_IMPLEMENT/{01_out,02_lock,03_stream,04_tape,05_snapshot,GATES.md,PHASE_GOAL.md}
festivals/CV0001/004_REVIEW/{BAR.md,GATES.md,PHASE_GOAL.md}
festivals/CV0001/{fest.yaml,FESTIVAL_GOAL.md,FESTIVAL_OVERVIEW.md,FESTIVAL_RULES.md,gates/implementation,TODO.md}
# CONTEXT.md, 001_INGEST/input_specs/, .fest/, .festival-checksums.json all absent

$ rg -i '<pattern 1 — CONTEXT.md §Professional grep>' README.md | rg -v '<allowed footer line>'
(no output)
$ rg -i '<pattern 2 — CONTEXT.md §Professional grep>' README.md docs/ tapes/ festivals/CV0001/
(no output)
$ rg -c 'fest.build' README.md
1

$ fest validate festivals/CV0001
✓ STRUCTURE / ✓ COMPLETENESS / ✓ Task Files / ✓ QUALITY GATES / ✓ Markers
Score 100/100
# 100 after `--exclude '.review-*'` landed. Before it, the two hidden reviewer scratch files
# (`.review-01_out.md`, `.review-02_lock.md`) shipped and cost 10 points on filename shape.

$ git check-ignore -v festivals/CV0001/fest.yaml
(no output — exit 1, the tree is not ignored)

$ git status --short
?? festivals/CV0001/
```

README untouched; no chrome added.

## Done when

- [x] `festivals/CV0001/` present, grep empty, `fest validate` green, `git status` shows only the new tree