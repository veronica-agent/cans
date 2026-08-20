---
fest_type: phase
fest_id: 003_IMPLEMENT
fest_name: IMPLEMENT
fest_parent: CA0001
fest_order: 3
fest_status: completed
fest_created: 2026-08-19T04:52:22Z
fest_updated: 2026-08-19T16:32:13.681405-06:00
fest_phase_type: implementation
fest_tracking: true
---


# Phase Goal: 003_IMPLEMENT

**Phase:** 003_IMPLEMENT | **Status:** Pending | **Type:** Implementation

## Phase Objective

**Primary Goal:** A working `cans` binary in the v1 worktree: say, keep, booth.

**Context:** Specs in 001_INGEST/output_specs. Decisions D001/D002. Execute in the cans worktree.

## Required Outcomes

- [ ] `just say "Put the cans on."` plays her clone
- [ ] `just keep` pins a wav
- [ ] `just run` opens the booth
- [ ] Tests pass without requiring the 0.6B model (mock sidecar)
- [ ] Commits are Veronica via `fest commit`

## Quality Standards

- [ ] No `veronica` Python package import
- [ ] NEVER.md on README
- [ ] Files under 500 lines

## Sequence Alignment

| Sequence | Goal | Key Deliverable |
|----------|------|-----------------|
| 01_scaffold | module + voices | go.mod, justfile, ref wav |
| 02_mouth | say + TTFA | sidecar + play |
| 03_keep | pin a wav | state/current.json |
| 04_booth | Charm TUI | internal/booth |
| 05_tape | GIF | tapes/booth.tape |
| 06_snapshot | fest ad | festivals/ copy, README |