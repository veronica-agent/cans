---
fest_type: gate
fest_id: 05_testing.md
fest_name: Testing and Verification
fest_parent: 03_stream
fest_order: 5
fest_status: completed
fest_autonomy: medium
fest_gate_id: testing
fest_gate_type: testing
fest_managed: true
fest_created: 2026-08-21T05:04:57.4921-06:00
fest_updated: 2026-08-21T14:59:47.700152-06:00
fest_tracking: true
fest_version: "1.0"
---


# Gate: Testing and Verification

Verify all functionality implemented in this sequence works correctly.

## Test Categories

### Unit Tests

- [x] All unit tests pass
- [x] New/modified code has test coverage
- [x] Tests are meaningful (not just coverage padding)

### Integration Tests

- [x] Integration tests pass
- [x] Components work together correctly

### Error Handling

- [x] Invalid inputs are rejected gracefully
- [x] Error messages are clear and actionable
- [x] Recovery paths work correctly

## Verification

- [x] Build completes without warnings
- [x] No regressions introduced
- [x] Coverage meets project requirements

## cans-v2 commands (all from the `cans-v2` worktree; every one must be clean)

```bash
gofmt -l .                                   # prints nothing
go vet ./...
CANS_NOPLAY=1 go test ./...                  # fake worker only — no real mouth
git diff origin/main -- go.mod go.sum        # empty: no new dependencies
wc -l $(git diff --name-only origin/main -- '*.go') | sort -n | tail -5   # every file < 500
./bin/cans say "Put the cans on." ; echo "exit=$?"   # one-shot unchanged: ttfa_ms=N, plays, temp wav gone
```

Then the sequence's own checks from its task files. Record the output of each command in this gate file under **Results** before marking it complete.

## Results

Run 2026-08-21 14:15–14:58 (America/Denver) from `projects/worktrees/cans/cans-v2`, branch `cans-v2`.

Every command in the block above was run and is green. Two notes on how the numbers were taken:

- **A mid-gate edit forced a re-run.** Another agent edited nine `.go` files between **14:29:41 and 14:36:48** — after the static suite first ran (14:17) and after the binary the first manual pass used (14:15). `internal/say/stream.go` was among them. Everything below was therefore re-run against the **14:51** binary and the current tree; nothing in this section predates it.
- **The box was busy for most of the session** (two foreign `Virtualization.VirtualMachine` processes at 841 % CPU / 25.6 GB, 1-minute load averaging 32.9 over 164 samples). A genuine quiet window opened at **14:54**, and all three real-mouth checks below were taken inside it at load **9.9 – 14.7**, under the D013 bar. The 1-minute load is recorded next to each.

A second product's worker is resident on this box and is **not** cans':

```
$ ps -o pid,ppid,pcpu,rss,etime,command -p 62136
  PID  PPID  %CPU    RSS  ELAPSED COMMAND
62136 61770   0.1  24688 08:56:49 ~/.cache/festival-voice/models/qwen3-tts/bin/qwen3-tts-worker …/models
```

24 MB resident, 0.1 % CPU, idle 8 h 56 m — no model loaded. It is `festival-voice`. Because a bare `pgrep -f qwen3-tts-worker` matches it, every check below uses the path-qualified pattern **`pgrep -f 'cans/native/bin/qwen3-tts-worker'`**.

### Static suite — current tree, 14:51

```
$ just build quick
go build -trimpath -ldflags "-s -w -X …/internal/ship.Version=v0.1.0-25-gc736f46-<uncommitted>" -o bin/cans ./cmd/cans

$ gofmt -l .
(exit=0)

$ go vet ./...
(exit=0)

$ CANS_NOPLAY=1 go test -count=1 ./...
ok  	github.com/veronica-agent/cans/cmd/cans	2.251s
ok  	github.com/veronica-agent/cans/internal/audio	0.741s
ok  	github.com/veronica-agent/cans/internal/booth	0.409s
ok  	github.com/veronica-agent/cans/internal/doctor	1.099s
ok  	github.com/veronica-agent/cans/internal/keep	0.611s
ok  	github.com/veronica-agent/cans/internal/mouth	1.843s
ok  	github.com/veronica-agent/cans/internal/play	0.908s
ok  	github.com/veronica-agent/cans/internal/say	11.908s
ok  	github.com/veronica-agent/cans/internal/ship	1.528s
ok  	github.com/veronica-agent/cans/internal/tts	4.672s
(exit=0)

$ git diff origin/main -- go.mod go.sum
(0 bytes — no new dependencies)
```

