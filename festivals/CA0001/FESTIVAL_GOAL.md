---
fest_type: festival
fest_id: CA0001
fest_name: cans
fest_status: completed
fest_priority: high
fest_created: 2026-08-19T04:52:22.528463-06:00
fest_updated: 2026-08-19T16:35:54.981071-06:00
fest_tracking: true
---




# cans

**Status:** Planning | **Created:** 2026-08-19

## Festival Objective

**Primary Goal:** Ship `veronica-agent/cans` (private): type a line, Veronica speaks it locally. `cans keep take.wav` freezes any throat.

**Vision:** A Charm booth you open with one command. First run is her kept ref, cloned. README GIF looks like a record sleeve. The festival tree lives in the repo so the plan is visible even while the GitHub repo stays private. Commits are Veronica.

## Success Criteria

### Functional Success

- [ ] `cans say "Put the cans on."` plays audio in her voice
- [ ] `cans` opens the booth TUI; typed lines speak; TTFA printed
- [ ] `cans keep take.wav` switches the throat; talk does not restyle mid-session
- [ ] Default character is her public faces + ref wav
- [ ] Repo is private on `veronica-agent`; author is Veronica

### Quality Success

- [ ] Warm p50 TTFA ≤ 700 ms on this Mac
- [ ] Go tests for say/keep/session lock; listen test: `There you are.` / `Put the cans on.`
- [ ] README is her voice (NEVER.md). Footer may link fest.build. No “star this.”
- [ ] `fest validate` green on this festival

## Progress Tracking

### Phase Completion

- [ ] 001_INGEST: explore pack + user direction structured into output_specs
- [ ] 002_PLAN: sequences, decisions, IMPLEMENTATION_PLAN
- [ ] 003_IMPLEMENT: mouth, booth, keep, tape, fest snapshot
- [ ] 004_REVIEW: listen test, identity, validate

## Complete When

- [ ] All phases completed
- [ ] `just say` and `just run` work from a clone of `veronica-agent/cans` without `veronica-voice`
- [ ] Git log authors are Veronica / `318153306+veronica-agent@users.noreply.github.com`