# BAR — commands and recorded output

Fill each block with the exact command and its output when `004_REVIEW` runs. Real mouth, idle machine (`uptime` first). Do not paraphrase output.

**Run 2026-08-21 18:17–18:40 by the review agent.** Binary `bin/cans` at `v0.1.0-28-g9136cbe`, worktree clean at `9136cbe`. Every real-mouth item was taken one at a time on a quiet box — the 1-minute load never left the 4.6–7.6 band for the whole session, so no item is contaminated and none had to be deferred.

Paths below are written `<scratch>` (a `mktemp -d` outside the repo) and `~`; the recorded commands had them spelled out. Grep patterns are referenced, never quoted — they are campaign-private (`CONTEXT.md §Professional grep`).

**Verdict: 13 / 13 pass.** Two known engine faults were seen again and are recorded where they appeared; neither is a `--stream` defect and both are already deferred in `CONTEXT.md`.

## 0. Preconditions

```bash
$ uptime
18:17  up 18 days, 12:29, 41 users, load averages: 5.97 6.41 6.85

$ ps -Ao command= | grep -c '^~/\.cans/native/bin/qwen3-tts-worker '
0

$ just build quick
mkdir -p <worktree>/bin
cd <worktree> && go build -trimpath -ldflags "-s -w -X github.com/veronica-agent/cans/internal/ship.Version=v0.1.0-28-g9136cbe" -o bin/cans ./cmd/cans
exit=0
```

**The worker counter.** `pgrep -f 'qwen3-tts-worker'` is not used anywhere in this bar. It self-matches and it cross-matches other concurrent samplers, which is how attempt 3 of `04_measure` read a false 2 and 3 with exactly one worker resident. Every count below comes from an anchored `ps` counter in its own script file, which also excludes the foreign `~/.cache/festival-voice/…/qwen3-tts-worker` (PID 62136, resident for this whole session, not cans and not touched):

```sh
ps -Ao pid,command= | grep '^ *[0-9][0-9]* /Users/…/\.cans/native/bin/qwen3-tts-worker '
```

It samples once a second and records **both the count and the PIDs**. The PID column is the stronger evidence for items 1–2: a second GGUF load means a second process, so a single unchanging PID across a whole run is proof of a single load in a way a count of 1 alone is not.

## 1. Stream — one worker, one load

```bash
$ ./bin/cans say --stream -o '<scratch>/out1/%03d.wav' --json < lines3.txt > stream1.jsonl 2> stream1.err
( ... )  111.70s user 10.30s system 117% cpu 1:43.89 total
exit=0

$ cat stream1.jsonl
{"line":1,"wav":"<scratch>/out1/001.wav","ttfa_ms":32591,"sample_rate":24000}
{"line":2,"wav":"<scratch>/out1/002.wav","ttfa_ms":32418,"sample_rate":24000}
{"line":3,"wav":"<scratch>/out1/003.wav","ttfa_ms":32346,"sample_rate":24000}

$ cat stream1.err
(nothing)

$ ls -l <scratch>/out1/
-rw-r--r--  53912 001.wav
-rw-r--r--   1484 002.wav
-rw-r--r-- 126530 003.wav
```

Sampler over the whole run — 98 one-second samples:

```
samples=98   max=1   counts seen: 0 1   distinct worker PIDs: 89253
```

Every sample from the second onward reads `1 89253`. **One process, start to finish** — the worker was started once, loaded the GGUF once and served all three lines. Worker gone after the run (`workers_after=0`).

The overhead confirms the single load arithmetically:

```
Σ ttfa_ms = 32591 + 32418 + 32346 = 97 355 ms = 97.4 s
wall                                          = 103.9 s
everything else (process start, doctor.Prepare, keep.Load, flock, GGUF load, 3 writes, 3 records)
                                              = 6.5 s for all three lines
```

6.5 s is one GGUF load (~6.6 s in the baseline), not three. **PASS.**

