# Festival Overview: cans-v2

## Problem Statement

**Current State:** `cans say` takes text from argv only, plays the wav and deletes it. There is no way to write a file, no stdin, and nothing coordinating concurrent calls. `tts.SayWith` opens a worker and defers `Close` per call (`internal/tts/synth.go`), and every `Open` loads the GGUF model — so a 200-line loop pays 200 model loads, and `xargs -P 8` puts eight workers' weights in unified memory with nothing refusing and nothing waiting. The booth already runs warm on one `Session`; the CLI does not.

**Desired State:** A script walks a document and hands cans one line at a time — from argv, from stdin, or as a stream — and gets a wav where it asked, a machine-readable record on stdout, and exactly one worker resident no matter how the loop is shaped.

**Why This Matters:** This is what makes cans useful beyond the booth, it is the honest demo for the second tape, and the festival tree that builds it becomes the second readable plan in the public repo.

## Scope

### In Scope

- `-o path` (one-shot) and `-o 'out/%03d.wav'` (stream) — wav written in Go by `audio.WritePCM16` to the caller's path; no `RemoveTemp`
- stdin as one utterance (`echo x | cans say`, `cans say -`)
- `--stream`: one utterance per stdin line over one `Session`, records flushed per line, per-line failure policy
- `--json` records on stdout; stdout is data, stderr is prose
- `--play` with `-o` (write and play)
- Exit codes 0 / 1 / 2 / 75
- Mouth lock: `flock` on `CANS_HOME/mouth.lock` held for the lifetime of one `Session`; `--nowait`, `--wait <dur>`; the booth takes it too
- Tests on the fake worker; measurements (load time, RSS, stream vs loop) recorded here
- Second VHS tape (`just vhs pipe`) and the README scripting section
- Snapshot of this festival into `projects/cans/festivals/CV0001/`

### Out of Scope

- Document ingestion: `cans read`, `-f`, markdown stripping, sentence chunking — the script owns the document
- A daemon (`cansd`), a job queue file, FIFO fairness
- Voice picker, SSML, config file, mic/VAD/replies, browser booth, radio
- Restyling mid-session; keep stays the only throat change
- A mid-utterance abort in the worker protocol (Ctrl-C terminates the worker instead — D014; documented)
- The `ttfa_ms` semantics bug and the shipped default ref text — flagged for the operator, not this festival's call
- Re-cutting `docs/booth.gif` / `booth.mp4`

## Planned Phases

### 001_INGEST

Structure the accepted design pack (`workflow/design/cans-v2`, `WI-a2e393`) and the operator's direction into purpose / requirements / constraints / context.

### 002_PLAN

Decide the contentious calls (lock lifetime and the booth, one PR, flag grammar), take the measurements, write STRUCTURE and IMPLEMENTATION_PLAN, scaffold the implementation sequences and task files, apply gates.

### 003_IMPLEMENT

`01_out` → `02_lock` → `03_stream` → `04_tape` → `05_snapshot`, in that order, on branch `cans-v2`. Safe before fast: the lock lands before stream mode so a stalled festival still leaves a safe tool.

### 004_REVIEW

The ship bar from `design-recommend.md`, identity, professional-surface greps, `fest validate`, and the PR.

## Notes

- Dependencies are satisfied: `WI-7eb171` (native mouth) gave the JSONL protocol, `Session`, and ctx; `WI-8b1c5d` (PR #13, `1e8cea2`) removed the Python payload.
- Execution is delegated to Opus/Sonnet subagents; the orchestrating agent plans, reviews, and approves on the operator's delegation. Every approval is logged in CONTEXT.md.
- The design pack said worktree-per-sequence; the operator asked for one PR. One worktree, one branch, one PR.
