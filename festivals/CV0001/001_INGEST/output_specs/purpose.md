# Purpose — cans v2 (CV0001)

## The slice in one line

**Cans becomes a unix primitive a script can drive over a document: text in from argv or stdin, wav out where you point it, and one mouth at a time no matter how the script loops.** (`design-recommend.md §The slice in one line`, `seed.md`)

The script owns the document. Cans speaks what it is handed. (`user-direction.md` #2, `design-pipes.md §The boundary`)

## What v2 is

| Item | Shape | Source |
|------|-------|--------|
| `-o take.wav` | write the wav instead of playing and deleting it | `design-recommend.md §v2 is` |
| stdin | `echo hi \| cans say`, `cans say -`; empty argv on a TTY stays a usage error | `design-pipes.md §Text in` |
| `--stream` | one utterance per stdin line over **one** `Session`, records flushed per line | `design-pipes.md §Text in` |
| `--json` | JSONL records on stdout; prose on stderr | `design-pipes.md §Streams and exit codes` |
| mouth lock | `flock` on `CANS_HOME/mouth.lock`, held for the lifetime of one `Session` | `design-queue.md §Lock mechanics` |

The lock is the only genuinely new machinery. Everything else is plumbing onto what `WI-7eb171` already shipped. (`design-queue.md §What the merge already gave us`)

## Why it matters

**The one-shot tax.** `tts.SayWith` opens a worker and defers `Close` per call (`internal/tts/synth.go:44-48`), and every `Open` loads the GGUF model before answering `ready`. A 200-line loop pays 200 model loads. The booth already runs warm on one `Session` (`internal/booth/booth.go:149-156`); the CLI does not. (`design-surface.md §What did not change: the one-shot tax`)

**The `xargs -P 8` failure.** Eight concurrent `cans say` calls put eight workers with model weights resident in unified memory. Nothing refuses, nothing waits — swap, thrash, wedge. Going native changed the mechanism (C++/GGML/Metal instead of Python/MLX); it did not change the arithmetic. (`design-queue.md §The problem, stated from the code`)

Stream mode fixes the first. The lock fixes the second. Neither adds a background process or persisted state. (`design-recommend.md §Why this shape`)

**The fest-ad tape.** The pipe tape is the honest demo of what v2 is for: a terminal, a document the *script* walks, and files landing in `out/`. No narration, no pitch. It sits next to the booth GIF; it does not replace it. (`design-fest-ad.md §What v2 adds`)

**The public festival tree.** The v2 festival lands under `projects/cans/festivals/CV0001/` as the second readable plan a stranger can read without installing `fest`. The ad compounds — one plan per slice. (`design-fest-ad.md §What v2 adds`)

## Done criteria

Functional — from `design-recommend.md §Ship verification` and `design-queue.md §The bar`:

1. `cans say "x" -o take.wav` writes the wav, does not play, does not delete; stdout is the path.
2. `echo x | cans say` and `cans say -` read one utterance from stdin; empty argv on a TTY is still exit 2.
3. `cat lines | cans say --stream -o 'out/%03d.wav'` writes one wav per line over **one** `Session`, with exactly one GGUF load.
4. `--json` emits one record per utterance on stdout, flushed as each finishes; prose is on stderr.
5. `xargs -P 8` over 50 lines: **one** `qwen3-tts-worker` at every `pgrep` sample, no swap.
6. Ctrl-C mid-stream keeps completed wavs, leaves no orphaned worker and no held lock; `kill -9` leaves the next run unblocked.
7. A booth session and a background script never interleave audio.
8. A second VHS tape shows a script piping lines in and wavs landing on disk.

Quality:

9. `cans say "x"` behaves exactly as at `1e8cea2`; every existing test passes.
10. Stream and lock paths are tested on the fake worker (`internal/tts/testdata/fakeworker`) — CI needs no real mouth.
11. A 200-line stream beats the 200-call loop by a **recorded** margin; both numbers live in this festival.
12. `just test unit`, `go vet ./...`, `gofmt -l .` clean; files under 500 lines, functions under 50.
13. README, docs, tapes and the snapshot pass the professional-surface grep (`CONTEXT.md §Professional grep`, campaign-private); exactly one Festival footer. (`design-fest-ad.md §Verification for the ad`)
14. `fest validate` green; `projects/cans/festivals/CV0001/` holds the readable snapshot.

## Not the point

Making cans a document reader, a daemon, or a config surface. Every "cans does not" in `design-pipes.md §The boundary` is the script's job — `sed`, `jq`, `awk`, and a `for` loop already do them better, and each refusal keeps the tool small enough to read in one sitting.
