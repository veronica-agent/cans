---
fest_type: sequence
fest_id: 02_lock
fest_name: lock
fest_parent: 003_IMPLEMENT
fest_order: 2
fest_status: completed
fest_created: 2026-08-21T05:04:55.867248-06:00
fest_updated: 2026-08-21T06:55:34.626135-06:00
fest_tracking: true
fest_working_dir: projects/worktrees/cans/cans-v2
---


# Sequence Goal: 02_lock

**Primary Goal:** Exactly one `qwen3-tts-worker` is resident, ever. A `flock` on `CANS_HOME/mouth.lock` is taken before the worker starts and released after it exits; the booth holds it for its whole run; `--nowait` exits 75 and `--wait` bounds the block.

Covers P0-11, P0-12, P0-13, P0-14, P0-15 (75), P0-19 (`kill -9`), P0-20. Decisions D001, D003, D011.

Creates `internal/mouth`; adds `tts.OpenWith`; the booth uses it. This is the only genuinely new machinery in v2.

Dependencies: 01_out (the `say.Options` the flags land in). Must be committed before 03_stream starts — safe before fast.