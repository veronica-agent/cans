---
fest_type: phase
fest_id: 003_IMPLEMENT
fest_name: IMPLEMENT
fest_parent: cans-v2-CV0001
fest_order: 3
fest_status: completed
fest_created: 2026-08-21T05:03:56.848822-06:00
fest_updated: 2026-08-21T18:47:47.966647-06:00
fest_phase_type: implementation
fest_tracking: true
---


# Phase Goal: 003_IMPLEMENT

**Phase:** 003_IMPLEMENT | **Status:** Pending | **Type:** Implementation

## Phase Objective

**Primary Goal:** Build the slice on branch cans-v2: 01_out, 02_lock, 03_stream, 04_tape, 05_snapshot — in that order, through fest commit.

**Context:** Requirements are `001_INGEST/output_specs/requirements.md`; decisions are `002_PLAN/decisions/`; the per-sequence contract is `002_PLAN/plan/IMPLEMENTATION_PLAN.md`. Everything lands on branch `cans-v2` in `projects/worktrees/cans/cans-v2` through `fest commit`, and `004_REVIEW` opens the PR.

## Required Outcomes

Deliverables this phase must produce:

- [ ] `cans say` takes `-o`, stdin, `--json`, `--play`, `--nowait`, `--wait`, `--stream` — and nothing else; `cans say "x"` is unchanged
- [ ] exactly one `qwen3-tts-worker` resident under any loop, `xargs -P`, or second terminal; the booth holds the lock
- [ ] `--stream` renders N lines over one `Session`; Ctrl-C keeps finished wavs and exits 130
- [ ] `docs/pipe.gif` from `just vhs pipe`; README scripting section; professional-surface grep clean
- [ ] `festivals/CV0001/` snapshot in the repo; measurements recorded with their commands

<!-- Add more required outcomes as needed -->

## Quality Standards

Quality criteria for all work in this phase:

- [ ] `gofmt -l .` empty, `go vet ./...` clean, `CANS_NOPLAY=1 go test ./...` green on the fake worker
- [ ] no new `go.mod` requires; files < 500 lines; functions < 50 lines
- [ ] `context.Context` first on I/O; stdout is data, stderr is prose
- [ ] reviewer is a different agent than the implementer

<!-- Add more quality standards as needed -->

## Sequence Alignment

| Sequence | Goal | Key Deliverable |
|----------|------|-----------------|
| 01_out | `-o`, stdin, `--json`, exit codes; `internal/say` | `cans say -o take.wav "x"` |
| 02_lock | one mouth at a time | `internal/mouth`, `tts.OpenWith`, exit 75 |
| 03_stream | one `Session`, N lines | `cans say --stream -o 'out/%03d.wav'`, measurements |
| 04_tape | the honest demo | `docs/pipe.gif`, README §Scripting |
| 05_snapshot | the public plan | `festivals/CV0001/`, recheck |

<!-- Add rows as sequences are created -->

## Pre-Phase Checklist

Before starting implementation:

- [ ] Planning phase complete
- [ ] Architecture/design decisions documented
- [ ] Dependencies resolved
- [ ] Development environment ready

## Phase Progress

### Sequence Completion

- [ ] 01_out
- [ ] 02_lock
- [ ] 03_stream
- [ ] 04_tape
- [ ] 05_snapshot

<!-- Track sequence completion here -->

## Notes

Order is load-bearing: `01 → 02 → 03 → 04 → 05`. Do not start `03_stream` before `02_lock` is committed. Measurements need the real mouth and an idle machine (`uptime` 1-min load below the core count) — if the box is loaded, record that and come back.

---

*Implementation phases use numbered sequences. Create sequences with `fest create sequence`.*