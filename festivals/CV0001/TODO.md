# Festival TODO - cans-v2

**Goal**: Cans becomes a unix primitive a script can drive over a document — text in from argv/stdin, wav out where you point it, one mouth at a time.
**Status**: Planning

Campaign-private and not in this public copy: `CONTEXT.md` (session memory) and `001_INGEST/input_specs/` (the design pack). Every decision they reference is under `002_PLAN/decisions/`.

---

## Festival Progress Overview

### Phase Completion Status

- [ ] 001_INGEST — design pack + direction → output_specs
- [ ] 002_PLAN — decisions, measurements, structure, scaffold, gates
- [ ] 003_IMPLEMENT — 01_out, 02_lock, 03_stream, 04_tape, 05_snapshot
- [ ] 004_REVIEW — ship bar, identity, greps, PR

### Current Work Status

```
Active Phase: 001_INGEST
Active Sequences: N/A (workflow phase)
Blockers: None
```

---

## Phase Progress

### 003_IMPLEMENT

**Status**: Not Started

#### Sequences

- [ ] 01_out — `-o`, stdin, `--json`, exit codes
- [ ] 02_lock — mouth lock, `--nowait` / `--wait`, booth
- [ ] 03_stream — `--stream`, `%03d` template, cancel, measurements
- [ ] 04_tape — pipe tape + README scripting section
- [ ] 05_snapshot — festival snapshot, professional recheck

---

## Blockers

None currently.

---

## Decision Log

See `CONTEXT.md` and `002_PLAN/decisions/`.

---

*Detailed progress available via `fest progress`*
