---
fest_type: sequence
fest_id: 01_out
fest_name: out
fest_parent: 003_IMPLEMENT
fest_order: 1
fest_status: completed
fest_created: 2026-08-21T05:04:55.63283-06:00
fest_updated: 2026-08-21T06:27:44.798239-06:00
fest_tracking: true
fest_working_dir: projects/worktrees/cans/cans-v2
---


# Sequence Goal: 01_out

**Primary Goal:** `cans say` writes a wav where it is told, reads one utterance from stdin, and can emit a JSON record — while `cans say "x"` stays byte-for-byte what it was at `1e8cea2`.

Covers P0-1, P0-2, P0-3, P0-4, P0-9, P0-10, P0-15, P0-16, P0-17, P0-18, P0-20. Decisions D004, D006, D010, D012.

Creates `internal/say` (the flow) and `cmd/cans/say_args.go` (the grammar); adds `tts.SayTo`. Nothing else moves. Exit codes in this sequence: 0, 1, 2.

Dependencies: none. Unblocks everything else.