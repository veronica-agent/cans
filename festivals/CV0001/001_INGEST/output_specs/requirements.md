# Requirements — cans v2 (CV0001)

Every requirement cites the input it came from. `CONTEXT.md` decisions D001–D007 are already made and are not re-opened here.

## P0 — the slice. Without these there is no v2.

| # | Requirement | Source |
|---|-------------|--------|
| P0-1 | `-o take.wav` writes the wav to the caller's path, does **not** play, does **not** delete. `Session.Say` already writes in Go via `audio.WritePCM16`; `-o` passes the caller's path to that writer and skips `RemoveTemp`. | `design-pipes.md §Audio out`, `design-recommend.md §v2 is` |
| P0-2 | `--play` with `-o` writes **and** plays, for a human watching a script run. | `design-pipes.md §Audio out` |
| P0-3 | stdin is one utterance when argv is empty: `echo hi \| cans say`, and `cans say -` explicitly. Whole input is one utterance. | `design-pipes.md §Text in` |
| P0-4 | Empty argv **and** stdin is a TTY → the usage error it is today (exit 2), not a hang. | `design-pipes.md §Text in` |
| P0-5 | `--stream` changes stdin from one utterance to one utterance per line, spoken over **one** `Session`: open once, loop `Say`, close at EOF. | `design-pipes.md §Text in`, `design-queue.md §Stream mode` |
| P0-6 | `-o 'out/%03d.wav'` in stream mode: one file per line, index substituted. A bad template is exit 2. | `design-pipes.md §Audio out`, `§Streams and exit codes` |
| P0-7 | Stream per-line policy: blank stdin lines are **skipped and do not consume an index**; any other per-line failure is reported on stderr (and as `{"line":N,"error":"…"}` under `--json`), the stream continues, and the exit code at EOF is 1 if any line failed. | `CONTEXT.md D005`, `design-pipes.md §Streams and exit codes` |
| P0-8 | `--stream` without `-o` plays each line through the speakers in order. | `CONTEXT.md D007` |
| P0-9 | `--json` emits JSONL on stdout, **flushed as each utterance finishes**. One-shot record: `{"wav":…,"ttfa_ms":…,"sample_rate":…}` (the existing `tts.Result` tags). Stream adds `"line":N`, 1-based on the stdin line. | `CONTEXT.md D006`, `design-pipes.md §Streams and exit codes` |
| P0-10 | `-o` without `--json` prints the wav path on stdout. No `-o`, no `--json`: `ttfa_ms=N` exactly as today. | `CONTEXT.md D006` |
| P0-11 | Mouth lock: `flock(2)` on `CANS_HOME/mouth.lock`, held for the lifetime of one `Session`. Any loop, any `xargs -P`, any second terminal: exactly one `qwen3-tts-worker` resident. | `design-queue.md §Lock mechanics`, `CONTEXT.md D003` |
| P0-12 | A blocked caller writes one line to **stderr** (`waiting for the mouth…`) so a script that looks hung explains itself. Never stdout. | `design-queue.md §Lock mechanics` |
| P0-13 | `--nowait` exits 75 immediately; `--wait <dur>` bounds the block; default waits forever, because for a document script blocking *is* correct. | `design-queue.md §Lock mechanics`, `§Limits` |
| P0-14 | The booth takes the lock and holds it for its whole run. A script started alongside it waits (stderr line) or exits 75 with `--nowait`; a booth started while a script holds the mouth waits the same way before the TUI opens. | `CONTEXT.md D001`, `design-queue.md §Lock mechanics` |
| P0-15 | Exit codes: `0` spoke it · `1` runtime failure (worker died or missing, bad wav, disk) · `2` usage (no text, unknown flag, bad `-o` template) · `75` busy with `--nowait` (`EX_TEMPFAIL`, so `xargs` and retry loops read it correctly). | `design-pipes.md §Streams and exit codes` |
| P0-16 | **stdout is data, stderr is prose. They never mix.** stdout carries only the v1 `ttfa_ms=` line, wav paths, or JSONL. Progress, waiting, and errors go to stderr. | `design-pipes.md §Streams and exit codes`, `FESTIVAL_RULES.md` |
| P0-17 | **`cans say "x"` behaves exactly as at `1e8cea2`.** One-shot stays one-shot: prints `ttfa_ms=N`, plays, deletes the temp wav. Every existing test passes unchanged. | `seed.md`, `user-direction.md §Standing rules`, `design-queue.md §Stream mode` constraint 1 |
| P0-18 | Flags interleave with text in **both** orders, the way `keep` already parses (`cmd/cans/main.go:109` `parseKeep`): `cans say "$line" -o out.wav` and `cans say -o out.wav "$line"` both work. The flag set is exactly `-o`, `--json`, `--stream`, `--play`, `--nowait`, `--wait <dur>`, and `-`. Nothing else. | `CONTEXT.md D004` |
| P0-19 | Ctrl-C mid-stream: stop reading stdin, `Close` the session (graceful `shutdown` if the worker is idle; SIGTERM then SIGKILL if it is mid-utterance — D014), leave completed wavs in place, release the lock, report the line it stopped on. `kill -9` leaves the next run unblocked — the kernel drops the `flock`. | `design-queue.md §Cancellation and partial output`, `§Lock mechanics` |
| P0-20 | The stream and lock paths are fakeable: tests run against `internal/tts/testdata/fakeworker` via `CANS_WORKER_BIN`, with `CANS_NOPLAY=1`, so `go test ./...` and CI need no real mouth. | `design-queue.md §Stream mode` constraint 2, `FESTIVAL_RULES.md §Code` |