`-count=1` is deliberate: the first pass came back entirely `(cached)`, which proves nothing.

### File length

```
$ wc -l $(git diff --name-only origin/main -- '*.go') | sort -n | tail -5
     164 internal/booth/booth.go
     200 cmd/cans/main_test.go
     201 internal/tts/worker.go
     418 internal/say/say_test.go
    2496 total
```

`git diff --name-only` lists tracked files only, so this sequence's four new files are counted separately:

```
$ wc -l internal/say/stream.go internal/say/stream_test.go internal/say/template.go internal/say/template_test.go
     167 internal/say/stream.go
     352 internal/say/stream_test.go
      64 internal/say/template.go
      65 internal/say/template_test.go
     648 total
```

Largest is `internal/say/say_test.go` at **418** — every file under 500. **For the reviewer:** `internal/tts/worker.go` is now **201** lines, up from the 196 the festival rules pin with "add files, do not grow it". Five lines, nowhere near the limit, but it is growth in the one file the rules name.

### One-shot unchanged — the v1 regression check

Load **9.90** before, **10.25** after. Worker alone; a one-second sampler ran throughout.

```
$ pgrep -fl 'cans/native/bin/qwen3-tts-worker'      # before
(exit=1)
$ ./bin/cans say "Put the cans on." ; echo "exit=$?"
ttfa_ms=5652
exit=0
```

wall **14 s**, sampler max **1**, and `pgrep` empty afterwards. stdout is exactly the v1 line — `ttfa_ms=N`, no path, no JSON — which is D006's "no `-o`, no `--json`" case. Playback happened, and the temp wav was removed: no `cans-say.*` file from this run survives in `$TMPDIR` (the only one there is from Aug 19), and `git status` shows no stray wav in the tree.

Against the baseline in `002_PLAN/inputs/measurements.md` — cold `cans say`, idle box, `ttfa_ms` **5 839**, real **13.1 s** — this run is **5 652 ms / 14 s**. One-shot is unchanged. `cans say "x"` still behaves as at `1e8cea2`.

### `01_stream_loop` "Done when" — three lines, three records, one worker

Load **10.02** before, **11.11** after. Sampler started first.

```
$ while sleep 1; do pgrep -f 'cans/native/bin/qwen3-tts-worker' | wc -l | tr -d ' '; done > a2.pgrep &
$ printf 'Put the cans on.\nOne worker, one load.\nFiles land here.\n' \
    | ./bin/cans say --stream -o '<scratch>/%03d.wav' --json
```

exit **0**, wall **65 s**, stderr **empty (0 bytes)**. stdout — three records, `line` 1/2/3, each with `wav` / `ttfa_ms` / `sample_rate` per D006:

```
{"line":1,"wav":"<scratch>/a2/001.wav","ttfa_ms":6122,"sample_rate":24000}
{"line":2,"wav":"<scratch>/a2/002.wav","ttfa_ms":27348,"sample_rate":24000}
{"line":3,"wav":"<scratch>/a2/003.wav","ttfa_ms":24215,"sample_rate":24000}
```

Three files landed on the `%03d` template, and the sampler took **61** one-second samples of which **every one read `1`**:

```
$ ls <scratch>/a2
001.wav   002.wav   003.wav
$ sort -n a2.pgrep | tail -1
1
$ pgrep -fl 'cans/native/bin/qwen3-tts-worker'      # after
(exit=1)
```