*Mouth fault, not a v2 defect:* `002.wav` is a well-formed but near-silent **0.03 s** / 1 484-byte file after 32.4 s of synthesis (`001` is 1.12 s, `003` is 2.64 s). Same shape as the 4/50, 5/50, 1/24 seen in `04_measure`; the worker reported success, so `--stream` had no error to report and correctly wrote what it was given. Deferred in `CONTEXT.md`, engine-side.

## 2. xargs -P 8 — one worker, no pageouts

The 24-line measurement is **not re-run here** — it is `002_PLAN/inputs/measurements.md §Stream → Attempt 4 (c)`, taken on a quiet box with 100 % of its load samples below 16:

| | Attempt 4 (c), 24 lines |
|---|---|
| Wall | 943.2 s (`real 15m43.217s`) |
| Exit / records / wavs | 0 / **24 of 24** / 24 |
| **Worker max** | **1** (878 samples, values `{0,1}`) |
| Pageouts | 875 724 → 875 724 — **delta 0** |
| stderr | 20 × `waiting for the mouth…`, nothing else |

Re-confirmed here with a short 8-line run under the same `-P 8` fan-out, on the same null-delimited pipeline (BSD `xargs` has no `-d`):

```bash
$ cat lines8.txt | nl -ba | sed 's/^[[:space:]]*//' | tr '\t\n' '\0\0' \
    | xargs -0 -P 8 -n 2 sh -c '"$CANS" say "$2" -o "$XOUT/$1.wav" --json' _ > x.jsonl 2> x.err
xargs_exit=0
xargs_wall_s=268.8

$ vm_stat | awk '/Pageouts/{print $NF}'     # before / after
875753. / 875753.

$ wc -l < x.jsonl ; ls <scratch>/out2/ | wc -l
8
8

$ cat x.err
waiting for the mouth…
waiting for the mouth…
waiting for the mouth…
waiting for the mouth…
waiting for the mouth…
waiting for the mouth…
waiting for the mouth…
```

Sampler across the run — 250 one-second samples:

```
samples=250   max=1   counts seen: 0 1
distinct worker PIDs: 9201 10239 10589 11762 12843 13339 14370 14837
```

Eight callers, eight worker processes, **never two at the same instant** — the PIDs are strictly sequential, each starting only after the previous exited. Seven `waiting for the mouth…` for eight callers is exactly right: one holds the lock, seven queue. **Pageouts delta 0. PASS.**

*Mouth fault again:* `5.wav` is 1 484 bytes (1 of 8).

## 3. Ctrl-C mid-stream → 130, wavs kept, no orphan

20 lines in, `SIGINT` sent to the `cans` process one second after `002.wav` appeared.

```
002.wav exists at 18:24:42 — sending SIGINT to 92731
exit_code=130
sigint_to_exit_s=0.05
worker_gone_after_exit_s=0.11
next_oneshot_exit=0
next_oneshot_wall_s=14.65
gap_sigint_to_next_start_s=0.18
```

```bash
$ cat stream3.err
interrupted after line 2

$ cat stream3.jsonl
{"line":1,"wav":"<scratch>/out3/001.wav","ttfa_ms":41986,"sample_rate":24000}
{"line":2,"wav":"<scratch>/out3/002.wav","ttfa_ms":41509,"sample_rate":24000}

$ ls -l <scratch>/out3/
-rw-r--r--  1484 001.wav
-rw-r--r-- 63972 002.wav
                              # no 003.wav — line 3 was in flight and was dropped

$ cat next3.out ; cat next3.err
ttfa_ms=5798
(nothing — no `waiting for the mouth…`)
```

Exit **130**; the two finished wavs stayed; the in-flight line 3 was terminated rather than waited out (D014); the worker was gone **0.11 s** after the process exited; the next one-shot began **0.18 s** after the Ctrl-C and completed normally. **PASS.**

## 4. kill -9 → next run unblocked

`SIGKILL` to the `cans` process ~8 s into a 20-line stream, with the lock file watched on either side.

