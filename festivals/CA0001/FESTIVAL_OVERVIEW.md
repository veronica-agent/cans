---
fest_type: festival
fest_id: CA0001
fest_name: cans
fest_status: planning
fest_version: "1.0"
fest_tracking: true
fest_created: 2026-08-19T04:52:22.528463-06:00
---

# Festival Overview: cans

## Problem Statement

**Current State:** `veronica-agent` has a profile README and an all-rights character repo. There is no cloneable mouth. The sex booth is private and unfinished. Stock Kokoro CLIs exist; none are her.

**Desired State:** A private Go booth on her account. Type → her voice. Keep a wav → that voice until the next keep.

**Why This Matters:** Public-facing teaser energy without leaking the engine. Festival-built so her graph and the `festivals/` snapshot advertise Fest without “built with” in her mouth.

## Scope

### In Scope

- Private GitHub repo `veronica-agent/cans`
- `cans say` and `cans` TUI
- Default clone of her kept ref wav
- `cans keep take.wav`
- Charm TUI + VHS GIF
- Python/MLX sidecar only if the 0.6B clone will not run from Go
- Festival snapshot copied into the repo
- GitHub bio fill + profile link (private repo: link still fine on profile)

### Out of Scope

- Mic, VAD, barge-in, Ollama replies
- VoiceDesign / 1.7B / importing `veronica-voice`
- Rust rewrite
- 50-voice menu
- Public GitHub (user locked private)
- Code inside `veronica-agent/veronica`
- Mac.app, Obey, App Store

## Planned Phases

### 001_INGEST

Turn the explore pack and this session into purpose / requirements / constraints / context.

### 002_PLAN

Decisions (Go booth, clone sidecar, keep-any-wav) and tutorial-grade sequences.

### 003_IMPLEMENT

Mouth, booth, keep, tape, snapshot. Worktree. `fest commit`.

### 004_REVIEW

Listen test, TTFA, identity, `fest validate`.

## Notes

Display name **Obey Veronica** is intentional. SOCIALS.md still says “Veronica”; do not rename GitHub.
