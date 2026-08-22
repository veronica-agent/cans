---
fest_type: task
fest_id: 01_pipe_tape.md
fest_name: pipe_tape
fest_parent: 04_tape
fest_order: 1
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:57.054275-06:00
fest_updated: 2026-08-21T16:16:00.050873-06:00
fest_tracking: true
---


# Task: pipe_tape

## Objective

`tapes/pipe.tape` → `docs/pipe.gif` via `just vhs pipe`: a terminal, a script the user owns, and files landing in `out/`. No narration, no pitch.

## Requirements

- [x] Model the tape on `tapes/booth.tape` (same theme block, `Set Width 680`, font size, `Env TERM` / `COLORTERM` / `CLICOLOR_FORCE`). Output `docs/pipe.gif` only (no mp4 — there is no audio to mux).
- [x] The visible script, typed: `printf 'Put the cans on.\nOne worker, one model load.\nFiles land where the script points.\n' > lines.txt`; `cat lines.txt | cans say --stream -o 'out/%03d.wav' --json`; `ls out`. Nothing else. Those three lines **are** the example text — do not change them.
- [x] Real mouth (no `CANS_SAY_BIN`, no fake worker). Hidden preamble as in `tapes/demo.tape.in`: `export PATH="<repo>/bin:$PATH" PS1="> "`, and `cd "$(mktemp -d)"` so `lines.txt` and `out/` never land in the repo.
- [x] Timing: if `vhs manual` shows a `Wait` command in the installed version, use `Wait+Screen /003.wav/` (with a timeout) before `ls out`; otherwise `Sleep` generously (30 s is fine).
- [x] `.justfiles/vhs.just`: `pipe:` recipe → `just build quick`, `just vhs record tapes/pipe.tape`. Keep `booth:` and `demo:` untouched.
- [x] `docs/pipe.gif` under 2 MB. Check a late frame: `ffmpeg -y -sseof -0.5 -i docs/pipe.gif -frames:v 1 frame.png` (in scratch) and look at it — the three JSON records and `ls out` must be readable.

## Implementation

1. `cp tapes/booth.tape tapes/pipe.tape`, change `Output`, replace the typed section, keep the `Hide` / `Show` preamble pattern from `demo.tape.in`.
2. Add the recipe; run `just vhs pipe`; inspect the frame; iterate on `Sleep` / `Wait` until the last frame shows the `ls out` listing.
3. Run the professional-surface grep (`CONTEXT.md §Professional grep`) over `tapes/` and `docs/` — it must print nothing.

## Done when

- [x] `just vhs pipe` regenerates `docs/pipe.gif` from a clean checkout with the mouth installed
- [x] Frame check done and recorded in the testing gate (describe what the last frame shows)
- [x] `git status` shows only `tapes/pipe.tape`, `.justfiles/vhs.just`, `docs/pipe.gif` — no `lines.txt`, no `out/`

## Result

`just vhs pipe` → `docs/pipe.gif`, **164 KB** (limit 2 MB), 123.4 s, 680×380, 12 fps, real mouth (no `CANS_SAY_BIN`, no fake worker). Recorded 16:12:31–16:14:58 on 2026-08-21, first attempt. Both `Wait+Screen` guards fired (`vhs` 0.11.0): `Wait+Screen@300s /003\.wav/` before `ls out`, `Wait+Screen@30s /(?m)^001\.wav/` on the listing.

**Load at record time: 1-minute load 20.88 mean / 16.26 min / 27.93 max over 33 five-second samples — above the 16 bar for the whole run.** The box was never under 16 while the tape ran. Consequence, visible in the gif: `ttfa_ms` reads **31377 / 34570 / 35202** for three short lines, against a 5652 ms quiet-box baseline for the same one-shot text, and the gif is 123 s long instead of ~40 s. Nothing is wrong with the tape or the code — the numbers on screen are the loaded box. **004_REVIEW should re-cut this gif with `just vhs pipe` when the 1-minute load is under 16**; the tape is deterministic and needs no edit to do it.

Frame check (`ffmpeg -y -sseof -0.5 -i docs/pipe.gif -frames:v 1 frame.png`, in scratch): the last frame shows the typed `printf … > lines.txt` (wrapped over two rows), the typed `cat lines.txt | cans say --stream -o 'out/%03d.wav' --json`, three JSON records — `{"line":1,"wav":"out/001.wav","ttfa_ms":31377,"sample_rate":24000}` and the same for lines 2 and 3, each wrapping its trailing `00}` onto a second row at 70 columns — then `> ls out` and the listing `001.wav 002.wav 003.wav`, then the prompt back. All readable.

Professional-surface grep over `README.md docs/ tapes/`: both patterns print nothing; `rg -c 'fest.build' README.md` prints `1`.

`git status --short` after the record: `M .justfiles/vhs.just`, `?? docs/pipe.gif`, `?? tapes/pipe.tape` — no `lines.txt`, no `out/` (the tape's hidden preamble does `cd "$(mktemp -d)"`).