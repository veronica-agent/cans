---
fest_type: gate
fest_id: 06_fest_commit.md
fest_name: Fest Commit Changes
fest_parent: 05_snapshot
fest_order: 6
fest_status: completed
fest_autonomy: high
fest_gate_id: fest-commit
fest_gate_type: commit
fest_managed: true
fest_created: 2026-08-21T05:04:57.567658-06:00
fest_updated: 2026-08-21T18:15:27.689584-06:00
fest_tracking: true
fest_version: "1.0"
---


# Gate: Commit Sequence Changes

Commit all changes from this sequence using the `fest commit` command.

## Pre-Commit Checklist

- [x] All tests pass — `03_testing`'s ten-package suite; no Go touched in this sequence or in `05_iterate`
- [x] Linting is clean — `gofmt -l .` silent, re-run after the iteration
- [x] No debug code or temporary files — the tree is 86 files, all `.md` plus `fest.yaml`; `.fest/`, `.workflow/`, `.festival-checksums.json` and the `.review-*` scratch notes are all excluded by D009
- [x] No secrets or credentials in staged changes — leak greps clean (nothing / nothing / `1`); home paths scrubbed to `~`; only the public `veronica-agent` commit identity appears

## Commit Command

You **MUST** use `fest commit` — not `git commit`. The `fest commit` command tags
commits with task reference IDs for tracking and metrics.

```bash
fest commit -m "<type>: <summary>"
```

**CRITICAL:** Do NOT use `git commit`, `git add && git commit`, or any other git
commit workflow. Always use `fest commit` so task references are preserved.

## Commit Message Format

```
<type>: <concise summary of changes>

<what changed — list concrete modifications>

<why it changed — purpose and motivation>
```

**Types:** `feat`, `fix`, `refactor`, `test`, `docs`, `chore`

The message should describe WHAT changed and WHY. Be specific about files,
functions, or features that were added, modified, or removed.

## Ethical Requirements

The following practices are **prohibited** in commit messages:

- NO "Co-authored-by" tags for AI assistants
- NO AI tool attribution or advertisements
- NO links to AI services or products

## Definition of Done

- [x] Pre-commit checklist verified
- [x] Commit created with `fest commit` (not `git commit`)
- [x] Message describes what changed and why
- [x] No prohibited content in commit message — no AI attribution, no co-author tag, no tool links

## Result

```
$ cd <festival>/003_IMPLEMENT/05_snapshot
$ gh auth switch -u veronica-agent
✓ Switched active account for github.com to veronica-agent
$ fest commit -m "feat: festivals/CV0001 — the v2 plan as a readable tree …"
Hash        9136cbe
Message     [veronica:ea389d71-FE-CV0001] feat: festivals/CV0001 — the v2 plan as a readable tree
Task        FE-CV0001
Campaign    [veronica:ea389d71-FE-CV0001]
Root Commit 21e5252
```

- **Worktree (`veronica-agent/cans`, branch `cans-v2`): `9136cbe`** — `86 files changed, 5541 insertions(+)`, all additions under `festivals/CV0001/`. Author `Veronica <318153306+veronica-agent@users.noreply.github.com>`, the campaign identity the rules require.
- **Campaign root: `21e5252`** — the festival's own task documents and statuses.
- Pushed from the worktree: `5e8123d..9136cbe  cans-v2 -> cans-v2`. `git status --short` clean afterwards.

The copy in the commit is current as of `05_iterate` being completed; only this gate's own
`fest_status` moved after the snapshot was taken. `004_REVIEW` re-runs the recorded rsync
(`01_snapshot.md` — six exclusions plus `--delete-excluded`) before the PR, which picks it up.