```
lock_before:      -rw-r--r--  0  Aug 21 06:32  ~/.cans/mouth.lock
kill -9 96323 at 18:25:24
killed_exit=137
workers_1s_after_kill9=1
lock_after_kill:  -rw-r--r--  0  Aug 21 06:32  ~/.cans/mouth.lock
next_oneshot_exit=0
next_oneshot_wall_s=35.36
gap_kill_to_next_start_s=1.10
lock_after_next:  -rw-r--r--  0  Aug 21 06:32  ~/.cans/mouth.lock
```

```bash
$ cat next4.out ; cat next4.err
ttfa_ms=27241
(nothing — no `waiting for the mouth…`)
```

The next `cans say` started **1.10 s** after the kill, never printed the wait line, and exited 0: the kernel dropped the `flock` with the process, exactly as D003 intends. The **lock file itself is never deleted** — same 0-byte inode and mtime before the kill, after the kill and after the next run. **PASS.**

Two things worth recording, neither a failure:

- **The orphaned worker outlives its killed parent for about a second.** `SIGKILL` cannot be handled, so cans gets no chance to shut the worker down; the worker exits on its own when the pipe closes. In the 42-sample trace exactly one sample reads 2 — `96327` (the killed stream's worker, on its way out) alongside `96546` (the new one-shot's) — and the next sample is back to 1. The overlap is a teardown transient during the second the new run spends on `doctor.Prepare`, not two workers synthesising. This is the price of `kill -9` and is what the design accepts in exchange for a lock the kernel always releases.
- `ttfa_ms=27241` for a one-shot whose quiet-box baseline is ~5 700 ms is the mouth's end-of-speech variance (see item 6), not the kill.

## 5. Booth holds the lock

The booth needs a real terminal. `script -q /dev/null ./bin/cans` **cannot** be driven from a fifo on macOS — `script` runs `tcgetattr` on its own stdin and dies with `script: tcgetattr/ioctl: Operation not supported on socket`, the booth never starts, and every `say` fired at it then wrongly succeeds. The booth here is given a real pty from `forkpty` instead, with keys written to the master fd.

```
booth pid=7843 started 18:31:34
booth_worker_up_after_s=5.2   workers=1

--- A_nowait ---
A_nowait_exit=75
A_nowait_wall_s=0.01
A_nowait_stdout=''
A_nowait_stderr='say: mouth busy\n'

--- B_wait2s ---
B_wait2s_exit=75
B_wait2s_wall_s=2.01
B_wait2s_stdout=''
B_wait2s_stderr='waiting for the mouth…\nsay: mouth busy\n'

--- C: end the booth ---
still alive after q — sending Esc
booth_exit_status=0
workers_after_booth=0

--- D_after ---
D_after_exit=0
D_after_wall_s=13.60
D_after_stdout='ttfa_ms=5627\n'
D_after_stderr=''
```

The booth took the lock 5.2 s after launch and held it for its whole run (D001). `--nowait` refused in **0.01 s** with **75** and `say: mouth busy` and never touched the mouth. `--wait 2s` printed `waiting for the mouth…`, polled for **2.01 s** — the requested budget, to the hundredth — then gave up with **75**. The booth quit on Esc (`q` is not its quit key), its worker went with it, and the next `say` ran immediately. Sampler max **1** over 27 samples; the two PIDs (booth's, then the post-booth say's) are sequential. **PASS.**

## 6. One-shot unchanged; tests green on the fake worker

```bash
$ gofmt -l .
(nothing)
$ go vet ./...
(nothing)
$ CANS_NOPLAY=1 go test -count=1 ./...
ok  	github.com/veronica-agent/cans/cmd/cans	2.111s
ok  	github.com/veronica-agent/cans/internal/audio	0.739s
ok  	github.com/veronica-agent/cans/internal/booth	1.002s
ok  	github.com/veronica-agent/cans/internal/doctor	1.205s
ok  	github.com/veronica-agent/cans/internal/keep	1.340s
ok  	github.com/veronica-agent/cans/internal/mouth	2.389s
ok  	github.com/veronica-agent/cans/internal/play	0.220s
ok  	github.com/veronica-agent/cans/internal/say	7.285s
ok  	github.com/veronica-agent/cans/internal/ship	1.519s
ok  	github.com/veronica-agent/cans/internal/tts	3.786s
exit=0
```

