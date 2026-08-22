---
fest_type: task
fest_id: 04_measure.md
fest_name: measure
fest_parent: 03_stream
fest_order: 4
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:56.825464-06:00
fest_updated: 2026-08-21T17:49:28.619192-06:00
fest_tracking: true
---




# Task: measure

## Objective

The numbers the README and the festival quote, taken on the real mouth with the machine otherwise idle, recorded with the commands that produced them (D013, P1-4).

## Requirements

- [x] Preconditions, checked and recorded: `uptime` 1-minute load below 16; `pgrep -fl qwen3-tts-worker` empty; `just build quick`. If the load is high, record `deferred: load X at <time>` in `002_PLAN/inputs/measurements.md §Stream` and retry later — do **not** take numbers on a loaded box.
- [x] Input: `lines.txt`, 50 lines, boring and technical, distinct: `for i in $(seq 1 50); do echo "Measurement line $i. The worker stays warm."; done > lines.txt`.
- [x] A sampler running for each measurement: `while sleep 1; do pgrep -f qwen3-tts-worker | wc -l | tr -d ' '; done > <name>.pgrep &` — max value must be 1.
- [x] (a) Stream: `time (cans say --stream -o 'out/stream/%03d.wav' --json < lines.txt > stream.jsonl)`.
- [x] (b) Loop: `time (i=0; while IFS= read -r l; do i=$((i+1)); cans say "$l" -o "out/loop/$(printf %03d $i).wav" --json; done < lines.txt > loop.jsonl)`.
- [x] (c) `xargs -P 8` over 24 lines: `head -24 lines.txt | nl -ba | xargs -P 8 -L1 sh -c 'cans say "$2" -o "out/x/$1.wav" --json' _ > x.jsonl` (adjust quoting until it works; record the final command). `vm_stat | grep Pageouts` before and after — delta must be 0.
- [x] From the JSONL: median and max `ttfa_ms` per mode (`jq -s 'map(.ttfa_ms) | sort | .[length/2|floor], max'`), count of records, count of `error` records.
- [x] Record everything in `002_PLAN/inputs/measurements.md §Stream`: walls, median/max, pgrep max per run, pageouts delta, the margin (loop wall − stream wall) and per-line saved cost (margin / 50). If any line's `ttfa_ms` exceeds 20 s, say so — it is the mouth's known end-of-speech variance, not v2's.

## Implementation

1. Work in a scratch dir **outside** the repo (`mktemp -d`) so `out/` and `*.jsonl` never touch git; call `cans` by absolute path to the worktree's `bin/cans`.
2. Run (a), (b), (c) strictly one after another; kill the sampler between runs; `pgrep -f qwen3-tts-worker` must find nothing between runs (if it does, record it — that is a bug).
3. Write the section with a table per run and a fenced block with the exact commands.

## Result — attempt 4 (2026-08-21 16:24–17:45)

All three runs completed back to back on a quiet box: **100 % of every run's 30-second load samples were below 16** (means 8.36 / 7.06 / 7.45), and the worker count never exceeded **1** in any run.

| | (a) stream 50 | (b) loop 50 | (c) `xargs -P 8` 24 |
|---|---|---|---|
| Wall | 1 512.6 s | 2 109.2 s | 943.2 s |
| Records / errors | 50 / 0 | 50 / 0 | 24 / 0 |
| `ttfa_ms` median / max | 35 730 / 37 614 | 35 842 / 77 932 | 35 736 / 36 714 |
| Lines over 20 s | 39 | 45 | 21 |
| Worker max | 1 | 1 | 1 |
| Pageouts delta | 637 | 3 000 | **0** |

**Margin = 2 109.2 − 1 512.6 = 596.7 s, i.e. 11.93 s per line, uncontaminated.** Of that, **6.44 s/line is the structural overhead `--stream` removes** (non-synthesis cost was 16.8 s total for the stream vs 338.6 s for the loop) — independently reproduced by run (c) at 6.69 s/line, and matching the ~6.6 s GGUF load in the baseline.

Notes: the binary was **not** rebuilt (`bin/cans` at `c443de9` was already built; the orchestrator directed no rebuild), so `just build quick` was not re-run. The prescribed `pgrep -f … | wc -l` sampler **cross-matches concurrent samplers** and read a false 2–3; an anchored `ps` counter was used instead. Full tables, the proof it was an artifact, the exact commands, the near-silent-wav indices and the run (b) wall caveat are in `002_PLAN/inputs/measurements.md §Stream → Attempt 4`.

## Done when

- [x] `002_PLAN/inputs/measurements.md §Stream` has the three runs with commands, walls, median/max, pgrep max = 1 everywhere, pageouts delta 0
- [x] The margin is positive and its per-line value is stated; the 50-line stream started exactly one worker