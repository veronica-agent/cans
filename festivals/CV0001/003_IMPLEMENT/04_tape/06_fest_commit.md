---
fest_type: gate
fest_id: 06_fest_commit.md
fest_name: Fest Commit Changes
fest_parent: 04_tape
fest_order: 6
fest_status: completed
fest_autonomy: high
fest_gate_id: fest-commit
fest_gate_type: commit
fest_managed: true
fest_created: 2026-08-21T05:04:57.555796-06:00
fest_updated: 2026-08-21T16:33:40.573977-06:00
fest_tracking: true
fest_version: "1.0"
---


# Gate: Commit Sequence Changes

Commit all changes from this sequence using the `fest commit` command.

## Pre-Commit Checklist

- [x] All tests pass
- [x] Linting is clean
- [x] No debug code or temporary files
- [x] No secrets or credentials in staged changes

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
- [x] No prohibited content in commit message

## Result

```
$ cd festivals/active/cans-v2-CV0001/003_IMPLEMENT/04_tape
$ gh auth switch -u veronica-agent
✓ Switched active account for github.com to veronica-agent
$ fest commit -m "feat: pipe tape and README scripting section" (+ What/Why body)
Hash         5e8123d
Task         FE-CV0001
Campaign     [veronica:ea389d71-FE-CV0001]
Root Commit  bbd6979
```

Project commit `5e8123d` on `cans-v2` in `projects/worktrees/cans/cans-v2`, campaign-root commit `bbd6979` for the festival files. Author `Veronica <318153306+veronica-agent@users.noreply.github.com>`. No Co-authored-by, no AI attribution.

Pushed **from the worktree**, which is what updates PR #14:

```
$ git push
To github-veronica-agent:veronica-agent/cans.git
   c443de9..5e8123d  cans-v2 -> cans-v2
```

Pre-commit: `gofmt -l .` empty, `go vet ./...` clean, `CANS_NOPLAY=1 go test -count=1 ./...` green in all ten packages, no debug code, no temp files (`lines.txt` / `out/` never land in the repo — the tape does `cd "$(mktemp -d)"`), no secrets. Six paths in the commit: `tapes/pipe.tape`, `docs/pipe.gif`, `.justfiles/vhs.just`, `README.md`, `cmd/cans/main.go`, `cmd/cans/say_args_test.go`. Worktree clean afterwards.

Not built and not run: another agent holds the real mouth, so this gate ran tests only — no `just build quick`, no `bin/cans`, no `just vhs`, no worker. `docs/pipe.gif` is committed as recorded (680×380, 123 s, `ttfa_ms` 31377/34570/35202) while `tapes/pipe.tape` now says `Set Height 300`; the two agree again after `004_REVIEW` re-cuts under load < 16.