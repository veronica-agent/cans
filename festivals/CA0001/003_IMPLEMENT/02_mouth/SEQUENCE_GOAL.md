---
fest_type: sequence
fest_id: 02_mouth
fest_name: mouth
fest_parent: 003_IMPLEMENT
fest_order: 2
fest_status: completed
fest_created: 2026-08-19T10:58:00Z
fest_updated: 2026-08-19T16:31:49.861955-06:00
fest_tracking: true
fest_working_dir: projects/worktrees/cans/v1
---


# Sequence Goal: 02_mouth

**Primary Goal:** `cans say` clones the current ref wav and plays audio, printing TTFA.

Clone like `projects/veronica-voice/src/veronica/tts/qwen.py` `_clone` (lines 252–268): Base 0.6B, `ref_audio`, `ref_text`. Sidecar must not `import veronica`.