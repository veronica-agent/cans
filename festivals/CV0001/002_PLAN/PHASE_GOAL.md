---
fest_type: phase
fest_id: 002_PLAN
fest_name: PLAN
fest_parent: cans-v2-CV0001
fest_order: 2
fest_status: completed
fest_created: 2026-08-21T04:32:49.135275-06:00
fest_updated: 2026-08-21T05:14:24.841899-06:00
fest_phase_type: planning
fest_tracking: true
---


# Phase Goal: 002_PLAN

**Phase:** 002_PLAN | **Status:** Pending | **Type:** Planning

## Phase Objective

**Primary Goal:** Plan architecture, design decisions, and task breakdown

**Context:** The pack left one contentious call open (does the booth hold the mouth lock?), asked for real measurements before any limit is sized, and the operator changed the delivery shape to one PR. Those have to be settled and the sequences scaffolded before a subagent can execute a task file.

## Exploration Topics

What areas need to be explored during this phase:

- Lock lifetime vs the booth (D001 in CONTEXT.md — confirm against `internal/booth/booth.go`)
- Where the lock lives in code so every `Session` is guarded by construction
- Flag parsing that tolerates interleaved flags and text (the `keep` precedent in `cmd/cans/main.go`)
- Real worker load time and resident memory on this Mac (`inputs/measurements.md`)
- What the fake worker can and cannot exercise (cancel, lock contention)

<!-- Add more exploration topics as identified -->

## Key Questions to Answer

Questions that must be answered before this phase is complete:

- Does the booth hold the lock for its whole run? (Yes — D001.)
- One PR or one per sequence? (One — D002.)
- Blank stdin lines in stream mode: error or skip? (Skip — D005.)
- Which sequence owns the 200-line stream-vs-loop measurement? (`03_stream`.)

<!-- Add more questions as they emerge -->

## Expected Documents

Documents that will be produced during this phase:

- `inputs/measurements.md` — worker load time, RSS, cold `cans say` wall, with commands
- `decisions/D001…D007` + `INDEX.md` — the decisions above, one file each
- `plan/STRUCTURE.md` — the tree
- `plan/IMPLEMENTATION_PLAN.md` — sequences, tasks, dependencies, verification per sequence
- `003_IMPLEMENT/` and `004_REVIEW/` scaffolded with task files and gates

<!-- Add more documents as planning progresses -->

## Success Criteria

This planning phase is complete when:

- [ ] Every decision has a file and a one-line rationale
- [ ] Measurements recorded with the command that produced them
- [ ] `003_IMPLEMENT` has five sequences, each with task files and four gates; `004_REVIEW` has the bar
- [ ] `fest validate` passes with zero markers

<!-- Add more success criteria as they become clear -->

## Notes

Planning is done by the orchestrating agent, which holds the full design context. The PRESENT checkpoint is approved on the operator's delegation and logged in CONTEXT.md.

---

*Planning phases use freeform structure. Create topic directories as needed.*