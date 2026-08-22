---
fest_type: sequence
fest_id: 03_stream
fest_name: stream
fest_parent: 003_IMPLEMENT
fest_order: 3
fest_status: completed
fest_created: 2026-08-21T05:04:55.99338-06:00
fest_updated: 2026-08-21T17:49:28.620777-06:00
fest_tracking: true
fest_working_dir: projects/worktrees/cans/cans-v2
---


# Sequence Goal: 03_stream

**Primary Goal:** `cans say --stream` speaks one utterance per stdin line over one warm `Session`, writes `-o 'out/%03d.wav'`, flushes a record per line, keeps going past a bad line, exits 130 on Ctrl-C with finished wavs intact — and the numbers that prove it are recorded.

Covers P0-5, P0-6, P0-7, P0-8, P0-9, P0-19, P0-20, P1-4. Decisions D005, D006, D007, D008, D013.

Dependencies: 02_lock committed (the stream holds the lock for its run).