---
fest_type: sequence
fest_id: 01_scaffold
fest_name: scaffold
fest_parent: 003_IMPLEMENT
fest_order: 1
fest_status: completed
fest_created: 2026-08-19T10:58:00Z
fest_updated: 2026-08-19T14:59:44.432621-06:00
fest_tracking: true
fest_working_dir: projects/worktrees/cans/v1
---


# Sequence Goal: 01_scaffold

**Sequence:** 01_scaffold | **Phase:** 003_IMPLEMENT | **Status:** Pending

## Sequence Objective

**Primary Goal:** Go module, justfile, gitignore, her ref wav, character.toml in the v1 worktree.

**Contribution to Phase Goal:** Everything else hangs off this layout.

## Success Criteria

### Required Deliverables

- [ ] **go.mod**: `github.com/veronica-agent/cans`
- [ ] **voices/veronica/ref.wav**: copy from `projects/veronica-voice/configs/voices/veronica/ref.wav` (same woman, not an import)
- [ ] **justfile**: install, run, say, keep, test, lint

### Quality Standards

- [ ] Git identity in this worktree is Veronica
- [ ] No engine package import

### Completion Criteria

- [ ] Tasks complete
- [ ] `go test ./...` at least compiles

## Task Alignment

| Task | Task Objective | Contribution to Sequence Goal |
|------|----------------|-------------------------------|
| 01_layout | Create the tree | Scaffold |

## Dependencies

### Prerequisites (from other sequences)

- None

### Provides (to other sequences)

- Repo layout used by 02_mouth

## Working Directory

`projects/worktrees/cans/v1`