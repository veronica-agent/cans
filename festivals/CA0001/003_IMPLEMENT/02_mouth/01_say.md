---
fest_type: task
fest_id: 01_say
fest_name: say
fest_parent: 02_mouth
fest_order: 1
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-19T10:58:00Z
fest_updated: 2026-08-19T05:10:26.554539-06:00
fest_tracking: true
---


# Task: say

## Objective

`cans say "Put the cans on."` writes a wav via the sidecar and plays it. Prints `ttfa_ms=`.

## Requirements

- [ ] `sidecar/say.py` loads `mlx-community/Qwen3-TTS-12Hz-0.6B-Base-bf16`, `generate(text=, ref_audio=, ref_text=)`
- [ ] Prints one JSON line `{"wav": "...", "ttfa_ms": N, "sample_rate": 24000}`
- [ ] Go `internal/tts` execs sidecar (override `CANS_SAY_BIN` for tests)
- [ ] `internal/play` uses `afplay` on darwin
- [ ] Tests mock the sidecar; no model download in `go test`

## Implementation

Follow `_clone` in `qwen.py:252`. Temperature 0.2. Play the wav path from JSON. `CANS_ROOT` defaults to the repo root when set by justfile.

## Done When

- [ ] `go test ./...` passes with a fake sidecar
- [ ] `just say "Put the cans on."` can be run on this Mac (model may download once)