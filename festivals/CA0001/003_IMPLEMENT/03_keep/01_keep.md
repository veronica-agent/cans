---
fest_type: task
fest_id: 01_keep
fest_name: keep
fest_parent: 03_keep
fest_order: 1
fest_status: completed
fest_autonomy: high
fest_created: 2026-08-19T10:58:00Z
fest_updated: 2026-08-19T05:10:26.568148-06:00
fest_tracking: true
---


# Task: keep

## Objective

Pin a wav as the current woman.

## Requirements

- [ ] `cans keep path.wav` copies into `$CANS_HOME/current/` (default `~/.cans`) and writes `current.json` `{wav, ref_text}`
- [ ] `--text` sets ref_text; default empty
- [ ] `cans say` uses current if present, else `voices/veronica`
- [ ] Missing wav → exit 2, stderr explains
- [ ] Tests cover missing file and default fallback

## Done When

- [ ] `go test ./...` covers keep
- [ ] After keep, say uses the new wav path