Ten packages, no cached results, the stream path exercised on the fake worker. The real-mouth one-shot, unchanged from `1e8cea2`:

```bash
$ ./bin/cans say "Put the cans on."
ttfa_ms=7157
( ... )  13.72s user 0.81s system 90% cpu 16.004 total
exit=0

sampler: max=1  counts seen: 0 1  (16 samples)   workers_after=0
temp wavs matching $TMPDIR/cans-say.*:  1 before, 1 after
```

One line of output in the v1 shape, exit 0, one worker, worker gone afterwards. The temp-wav count is unchanged across the run — the single match is a 355-byte leftover dated **2026-08-19**, older than this branch and untouched; this run created no temp file and left none behind. `ttfa_ms` 7 157 against a 5 652 / 5 669 / 5 798 / 5 627 baseline across this session's other one-shots. **PASS.**

## 7. Margin (from measurements.md)

Quoted from `002_PLAN/inputs/measurements.md §Stream → Attempt 4`, the one uncontaminated attempt — three back-to-back runs with **100 % of every run's load samples below 16** (means 8.36 / 7.06 / 7.45, max 14.97) and worker max 1 in each. Not re-measured here: nothing in this phase changed the code, and a 60-minute re-run would only re-derive it.

| | (a) stream, 50 lines | (b) loop, 50 calls |
|---|---|---|
| Wall | **1 512.6 s** | **2 109.2 s** |
| Records / errors | 50 of 50 / 0 | 50 of 50 / 0 |
| Σ reported synthesis | 1 495.8 s | 1 770.6 s |
| **Everything else** | **16.8 s — 0.34 s/line** | **338.6 s — 6.77 s/line** |

```
margin        = 2 109.2 − 1 512.6 = 596.7 s over 50 lines
per-line      = 596.7 / 50        = 11.93 s per line
structural    = 338.6 − 16.8      = 321.8 s  →  6.44 s per line
```

The **6.4 s/line** figure is the defensible one and the one the README quotes: it is the per-call GGUF load `--stream` removes, reproduced independently by run (c) at 6.69 s/line and matching the ~6.6 s load in the baseline. The remaining 274.9 s of the margin is the loop's higher reported synthesis time — real, but the mouth's variance rather than something `--stream` engineered away. Item 1 above reproduces the same structure in miniature: 6.5 s of non-synthesis for a whole 3-line stream. **PASS.**

## 8. Professional-surface grep + one footer

Three patterns from `CONTEXT.md §Professional grep`, run from the worktree over `README.md docs/ tapes/ festivals/CV0001/`. The patterns are campaign-private and are referenced, not quoted.

```bash
$ rg -i '<pattern 1 — CONTEXT.md §Professional grep>' README.md | rg -v '<allowed footer line>'
(nothing — exit 1)

$ rg -i '<pattern 2 — CONTEXT.md §Professional grep>' README.md docs/ tapes/ festivals/CV0001/
(nothing — exit 1)

$ rg -c 'fest.build' README.md
1
```

Nothing / nothing / exactly one Festival footer. **PASS.**

## 9. Fresh CANS_HOME doctor

`bin/cans` copied to a scratch directory outside the checkout, run from a directory outside the checkout, against an empty `CANS_HOME`.

