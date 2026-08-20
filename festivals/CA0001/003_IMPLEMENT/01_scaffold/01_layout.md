---
fest_type: task
fest_id: 01_layout
fest_name: layout
fest_parent: 01_scaffold
fest_order: 1
fest_status: completed
fest_autonomy: high
fest_created: 2026-08-19T10:58:00Z
fest_updated: 2026-08-19T05:10:26.540837-06:00
fest_tracking: true
---


# Task: layout

## Objective

Create the cans worktree layout so later sequences have a module to compile.

## Requirements

- [ ] `go.mod` module `github.com/veronica-agent/cans`
- [ ] `.gitignore` for bin/, .venv/, state/, *.wav except voices/veronica/ref.wav
- [ ] Copy `projects/veronica-voice/configs/voices/veronica/ref.wav` and `meta.json` into `voices/veronica/`
- [ ] `character.toml` with name Veronica and pull-quote `Put the cans on.`
- [ ] `justfile` stubs: install, run, say, keep, test, lint
- [ ] Worktree git user is Veronica

## Implementation

Work in `projects/worktrees/cans/v1`. Do not edit `projects/veronica-voice` except reading the ref wav.

1. `go mod init github.com/veronica-agent/cans`
2. Copy wav+meta with `cp` (binary). meta.json `ref_text` must stay `Just like that, feel the rhythm of my voice.`
3. `character.toml`:
   ```
   name = "Veronica"
   quote = "Put the cans on."
   voice = "voices/veronica"
   ```
4. justfile recipes call `go run ./cmd/cans` once that package exists; for this task a `cmd/cans/main.go` that prints usage is enough.

## Done When

- [ ] All requirements met
- [ ] `ls voices/veronica/ref.wav` exists in the worktree