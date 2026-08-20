---
fest_type: phase
fest_id: 001_INGEST
fest_name: INGEST
fest_parent: cans-CA0001
fest_order: 1
fest_status: completed
fest_created: 2026-08-19T04:52:22.535895-06:00
fest_updated: 2026-08-19T16:29:06.815139-06:00
fest_phase_type: ingest
fest_tracking: true
---


# Phase Goal: 001_INGEST

**Phase:** 001_INGEST | **Status:** Pending | **Type:** Ingest

## Phase Objective

**Primary Goal:** Turn the public-oss explore pack and the 2026-08-19 user direction into structured specs for planning.

**Context:** WI-9608d7 (`workflow/explore/public-oss`) already ranked `cans`. This session locked: her voice, keep-any-wav, Go booth, private repo on `veronica-agent`.

## Input Sources

- [x] `input_specs/user-direction.md` — this session
- [x] `input_specs/explore-*.md` — WI-9608d7 pack
- [x] `input_specs/docs-voice.md` / `docs-character.md` / `docs-never.md` / `docs-socials.md` / `docs-tone.md` / `docs-github.md`
- [x] `input_specs/seed.md` — festival create seed (recommend)

## Expected Outputs

| Output | Purpose |
|--------|---------|
| `purpose.md` | Why cans exists and what done means |
| `requirements.md` | P0/P1/P2 |
| `constraints.md` | Identity, stack, never-list |
| `context.md` | Prior art and this campaign |
| `PRESENTATION.md` | Summary for the ingest gate |

## Success Criteria

- [ ] All input sources reviewed
- [ ] Output specs created
- [ ] User direction captured (private repo, her voice, Go, keep)
- [ ] No leftover “should we use Kokoro Bella / Rust” questions