```bash
$ just test unit
cd <worktree> && CANS_NOPLAY=1 go test ./...
ok  	…/cmd/cans   …/internal/audio   …/internal/booth   …/internal/doctor   …/internal/keep
ok  	…/internal/mouth   …/internal/play   …/internal/say   …/internal/ship   …/internal/tts
exit=0

$ export CANS_HOME=<scratch>/freshhome \
         CANS_WORKER_BIN=~/.cans/native/bin/qwen3-tts-worker \
         CANS_WORKER_MODELS=~/.cans/native/models
$ <scratch>/binout/cans doctor
  machine  ok  darwin/arm64
  worker   ok  ~/.cans/native/bin/qwen3-tts-worker
  payload  ok  <scratch>/freshhome/shipped
  throat   ok  <scratch>/freshhome/shipped/voices/veronica/ref.wav
  play     ok  /usr/bin/afplay
put the cans on.
exit=0

$ <scratch>/binout/cans say "Put the cans on." -o <scratch>/fresh.wav
<scratch>/fresh.wav
( ... )  12.29s user 0.67s system 103% cpu 12.473 total
exit=0

$ ffprobe -v error -show_entries stream=codec_name,sample_rate,channels -show_entries format=duration,size …
codec_name=pcm_s16le
sample_rate=24000
channels=1
duration=1.041583
size=50040

$ ls $CANS_HOME
mouth.lock
shipped

sampler: max=1  counts seen: 0 1  (13 samples)
```

Five doctor rows ok, a real 24 kHz mono wav in 12.5 s, one worker, and a fresh home containing exactly the lock file and the unpacked payload. **PASS.**

## 10. Identity

```bash
$ git log origin/main..cans-v2 --format='%h %an <%ae>'
9136cbe Veronica <318153306+veronica-agent@users.noreply.github.com>
5e8123d Veronica <318153306+veronica-agent@users.noreply.github.com>
c443de9 Veronica <318153306+veronica-agent@users.noreply.github.com>
c736f46 Veronica <318153306+veronica-agent@users.noreply.github.com>
47d39d9 Veronica <318153306+veronica-agent@users.noreply.github.com>
0f4e436 Veronica <318153306+veronica-agent@users.noreply.github.com>

$ git log origin/main..cans-v2 --format='%cn <%ce>' | sort -u
Veronica <318153306+veronica-agent@users.noreply.github.com>

$ git log origin/main..cans-v2 --format=%B | rg -i 'co-authored|claude|gpt|grok'
(nothing — exit 1)
```

Six commits, one author and one committer on every one, no assistant attribution anywhere in the bodies. **PASS.**

## 11. fest validate (festival + snapshot)

```bash
$ fest validate                       # the campaign festival
● STRUCTURE
⚠ Task filename should match NN_name.md: .review-01_out.md
⚠ Task filename should match NN_name.md: .review-02_lock.md
✓ COMPLETENESS  ✓ Task Files  ✓ QUALITY GATES  ✓ Markers  ✓ ORDERING  ✓ AUTO-LINK  ✓ HOOKS  ✓ WORKFLOW
Score 90/100
VALIDATION PASSED WITH WARNINGS
exit=0

$ fest validate festivals/CV0001     # the public snapshot, from the worktree
✓ STRUCTURE  ✓ COMPLETENESS  ✓ Task Files  ✓ QUALITY GATES  ✓ Markers  ✓ ORDERING  ✓ AUTO-LINK  ✓ HOOKS  ✓ WORKFLOW
Score 100/100
VALIDATION PASSED
exit=0
```

Both pass. The campaign festival's two warnings are the cold reviewers' hidden `.review-01_out.md` / `.review-02_lock.md` scratch notes — the exact files D009 was amended to exclude, which is why the snapshot scores 100. Expected, not a defect. **PASS.**

## 12. Snapshot re-sync

The recorded command from `003_IMPLEMENT/05_snapshot/01_snapshot.md`, re-run after this phase's own statuses and results were written. `--delete-excluded` is load-bearing: plain `--exclude` will not remove a file already present at the destination.

