---
fest_type: task
fest_id: 01_tui
fest_name: tui
fest_parent: 04_booth
fest_order: 1
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-19T10:58:00Z
fest_updated: 2026-08-19T05:10:26.582373-06:00
fest_tracking: true
---


# Task: tui

## Objective

Full-screen booth. Dark. Quote from character.toml. Input field. Status idle/speaking. TTFA after a turn.

## Requirements

- [ ] `internal/booth` bubbletea + lipgloss + textinput
- [ ] Enter synthesizes via the same tts package as say
- [ ] Esc or ctrl+c quits
- [ ] Empty enter does nothing
- [ ] Synth error shows in the status line, does not crash

## Done When

- [ ] `just run` opens the booth
- [ ] `go test ./internal/booth` at least builds