## P1 — the slice is not delivered without these, but the tool works.

| # | Requirement | Source |
|---|-------------|--------|
| P1-1 | A second VHS tape (`just vhs pipe`) showing a script piping lines in and wavs appearing in `out/`. Sits next to the booth GIF; does not replace it. Reproducible: `just vhs <name>` regenerates it. | `design-fest-ad.md §What v2 adds`, `§Verification for the ad` #5 |
| P1-2 | README scripting section: the loops from `design-pipes.md §Loops this makes possible`, the exit-code table, and a plain statement of what Ctrl-C does — per D014: the line being spoken is dropped, finished wavs stay, exit 130, a second Ctrl-C stops at once. | `design-pipes.md`, `design-queue.md §Stream mode` constraint 5, `CONTEXT.md §Deferred` |
| P1-3 | Snapshot of this festival into `projects/cans/festivals/CV0001/` — the second readable plan in the public tree; `cans` still runs with `fest` uninstalled. | `design-fest-ad.md §What v2 adds`, `§Verification for the ad` #4 |
| P1-4 | Measurements recorded in the festival with the command that produced each: real resident memory (RSS) for one `qwen3-tts-worker`, real GGUF load time, 200-line stream vs the equivalent 200-call loop, and `pgrep` samples under `xargs -P 8` over 50 lines. Every limit and every README claim quotes these numbers, not the pack's estimates. | `design-recommend.md §Festival shape` ("Measure before you size"), `design-queue.md §The bar`, `FESTIVAL_RULES.md §Process` |
| P1-5 | Professional-surface grep (pattern in `CONTEXT.md §Professional grep`, campaign-private) is empty over README, docs/, tapes/ and the festival snapshot; exactly one Festival footer, unchanged from `c72d4a0`. Example text in README, tape, and fixtures is boring and technical. | `design-fest-ad.md §The professional lock`, `§Verification for the ad` |

## P2 — parked. **Do not build these.**

| # | Item | Why parked | Source |
|---|------|-----------|--------|
| P2-1 | `cansd` daemon — warmth *across* invocations | Real lifecycle, real version skew, and not what makes v2 useful. `Session` is already the client it would wrap, so it stays cheap to add later. Revisit only if a real user asks. | `design-queue.md §Options`, `CONTEXT.md §Deferred` |
| P2-2 | Mid-utterance abort in the worker protocol | Needs worker support. Ctrl-C terminates a busy worker instead (D014); documented instead of pretended away (P1-2). | `design-queue.md §Stream mode` constraint 5, `CONTEXT.md §Deferred` |
| P2-3 | `ttfa_ms` semantics fix (`worker_pcm.go` stamps `wall` at `final`, so the field is total synthesis time, not first-audio) | Separate fix; the festival keeps the field's current meaning and says so. | `design-surface.md §Two findings for the operator` #1, `CONTEXT.md §Deferred` |
| P2-4 | Changing the shipped default ref text (`internal/keep/keep.go:42`) | Changing the text without changing `ref.wav` would make it a false transcript. A voice-lock call for the operator, not this festival's. | `design-fest-ad.md §The professional lock`, `CONTEXT.md §Deferred` |
| P2-5 | Re-cutting `docs/booth.gif` / `booth.mp4` with the native mouth | Out of the slice. | `CONTEXT.md §Deferred` |
| P2-6 | A fair / FIFO queue among waiters | `flock` is not FIFO, and for a document script every waiter is doing the same work, so fairness buys nothing. Do not build a fair queue for a problem nobody has. | `design-queue.md §Lock mechanics` |

## Negative requirements — these must NOT exist when the festival is done

From `seed.md §v2 is not`, `design-pipes.md §Rejected`, `design-recommend.md §v2 is not`, `FESTIVAL_RULES.md`.

| Must not exist | Why |
|----------------|-----|
| `cans read <file>` | Document ingestion. The script owns the document. |
| `-f file.txt` | The same thing wearing a flag. |
| Markdown / code-fence stripping | A text filter, not a mouth. `sed` already exists. |
| A sentence chunker | Cans speaks what it is handed; the caller decides what a line is. |
| `--voice name`, a voice picker, SSML | Keep is the only throat change. |
| A config file | Four flags do not need a config surface. |
| Any flag not in `design-pipes.md` | "If a task adds a flag that is not in `design-pipes.md`, the task is wrong." (`FESTIVAL_RULES.md`) |
| A daemon, a job queue file, job IDs, persisted lock state | The lock is `flock`; the kernel cleans up. No state to garbage-collect. |
| Restyle mid-session | Same woman does not restyle mid-scene; keep stays the only throat change. |
| Mic, VAD, replies, browser booth, radio | Still parked from v1. |
| New module dependencies | The lock uses stdlib `syscall.Flock`. `golang.org/x/sys` stays indirect. |
| Python in the shipped payload | `tapes/render-demo-tape.py` is dev tooling and stays; nothing else. |
| `projects/veronica-voice` or `qwen3-tts-native` in `go.mod` | Engine stays out of the CLI module. |
| A second Festival footer line, or re-litigating the first | It shipped at `c72d4a0` as chrome and stays as shipped. |
