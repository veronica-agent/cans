---
fest_autonomy: high
fest_created: 2026-08-21T05:04:57.567658-06:00
fest_gate_id: fest-commit
fest_gate_type: commit
fest_id: 06_fest_commit.md
fest_managed: true
fest_name: Fest Commit Changes
fest_order: 6
fest_parent: 05_snapshot
fest_status: pending
fest_tracking: true
fest_type: gate
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

- [ ] Pre-commit checklist verified
- [ ] Commit created with `fest commit` (not `git commit`)
- [ ] Message describes what changed and why
- [ ] No prohibited content in commit message