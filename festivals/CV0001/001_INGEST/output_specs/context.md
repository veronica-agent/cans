# Context — cans v2 (CV0001)

## Where this sits

The design pack (`workflow/design/cans-v2`, `WI-a2e393`) was reviewed and accepted. This festival builds it. Base: `origin/main` at or after `1e8cea2`. Worktree `projects/worktrees/cans/cans-v2`, branch `cans-v2`, linked to `WI-a2e393`.

## Prior art — the dependencies that made this slice small

### `WI-7eb171` — the native mouth (PR #7/#8, main at `07ddf48`)

Replaced the Python/MLX sidecar with `qwen3-tts-worker` (C++/GGML/Metal). It shipped **three things the earlier draft of this pack listed as v2 work**:

| Previously designed as v2 work | Delivered by `WI-7eb171` |
|---|---|
| A JSONL request/response protocol over a long-lived child | `qwen3-tts-worker/v1` — `ready` handshake, `{"type":"synthesize","id","text","ref_wav"}`, `pcm_meta`/`final`, `{"type":"shutdown"}` (`internal/tts/worker.go`) |
| A warm-worker abstraction | `Session` — *"a warm worker for one booth (or one say)"* (`internal/tts/session.go:16`), `Open(ctx)` / `Say(ctx,text,throat)` / `Close()` |
| `ctx` through to the child | `SayWith(ctx,…)` (`synth.go:37`), `exec.CommandContext` (`worker.go:46`), `readLine` selecting on `ctx.Done()` (`worker.go:158-163`) |
| A streaming reader replacing buffered `lastJSONLine` | `bufio.Reader` line loop (`worker_pcm.go:13-52`) |

So **stream mode is no longer a protocol build** — it is open once, loop, close. What is genuinely new in v2 is the mouth lock plus the CLI surface. (`design-queue.md §What the merge already gave us`)

### `WI-8b1c5d` — drop the Python sidecar (PR #13, `1e8cea2`, merged 2026-08-21)

The shipped payload is Go + native worker only. `internal/ship/fs/` no longer unpacks `pyproject.toml`, `uv.lock`, `sidecar/say.py` for a path nothing calls. This closed finding #2 from `design-surface.md §Two findings` and is the base commit for every v2 worktree. (`seed.md §Dependencies — satisfied`)

## Key code locations the design cites

| Location | What it is / why v2 touches it |
|---|---|
| `internal/tts/synth.go:37-50` `SayWith` | Opens a worker and **defers `Close` per call** — the one-shot tax. Also the `CANS_SAY_BIN` fake-binary hook (`synth.go:41-43` → `synth_bin.go`). |
| `internal/tts/session.go:16` `Session`, `:54-60` | The warm-worker abstraction stream mode loops over. `Say` receives PCM and calls `audio.WritePCM16` to a **generated temp path** — `-o` passes the caller's path to that writer instead and skips `RemoveTemp`. Fewer moving parts than the pack's original sidecar-`--out` plan. |
| `internal/booth/booth.go:143-161` `Run` | Opens **one** `Session` in `Run` and reuses it for the whole TUI session (`defer sess.Close()`). The proof the warm pattern is cheap, and the place the lock is held for a human-length span (D001). |
| `internal/tts/worker.go` | `StartWorker` (`:46`), the `qwen3-tts-worker/v1` protocol, `readLine` ctx select. **196 lines — add files, do not grow it.** |
| `internal/tts/worker_pcm.go:13-52` | The PCM read loop. `wall` is stamped at `final` (`:42-47`), which is why `ttfa_ms` is really total synthesis time (P2-3). |
| `internal/tts/testdata/fakeworker/main.go` | The fake JSONL worker. `CANS_WORKER_BIN` points at it so CI runs the stream and lock paths without a real mouth. |
| `cmd/cans/main.go:56-78` (`say` case) | Today: argv joined, `doctor.Prepare`, `tts.Say`, print `ttfa_ms=N`, `play.File`, `tts.RemoveTemp`. This is the behavior P0-17 freezes. |
| `cmd/cans/main.go:109-140` `parseKeep` | **The interleaved-flag precedent.** It already accepts `keep take.wav -text words` and `keep -text words take.wav` by collecting positionals in a hand-rolled loop. `say` follows the same shape (D004) because stdlib `flag` stops at the first positional and the `design-pipes.md` loop examples put the text first. |
| `internal/play/play.go:12-30` `File` | `CANS_NOPLAY=1` skips playback **after** validating the wav header. It stays what it is: a test hook. `-o` is the supported way to be headless. |
| `internal/keep/keep.go` | The frozen throat (`keep.Load`) every path uses. `:42` holds the shipped default ref text flagged for the operator (P2-4). |
| `CANS_HOME` (default `~/.cans`) | `current/` (throat `ref.wav` + `current.json`), `native/bin/qwen3-tts-worker` (`CANS_WORKER_BIN` overrides), `native/models` (`CANS_WORKER_MODELS`), `shipped/`. **`mouth.lock` is new and lives here.** |

## The snapshot precedent — CA0001

The v1 festival tree is already committed in the public repo under `projects/cans/festivals/CA0001/` (`3c20108`) with the standard shape (`fest.yaml`, `FESTIVAL_GOAL.md`, `FESTIVAL_OVERVIEW.md`, `FESTIVAL_RULES.md`, `TODO.md`, `gates/`, `001_INGEST` … `004_REVIEW`). A stranger who clones can read how v1 was planned without installing `fest`. `05_snapshot` lands **CV0001** the same way, as the second readable plan. (`design-fest-ad.md §What v2 adds`)

## The fest-ad chrome rules

The repo advertises Festival through **chrome** — a topic, a readable `festivals/` tree, one footer line — never through a pitch in her mouth. (`design-fest-ad.md`)

| Lever | v2 |
|---|---|
| Tape | **Adds one**: a script piping lines in, wavs appearing on disk. Next to the booth GIF, not replacing it. No narration, no pitch, no claim about what she is. |
| `festivals/` | **Adds one plan** (CV0001). The ad compounds, one plan per slice. |
| Topic | Unchanged — `festival-methodology` plus `tts` / `local-ai`. |
| Footer | Unchanged — **Built with [Festival](https://fest.build)**, one line, shipped at `c72d4a0`. Not re-litigated, not duplicated. |
| Author | `veronica-agent`, via `fest commit` / `camp p commit` in the worktree. |

Verification: the professional-surface grep (`CONTEXT.md §Professional grep`, campaign-private) is empty over README, docs/, tapes/ and the snapshot; exactly one footer; `festivals/` holds two readable plans; `cans` still runs with `fest` uninstalled; `just vhs <name>` regenerates the new tape.

## Open questions carried in (both low, both for the operator — not this festival's call)

- The shipped default ref text `"Just like that, feel the rhythm of my voice."` (`internal/keep/keep.go:42`) predates the professional lock. The native mouth no longer sends it, but `keep` stores it and `doctor` can surface it. Changing the text without changing `ref.wav` would make it a false transcript — a voice-lock call.
- Whether `--wait` should have a non-infinite default for the booth. The festival ships infinite with the stderr line; revisit if it annoys.
