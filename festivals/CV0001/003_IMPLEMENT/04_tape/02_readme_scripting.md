---
fest_type: task
fest_id: 02_readme_scripting.md
fest_name: readme_scripting
fest_parent: 04_tape
fest_order: 2
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:57.057442-06:00
fest_updated: 2026-08-21T16:18:13.148031-06:00
fest_tracking: true
---


# Task: readme_scripting

## Objective

A README "Scripting" section and an updated `usage` const: the loops, the flags, the exit codes, the one honest line about Ctrl-C (per D014, not D008) — boring, technical, and clean on the grep.

## Requirements

- [x] New `## Scripting` section after `## Commands`, in this order: one sentence (`The script owns the document. cans speaks what it is handed.`); `docs/pipe.gif` embedded at width 680 like the booth gif; the three loops from `design-pipes.md §Loops this makes possible` with boring text (measurement-style lines or the three tape lines — never product copy); a flag table (`-o`, `--stream`, `--json`, `--play`, `--nowait`, `--wait`, `-`); an exit-code table (`0`, `1`, `2`, `75`, `130`); the line `Ctrl-C stops the stream: the line being spoken is dropped, finished wavs stay, exit 130. A second Ctrl-C stops at once.` (D014); the line `Without -o the wav is a temp file removed after playback.`
- [x] Keep the section under 60 lines. No new headings elsewhere. The footer `Built with [Festival](https://fest.build)` stays exactly once — do not add a second Festival line and do not reword the first.
- [x] `cmd/cans/main.go` `usage` const: add `cans say [-o out.wav] [--json] [--play] [--nowait|--wait 30s] <text>`, `echo text | cans say`, and `cans say --stream -o 'out/%03d.wav' < lines.txt` lines; keep the const under 20 lines.
- [x] `docs/phrases` is campaign-private — do not reference it from the README. Run the professional-surface grep from `CONTEXT.md §Professional grep`; it must print nothing; `rg -c 'fest.build' README.md` must print `1`.

## Implementation

1. Write the section; keep the existing README voice (short lines, no exclamation marks, no "we").
2. Update `usage`; run `./bin/cans --help` and `./bin/cans say -h` and read them.
3. Run the greps; fix anything they catch.

## Done when

- [x] Grep empty; footer count 1; README read top to bottom once
- [x] `CANS_NOPLAY=1 go test ./...` green (the `usage` change may touch `main_test.go` expectations)

## Result

`## Scripting` added as the last section of `README.md`, after `## Commands` and its two closing paragraphs, immediately above the `---` footer. **46 lines** (limit 60). Contents in the required order: the one sentence, `docs/pipe.gif` at `width="680"` in the same centred `<p>` shape as the booth gif, one mechanism line (`one model load for the whole document` — no number, no margin, since `04_measure` is still blocked), the three loops from `design-pipes.md §Loops this makes possible` verbatim in shape (`chapter.md`, `lines.txt`, `manifest.txt`), a 7-row flag table, a 5-row exit-code table, the stdout/stderr line, the D014 Ctrl-C line word for word, and the temp-file line.

`ttfa_ms` is described as "the worker's total synthesis time for that line", which is what `worker_pcm.go` actually stamps — the deferred semantics item in `CONTEXT.md` says the festival keeps the field's meaning and says so.

`cmd/cans/main.go` `usage`: three lines added after the command list — `cans say [-o out.wav] [--json] [--play] [--nowait|--wait 30s] <text>`, `echo text | cans say`, `cans say --stream -o 'out/%03d.wav' < lines.txt`. The const is **16 lines** (limit 20). `usage` is only ever passed to `fmt.Fprint` and `fmt.Errorf("%s", …)`, so the `%03d` is not a format hazard. `./bin/cans --help` (exit 0) and `./bin/cans say -h` (exit 2, usage on stderr — pre-existing `parseSay` behaviour, untouched) both print it.

Checks: `rg -i '<pattern 1 — CONTEXT.md §Professional grep>' README.md | rg -v '<allowed footer line>'` — nothing. `rg -i '<pattern 2 — CONTEXT.md §Professional grep>' README.md docs/ tapes/` — nothing. `rg -c 'fest.build' README.md` — `1`. `rg -n '\bwe\b|!' README.md` — nothing. No reference to `docs/phrases`. `gofmt -l .` empty, `go vet ./...` clean, `CANS_NOPLAY=1 go test ./...` green in all ten packages (`cmd/cans` re-ran at 0.538s after the `usage` change; no test asserted on the const's text).