```bash
$ rsync -a --delete --delete-excluded \
    --exclude CONTEXT.md --exclude '001_INGEST/input_specs' \
    --exclude .fest --exclude .workflow --exclude .festival-checksums.json \
    --exclude '.review-*' \
    festivals/active/cans-v2-CV0001/ projects/worktrees/cans/cans-v2/festivals/CV0001/
exit=0

$ find festivals/CV0001 -type f | wc -l
      86

$ fest validate festivals/CV0001
Score 100/100 — VALIDATION PASSED

$ rg -i '<pattern 1 — CONTEXT.md §Professional grep>' README.md | rg -v '<allowed footer line>'
(nothing)
$ rg -i '<pattern 2 — CONTEXT.md §Professional grep>' README.md docs/ tapes/ festivals/CV0001/
(nothing)
$ rg -c 'fest.build' README.md
1

$ diff -r -x CONTEXT.md -x input_specs -x .fest -x .workflow -x .festival-checksums.json -x '.review-*' \
    festivals/active/cans-v2-CV0001 projects/worktrees/cans/cans-v2/festivals/CV0001
(nothing — the trees are identical under the six exclusions)
```

**PASS.** Committed with this phase.

## 13. PR

PR **#14** on `veronica-agent/cans`, `cans-v2` → `main`, opened under the `veronica-agent` account. The body was rewritten from this file at the end of the phase: why the change exists, the five features with their flags and exit codes, the Attempt 4 numbers with a citation to `festivals/CV0001/002_PLAN/inputs/measurements.md`, one line per bar item, the two known engine behaviours, and what is deliberately not in the PR. The three professional-surface patterns were run over the body before it was applied and all three printed nothing.

CI and the final head SHA are recorded in the phase report.

## Carry-over from `04_tape` — `docs/pipe.gif` re-cut

`04_tape` cut the gif at 1-minute load 16.3–27.9 and assigned the re-cut here. Re-cut on a quiet box with no change to the tape:

```bash
$ just vhs pipe
Creating docs/pipe.gif...
( ... )  18.94s user 6.36s system 23% cpu 1:46.40 total

load during the record: n=22  mean=4.89  min=4.54  max=5.35  below16=22 (100%)

$ ffprobe -v error -show_entries format=duration,size -show_entries stream=width,height …
width=680
height=300
duration=91.520000
size=148556
```

|  | committed by `04_tape` | re-cut here |
|---|---|---|
| Geometry | 680 × **380** | 680 × **300** — now matches `tapes/pipe.tape:13` |
| Duration | 123.44 s | **91.52 s** |
| Size | 168 370 B | **148 556 B** (0.14 MB, bar is 2 MB) |
| `ttfa_ms` on screen | 31377 / 34570 / 35202 | **26607 / 29461 / 14174** |
| Load while recording | 16.26–27.93 | 4.54–5.35 |

The last frame was extracted (`ffmpeg -y -sseof -0.5 -i docs/pipe.gif -frames:v 1`) and read: the three JSON records, then `ls out` returning `001.wav 002.wav 003.wav`, then the prompt — the whole run fits the 300 px frame with no dead space, which is what the height change was for.

The on-screen `ttfa_ms` are still 14–29 s for three short lines on an idle box. That is the mouth's end-of-speech variance, not load and not `--stream`: the same binary in item 6 did a one-shot in 7 157 ms and in item 5 in 5 627 ms. The gif is honest about what the field means — total synthesis time — and the README and the PR body both say so.

## Known engine behaviour seen during this bar

Neither is a defect in this branch; both are already deferred in `CONTEXT.md`, and both are named in the PR body so a stranger reading the numbers is not surprised.

1. **End-of-speech variance.** The same one-shot text cost 5 627 / 5 652 / 5 669 / 5 798 / 7 157 ms in some runs and 27 241 ms in another, and stream lines ran 32–42 s. `ttfa_ms` is stamped at `final`, so it is total synthesis time; the variance is the worker's, not the caller's.
2. **Near-silent wavs.** 1 484-byte, ~0.03 s files returned as success after full synthesis time — one in item 1 (1 of 3), one in item 2 (1 of 8), one in item 3 (1 of 2). Matches the 4/50, 5/50, 1/24 rates measured in all three modes in `04_measure`, which is what identifies it as the mouth rather than `--stream`. A length floor that warns on a suspiciously short return is the cheap future guard.
