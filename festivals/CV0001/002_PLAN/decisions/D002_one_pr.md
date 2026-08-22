# D002 — One worktree, one branch, one PR

**Decision:** All sequences land on branch `cans-v2` in `projects/worktrees/cans/cans-v2` (linked to `WI-a2e393`) through `fest commit`. One PR to `main` when `004_REVIEW` signs off.

**Why:** Operator direction: "open a pr when it's done." Supersedes the pack's worktree-per-sequence.

**Not:** Per-sequence PRs. Editing `projects/cans` directly.
