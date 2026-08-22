---
fest_type: phase
fest_id: 001_INGEST
fest_name: INGEST
fest_parent: cans-v2-CV0001
fest_order: 1
fest_status: completed
fest_created: 2026-08-21T04:32:49.132577-06:00
fest_updated: 2026-08-21T05:00:07.75662-06:00
fest_phase_type: ingest
fest_tracking: true
---


# Phase Goal: 001_INGEST

**Phase:** 001_INGEST | **Status:** Pending | **Type:** Ingest

## Phase Objective

**Primary Goal:** Ingest and structure input materials into actionable specifications

**Context:** Seeded input is available in input_specs/seed.md and should be transformed into structured output specs for planning.

## Input Sources

Place all raw input materials in `input_specs/`:

- [x] `seed.md` — the slice in one line, v2-is / v2-is-not, satisfied dependencies
- [x] `design-surface.md`, `design-pipes.md`, `design-queue.md`, `design-fest-ad.md`, `design-recommend.md` — the accepted design pack (`WI-a2e393`), copied verbatim
- [x] `user-direction.md` — the operator's verbatim direction and the rules that do not move

## Expected Outputs

The following structured documents will be created in `output_specs/`:

| Output | Purpose |
|--------|---------|
| `purpose.md` | Festival purpose, success criteria, motivation |
| `requirements.md` | Prioritized requirements (P0/P1/P2) with traceability |
| `constraints.md` | Technical and process constraints |
| `context.md` | Prior art, related systems, key references |

## Success Criteria

This ingest phase is complete when:

- [ ] All input sources reviewed and understood
- [ ] Output specs created following standard structure
- [ ] User has approved the structured output
- [ ] No unresolved questions or ambiguities

## Workflow

This phase uses step-based workflow guidance. See `WORKFLOW.md` for the step-by-step process.

Use `fest next` to see the current step.
Use `fest workflow advance` to move to the next step.

## Notes

The design pack was already reviewed and accepted by the operator, so ingest is restructuring, not discovery. Do not re-open decisions the pack closed (no ingestion, no daemon). The PRESENT checkpoint is approved by the orchestrator on the operator's delegation — log it in CONTEXT.md.

---

*Ingest phases transform unstructured input into structured specifications ready for planning.*