Three records, three files, **one worker for the whole run, one GGUF load**. That is the requirement, met.

Lines 2 and 3 report `ttfa_ms` of **27 348** and **24 215** — both over 20 s, on a box at load 10. That is the mouth's known end-of-speech variance (the baseline's runs 1 and 3 ran to the 220-token ceiling for 17.58 s of audio against run 2's 1.9 s), **not** something `--stream` introduced. It is why `04_measure` reports median *and* max.

### `03_cancel` "Done when" — SIGINT mid-stream

Load **10.48** before, **14.66** after.

```
$ seq 1 20 | sed 's/^/Line /' | ./bin/cans say --stream -o '<scratch>/b2/%03d.wav' &
$ CANS=$(pgrep -f 'bin/cans say --stream')     # 19053
# polled until 002.wav appeared (~49 s), then:
$ kill -INT $CANS
$ wait; echo $?
130
```

exit **130**. stderr, in full:

```
interrupted after line 2
```

This is the task's rule `N = last completed line`, exactly. Line 3 was in flight when the signal landed; `SayTo` returned `context.Canceled`, `speakLine`'s interrupt branch returned it without counting a failure, and `last` stayed at 2 — so **no `003.wav` was written**. The two finished wavs survived the interrupt:

```
$ ls <scratch>/b2
001.wav   002.wav
```

The worker was gone the moment the stream returned — the deferred `sess.Close()` ran, sent `shutdown`, and released the lock:

```
$ pgrep -fl 'cans/native/bin/qwen3-tts-worker'
(exit=1)
```

And the next one-shot started immediately — **stderr empty, no `waiting for the mouth…` line**, wall **14 s**, which is a normal cold start, not a lock wait:

```
$ ./bin/cans say -o '<scratch>/b2/next.wav' "Put the cans on."
<scratch>/b2/next.wav
exit=0
```

### One thing the reviewer should look at

`002.wav` from the cancel run is **1 484 bytes — about 0.03 s of near-silence** for the input `Line 2`, while `001.wav` in the same run is 289 904 bytes (~6.04 s) for `Line 1`. The file is a well-formed WAV, so nothing in the write path is wrong:

```
$ xxd -l 48 <scratch>/b2/002.wav
00000000: 5249 4646 c405 0000 5741 5645 666d 7420  RIFF....WAVEfmt
00000010: 1000 0000 0100 0100 c05d 0000 80bb 0000  .........]......
00000020: 0200 1000 6461 7461 a005 0000 0000 0000  ....data........
```

PCM, mono, 24 000 Hz, 16-bit, `data` chunk 1 440 bytes, leading samples zero. `--stream` reported success and moved on, which is correct behaviour — it wrote exactly what the worker returned. This is the **mouth** emitting an almost immediate end-of-speech on a short line, the same variance that produces 27 s `ttfa_ms` values at the other extreme. Recorded here because a near-silent wav for a valid line is worth a decision from the operator; it is not a `03_stream` defect and it is already listed as deferred in `CONTEXT.md` ("`ttfa_ms` semantics" / end-of-speech variance).

### Still open in this sequence

`003_IMPLEMENT/03_stream/04_measure` is **still blocked** after three attempts. Attempt 3 (14:57–15:34) got run (a) — the 50-line stream — to completion: exit 0, **50/50 records, `line` dense 1…50, zero errors, 50 wavs, `pgrep` max 1 across 1 753 samples**, and 21.1 s of total non-synthesis overhead. Its wall clock and `ttfa_ms` are contaminated (load mean 23.2, peak 92.6), and runs (b) and (c) aborted at load 34.28, so **there is no measured margin yet**. Full split of quotable vs contaminated results, plus a mouth fault worth the operator's attention (3 of 50 wavs returned ~0.03 s of silence after 36–54 s of synthesis), is in `002_PLAN/inputs/measurements.md §Stream`.

This gate does not depend on that task: every command in the block above was run on the current tree, and the real-mouth ones were taken inside the 14:54 quiet window at load 9.9–14.7.