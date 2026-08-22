---
fest_type: gate
fest_id: 04_review.md
fest_name: Code Review
fest_parent: 04_tape
fest_order: 4
fest_status: completed
fest_autonomy: low
fest_gate_id: review
fest_gate_type: review
fest_managed: true
fest_created: 2026-08-21T05:04:57.532466-06:00
fest_updated: 2026-08-21T16:28:55.3774-06:00
fest_tracking: true
fest_version: "1.0"
---


# Gate: Code Review

Reviewed cold by a second agent: `git diff HEAD` plus untracked `tapes/pipe.tape` and `docs/pipe.gif`. Nothing was executed against the real mouth.

## Review Checklist

### Code Quality

- [x] Code is readable and well-organized
- [x] Functions are focused (single responsibility)
- [x] Naming is clear and consistent
- [x] No unnecessary complexity or duplication

### Standards Compliance

- [x] Linting passes without warnings
- [x] Formatting is consistent
- [x] Project conventions are followed

### Error Handling & Security

- [x] Errors are handled appropriately
- [x] No secrets in code
- [x] Input validation present where needed
- [x] No obvious security issues

### Alignment

- [x] Changes align with sequence goal
- [x] No scope creep beyond what was requested

The boxes cover the Go change — three lines in the `usage` const, clean. The Criticals are in the
README's shell examples: shipped documentation that does not do what its comment says.

## Findings

**Critical Issues:** (must fix)

- `README.md:87` — `cans say "$line" -o "out/$(printf %03d $((++i))).wav"` writes `out/001.wav` on
  every iteration. `$((++i))` is evaluated inside the `$( … )` subshell, so `i` never changes in the
  parent shell. Verified: `while IFS= read -r line; do echo "out/$(printf %03d $((++i))).wav"; done`
  over a 3-line file prints `out/001.wav` three times. Anyone who copies the loop silently overwrites
  every wav but the last — the exact failure `-o` exists to prevent. Fix: increment in the parent —
  `i=0` before the loop, `i=$((i+1))` inside, `-o "$(printf 'out/%03d.wav' "$i")"`.
- `README.md:91` — `cans say --stream --json < lines.txt | jq -r 'select(.error==null) | .wav' >
  manifest.txt` writes a manifest of files that no longer exist. With no `-o`, `playTail`
  (`internal/say/say.go:65-70`) plays the temp wav and calls `tts.RemoveTemp` immediately after the
  record is emitted (`internal/say/stream.go:122-129`), so every path in `manifest.txt` is deleted
  before the next line starts — and every line is played through the speakers, which a build step
  does not want. Fix: add the output template —
  `cans say --stream --json -o 'out/%03d.wav' < lines.txt | jq -r 'select(.error==null) | .wav' > manifest.txt`.
- `README.md:82-83` — `awk -v RS='' '{print}' chapter.md | cans say --stream -o 'out/%03d.wav'` does
  not give "one wav per paragraph". In paragraph mode `$0` keeps the paragraph's embedded newlines,
  so `print` emits them and `--stream` speaks one wav per *source* line. Verified: a 2-line paragraph
  plus a 1-line paragraph yields 3 stdin lines, so wrapped prose fragments into per-line wavs. Fix:
  flatten the record — `awk -v RS='' '{gsub(/\n/," "); print}' chapter.md | cans say --stream -o 'out/%03d.wav'`.

All three are verbatim from `design-pipes.md:99-110`; the pack carries the same bugs. Fix them here.

**Suggestions:** (should consider)

- `cmd/cans/main.go:39` — `exit 75 when another cans holds the mouth and --nowait was set` is now
  incomplete: line 35 advertises `--wait 30s`, and an expired `--wait` returns 75 too
  (`internal/mouth/lock.go` → `ErrBusy` → `say.ExitBusy`, `internal/say/say.go:99-103`). The README
  exit table gets it right ("refused or ran out"); make `usage` match.
- `cmd/cans/main.go:30,35` — `cans say <text>` is listed twice, bare and with flags; drop the bare one.
- `tapes/pipe.tape:13` — `Set Height 380` leaves ~130 px of empty terminal under the returned prompt
  in the last frame. `Set Height 300`, as in `booth.tape:12`, frames the run and trims bytes.
- Carry-forward, already recorded, not a defect: `docs/pipe.gif` was cut at 1-min load 16.3–27.9, so
  the visible `ttfa_ms` reads **31377 / 34570 / 35202** and the gif runs 123 s. `01_pipe_tape.md` and
  `03_testing.md` both record this and assign the re-cut to `004_REVIEW`. Re-cut with `just vhs pipe`
  under load < 16 before the PR; the tape is deterministic and needs no edit.

## Verified

- **README**: `## Scripting` at line 71, after `## Commands`; **46 lines** (< 60). Order as specified:
  sentence, `docs/pipe.gif` at `width="680"` in the booth gif's centred `<p>`, mechanism line, three
  loops, 7-row flag table, 5-row exit table (`0/1/2/75/130`), stdout/stderr line, the D014 Ctrl-C and
  temp-wav lines **verbatim**. No new headings.
- **Tables match the code**: `--play` needs `-o` (`say_args.go:91-93`); `--stream` + argv text is
  usage (`:97-99`); `-` + text is usage (`:94-96`); 75 on `--nowait` and on expired `--wait`
  (`exit.go:9`, `say.go:100-102`); 130 on interrupt (`say.go:104-107`, `stream.go:135-142`); stream
  exits 1 when any line failed (`stream.go:104-106`); `-o` under `--stream` needs exactly one `%d`
  (`template.go:12-18`). No flag outside the seven exists.
- **No speed number or margin** anywhere in the new section — only the one-worker invariant. The one
  number in `README.md` is the pre-existing 1.6 GB download at line 27.
- **Greps** (`CONTEXT.md §Professional grep`) over `README.md docs/ tapes/`: patterns 1 and 2 print
  nothing; `rg -c 'fest.build' README.md` prints `1`; `rg -n '\bwe\b|!' README.md` and
  `rg -n 'docs/phrases' …` print nothing. The Festival footer appears exactly once.
- **`usage` const**: `cmd/cans/main.go:25-40`, **16 lines** (< 20), boring; `%03d` is safe because
  `usage` only reaches `fmt.Fprint` and `fmt.Errorf("%s", …)`.
- **`tapes/pipe.tape`**: three typed lines exact; real mouth (no `CANS_SAY_BIN`, no `CANS_NOPLAY`, no
  fake worker); preamble exports `PATH` from `$PWD/bin` **before** `cd "$(mktemp -d)"`; both
  `Wait+Screen` guards present; `Output docs/pipe.gif` only; theme line byte-identical to `booth.tape:17`.
- **`.justfiles/vhs.just`**: `pipe:` added at 35-38 (`just build quick`, `just vhs record`);
  `booth:` and `demo:` untouched — the diff is +5 lines, nothing else.
- **`docs/pipe.gif`**: 168 370 B (164 KB, limit 2 MB), 680×380, 123.44 s. Late frame extracted to
  scratch and read: three JSON records, `> ls out`, `001.wav 002.wav 003.wav`, prompt returned.
- **Working tree**: exactly the five expected paths. No `lines.txt`, no `out/`.
- **Go**: `gofmt -l .` empty, `go vet ./...` clean, `CANS_NOPLAY=1 go test ./...` green in all ten
  packages. Only `usage` changed, so `cans say "x"` is byte-identical to `1e8cea2`.

**Verdict: iterate.** Fix the three README examples in `05_iterate.md`, then commit.