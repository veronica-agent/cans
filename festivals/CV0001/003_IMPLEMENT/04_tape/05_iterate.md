---
fest_type: gate
fest_id: 05_iterate.md
fest_name: Review Results and Iterate
fest_parent: 04_tape
fest_order: 5
fest_status: completed
fest_autonomy: medium
fest_gate_id: iterate
fest_gate_type: iterate
fest_managed: true
fest_created: 2026-08-21T05:04:57.537114-06:00
fest_updated: 2026-08-21T16:32:54.885927-06:00
fest_tracking: true
fest_version: "1.0"
---


# Gate: Review Results and Iterate

Address all findings from testing and code review. Iterate until the sequence meets quality standards.

## Findings to Address

### From Testing

- [x] No defects. `03_testing.md` raised two carry-forwards, both already recorded and neither a `04_tape` fix: the gif was cut at 1-min load 16.3–27.9 (re-cut assigned to `004_REVIEW`), and two of three quiet-box one-shots cost ~27 s against a 5 652 ms baseline (the mouth's end-of-speech variance, deferred in `CONTEXT.md`, engine-side).

### From Code Review

**Critical**

- [x] `README.md:87` — `$((++i))` inside `$( … )` never advanced `i`, so every iteration wrote `out/001.wav`. Replaced with a parent-shell counter: `i=0` before the loop, `i=$((i+1))` as the first statement inside, and `-o "$(printf 'out/%03d.wav' "$i")"`. Proof below.
- [x] `README.md:91` — the manifest example had no `-o`, so `playTail` played each line and `tts.RemoveTemp` deleted the wav before the next line, leaving `manifest.txt` full of dead paths. Added `-o 'out/%03d.wav'` so the recorded paths exist and nothing goes to the speakers.
- [x] `README.md:82-83` — `awk -v RS='' '{print}'` kept the paragraph's embedded newlines, so `--stream` spoke one wav per *source* line, not per paragraph. Now `awk -v RS='' '{gsub(/\n/," "); print}'`. Proof below.

**Suggestions**

- [x] `cmd/cans/main.go:39` — exit-75 line now reads `exit 75 when another cans holds the mouth and --nowait was set or --wait ran out`, matching `mouth.ErrBusy` → `say.ExitBusy` on both paths and the README's "refused or ran out".
- [x] `cmd/cans/main.go:30` — dropped the bare `cans say <text>              speak one line`; the flagged synopsis at line 10 of the const is the only `say` entry now. Const is **14 lines** (limit 20).
- [x] `tapes/pipe.tape:13` — `Set Height 380` → `Set Height 300`, same as `booth.tape:12`.
- [x] Carry-forward acknowledged, not actioned here: `docs/pipe.gif` stays as recorded (680×380, 123 s, `ttfa_ms` 31377/34570/35202). Not re-recorded — another agent holds the real mouth, and `004_REVIEW` owns the re-cut under load < 16.

## Iteration

### Not in the review, found while fixing it

Dropping the bare `cans say <text>` line broke a test the review had not flagged and that `03_testing.md` wrongly reported as absent: `cmd/cans/say_args_test.go:19-20` asserted the `-h` / `--help` error contains the literal `"cans say <text>"`.

```
--- FAIL: TestParseSayErrors/help_short
    say_args_test.go:38: error "cans — put the cans on. …" does not contain "cans say <text>"
--- FAIL: TestParseSayErrors/help_long
FAIL	github.com/veronica-agent/cans/cmd/cans	2.172s
```

Fixed by moving both rows onto the surviving synopsis — `"cans say [-o out.wav]"` — which keeps the test's intent (`-h` prints the say usage) and is the line the review asked to keep. My statement in `03_testing.md` that "nothing asserts on the text of the `usage` const" was wrong: I had grepped `main_test.go` only. The claim in `03_testing.md` should be read as corrected here.

### Proof for the two shell findings

Run in scratch against the **verbatim** block now in `README.md`, with a `cans` stub that echoes only its `-o` argument (nothing touched the real mouth):

Fixture — `chapter.md` is two paragraphs, the first wrapped over two source lines; `lines.txt` is the three tape lines.

```
== (a) BROKEN: $((++i)) inside $( ) ==
out/001.wav
out/001.wav
out/001.wav

== (a) FIXED (shipped text) ==
out/001.wav
out/002.wav
out/003.wav

== (c) BROKEN: awk -v RS='' '{print}' ==
A paragraph that the author
wrapped over two source lines.
A second paragraph on one line.
stdin lines = 3          # 3 wavs for 2 paragraphs

== (c) FIXED (shipped text): gsub flattens the record ==
A paragraph that the author wrapped over two source lines.
A second paragraph on one line.
stdin lines = 2          # 2 wavs for 2 paragraphs
```

Finding (b) needs no shell proof: it is the presence of `-o 'out/%03d.wav'`, which is what stops `playTail` from taking the temp-wav branch (`internal/say/say.go:65-70`).

### Checks after the changes

```
$ gofmt -l .                                  (no output)
$ go vet ./...                                (no output)
$ CANS_NOPLAY=1 go test -count=1 ./...        ok in all ten packages
$ rg -i '<pattern 1 — CONTEXT.md §Professional grep>' README.md | rg -v '<allowed footer line>'
                                              (no output)
$ rg -i '<pattern 2 — CONTEXT.md §Professional grep>' README.md docs/ tapes/
                                              (no output)
$ rg -c 'fest.build' README.md                1
$ rg -n '\bwe\b|!' README.md                  (no output)
$ rg -n 'docs/phrases' README.md              (no output)
```

`## Scripting` is now **48 lines** (was 46; the parent-shell counter adds two), still under 60. No build, no `bin/cans`, no `just vhs`, no worker — another agent holds the real mouth.

## Definition of Done

- [x] All critical findings fixed
- [x] All tests pass after changes
- [x] Linting passes
- [x] Code review findings addressed
- [x] Ready to commit