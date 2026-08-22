# Measurements — the real mouth on this Mac

Machine: Apple M4 Max, 16 cores, 128 GB, macOS 26.5. Worker: `~/.cans/native/bin/qwen3-tts-worker` (cans overlay, 2026-08-20) with `~/.cans/native/models`. Ref: `voices/veronica/ref.wav`. Text: `Put the cans on.` (4 words → cans sends `max_tokens 220`, `temperature 0.2`).

**Rule (D013): no number is quoted without its command, and nothing else may be running.** Two workers on Metal at once inflated an early probe by 8×; a 1-minute load average of 282 (other sessions on the box) inflated a later one by 10×. Check `uptime` first; 1-minute load must be below 16.

## Baseline — worker alone, idle machine (2026-08-21 ~04:20, load < 16)

Probe: a 60-line Go program (`scratchpad/measure/main.go`) that starts the worker, stamps `ready`, sends one `synthesize` with cans' exact request fields, stamps `final`, sends `shutdown`, and reads the worker's `Rusage.Maxrss`. The probe was throwaway scaffolding and lives campaign-side, not in this repo; the §Stream numbers below are the reproducible ones, taken with the shipped `cans` binary.

```
go run . ~/.cans/native/bin/qwen3-tts-worker ~/.cans/native/models <ref.wav>
```

| Run | ready (GGUF load) | synth | audio | worker max RSS |
|-----|-------------------|-------|-------|----------------|
| 1 | 7 635 ms | 29 932 ms | 17.58 s | 3 758 MB |
| 2 | 6 517 ms | 6 185 ms | 1.90 s | 2 853 MB |
| 3 | 6 568 ms | 36 760 ms | 17.58 s | 3 761 MB |

Without `max_tokens`/`temperature` (worker defaults): ready 6 626–6 762 ms, 28.78 s of audio every time (the 360-token ceiling), RSS 4 397 MB.

**Read:** GGUF load is **~6.6 s** and is the cost a loop pays per call. Resident memory for one worker is **2.8–4.4 GB** depending on how long it generates — eight of them is 22–35 GB of weights. Run 2 is the intended case (stops at end of speech, 1.9 s of audio for four words); runs 1 and 3 ran to the 220-token budget (17.58 s). That variance is **pre-existing mouth behavior** (flagged for the operator in `CONTEXT.md`), not v2's, and it is why stream measurements report median and max.

## Baseline — cold `cans say`, idle machine

```
cd <the drop-sidecar worktree — campaign-side, pre-dates this branch> && just build quick
CANS_NOPLAY=1 /usr/bin/time -l ./bin/cans say "Put the cans on."
```

| ttfa_ms (cans' field = total synth wall) | real | note |
|------------------------------------------|------|------|
| 5 839 | 13.1 s | 04:21, load < 16, run 2-style end of speech |

So a one-shot `cans say` is **~6.6 s load + ~6 s synth ≈ 13 s** cold for four words, of which the load is what `--stream` removes from every line after the first.

## Contaminated runs (kept so nobody repeats the mistake)

- 04:35 — probe and `cans say` launched **concurrently**: `cans say` real 106–120 s, `ttfa_ms` 97 524–110 290; probe synth 30–77 s. Two workers on Metal.
- 04:45 — "alone" but 1-minute load average **282** (35 sessions on the box): ready 11–32 s, synth 16–77 s for 1.5–2.2 s of audio.

## Stream — filled by `003_IMPLEMENT/03_stream/04_measure`

**Status: COMPLETE — taken on attempt 4, 2026-08-21 16:24–17:45.** Four attempts. Attempts 1 and 2 deferred before starting; attempt 3 got run (a) only, on a box averaging load 23. **Attempt 4 ran all three — (a) stream, (b) loop, (c) `xargs -P 8` — back to back on a genuinely quiet machine, with 100 % of every run's load samples below 16, so the margin is measured and uncontaminated: `loop wall − stream wall` = **596.7 s over 50 lines, 11.93 s per line**, of which **6.44 s per line is the structural overhead `--stream` removes**. Attempt 4 is the quotable attempt. Everything above it is kept as history so nobody repeats the mistakes.

Read the split below carefully. Attempt 3 produced two kinds of result, and only one kind is quotable:

- **Structural results are solid.** Worker count, record count, index density and error count do not depend on machine load. These are the claims `--stream` actually makes, and they hold.
- **Every wall clock and every `ttfa_ms` from attempt 3 is contaminated** and must not reach the README. During the 31-minute run the 1-minute load averaged **23.2** and peaked at **92.6**; only 56 % of samples were under the D013 bar of 16.

### Attempts 1 and 2 — deferred before starting

- **Attempt 1 — deferred: load 39.21 at 2026-08-21 07:02.** 1-minute load 39.21 (must be < 16); a `festival-voice` worker (PID 62136, not cans) resident.
- **Attempt 2 — deferred: load 17.38 at 2026-08-21 14:41**, after a 15-minute wait for a window that never opened. 164 load samples over 14 minutes: mean **32.9**, min 15.84, max 60.04, **only 0.6 % below 16**. Source was two foreign VMs holding 841 % CPU and 25.6 GB:

```
$ ps -Ao pcpu,rss,pid,comm -r | head -3
 %CPU      RSS   PID COMM
479.7  5587568 74004 …/com.apple.Virtualization.VirtualMachine
361.5 20017296 13387 …/com.apple.Virtualization.VirtualMachine
$ ps -Ao pcpu -r | awk 'NR>1{s+=$1} END{print s"%"}'
1032.6%
```

### Attempt 3 — run (a) completed, runs (b) and (c) not taken

Started 14:57 when the box briefly quieted; preconditions passed at **14:58:07 with load 15.41** and no cans worker resident.

```
$ for i in $(seq 1 50); do echo "Measurement line $i. The worker stays warm."; done > lines.txt
$ while sleep 1; do pgrep -f 'cans/native/bin/qwen3-tts-worker' | wc -l | tr -d ' '; done > a.pgrep &
$ time ( $CANS say --stream -o "out/stream/%03d.wav" --json < lines.txt > stream.jsonl )
```

#### (a) Stream, 50 lines — structural results (quotable)

| Property | Result | Why it is load-independent |
|----------|--------|----------------------------|
| Exit code | **0** | — |
| Records emitted | **50 / 50** | one per stdin line |
| `line` values | **exactly 1…50, dense** (`jq -s 'map(.line) == [range(1;51)]'` → `true`) | D005/D006 |
| Error records | **0** | — |
| Wav files written | **50 / 50** on the `%03d` template | — |
| **Workers resident, max** | **1** | 1 753 one-second samples, distinct values `{0, 1}` — **never 2** |
| GGUF loads | **1**, amortised across 50 lines | see overhead below |

```
$ sort -n a.pgrep | tail -1
1
$ sort -u a.pgrep | tr '\n' ' '
0 1
$ wc -l < a.pgrep
1753
$ pgrep -fl 'cans/native/bin/qwen3-tts-worker'   # between runs
(exit=1)
```

**This is the result the sequence exists to prove.** A 50-line document went through one `Session`, one `flock`, one worker process and one GGUF load, and the worker was gone the moment the stream ended.

The overhead number is the strongest evidence, and it survives the contamination because it is a *ratio* of two quantities measured on the same loaded box:

```
$ jq -s 'map(.ttfa_ms) | add / 1000' stream.jsonl
1855.906
```

**1 855.9 s of reported synthesis inside a 1 877 s wall — 21.1 s of everything else, for all 50 lines.** Process start, `doctor.Prepare`, `keep.Load`, lock acquire, one GGUF load, 50 file writes and 50 flushed records together cost about 21 seconds. 98.9 % of the wall was the mouth. A 50-call loop would have paid the GGUF load 50 times instead of once.

#### (a) Stream — timing (CONTAMINATED, do not quote)

| Metric | Value | |
|--------|-------|--|
| Wall | 1 877 s (`real 31m16.462s`, `user 32m45.870s`, `sys 3m32.738s`) | contaminated |
| `ttfa_ms` min | 9 143 | contaminated |
| `ttfa_ms` q1 / **median** / q3 | 35 688 / **37 151** / 41 255 | contaminated |
| `ttfa_ms` max | **105 278** | contaminated |
| Lines over 20 s | **41 of 50** | contaminated |
| Pageouts | 866 170 → 869 796, **delta 3 626** | contaminated; not the (c) check |

```
$ jq -s 'map(.ttfa_ms) | sort | .[length/2|floor], max' stream.jsonl
37151
105278
```

Load traced once per 5 s for the whole run:

```
$ awk '$1>="14:58:07" && $1<="15:29:30" {n++; v=$2+0; s+=v; if(v<16)k++; if(v>mx)mx=v} \
    END {printf "samples=%d mean=%.1f max=%.1f below16=%d (%.0f%%)\n", n,s/n,mx,k,(k/n)*100}' loadtrace2.log
samples=373  mean=23.2  max=92.6  below16=208 (56%)
```

A median of 37 s per line is not this machine's behaviour. The same binary, minutes earlier on a quiet box, did the three-line manual check at `ttfa_ms` 6 122 / 27 348 / 24 215 and a one-shot at 5 652. The pageouts delta of 3 626 is likewise the two foreign VMs paging, not cans — the real pageouts check belongs to run (c), which did not happen.

#### (b) and (c) — not taken

```
ABORT after A: load 34.28 >= 16, will not measure B on a loaded box
```

Run (b) waited its full 5-minute settling window at load 31–34 and gave up; run (c) never started. **Without (b) there is no margin and no per-line saved cost.** The value can be *predicted* from the baseline — one GGUF load is ~6.6 s, so 49 avoided loads ≈ 323 s ≈ **6.5 s/line** — but that is arithmetic on an old measurement, not a measured margin, and it must not be quoted as one.

### What a valid run still needed (written before attempt 4; attempt 4 met all three)

1. `uptime` 1-minute load below 16 **and staying there** for roughly 60–90 minutes. Attempt 3 proves a momentary dip is not enough: it passed the precondition at 15.41 and was at 92.6 twenty minutes later.
2. All three runs — (a) stream, (b) 50-call loop, (c) `xargs -P 8` over 24 lines — back to back, with `pgrep` empty between them.
3. `vm_stat` pageouts delta of 0 across run (c).

The `xargs` quoting is solved; BSD `xargs` has no `-d`, so the pipeline is null-delimited. Verified against a stand-in before it was pointed at the mouth:

```
$ head -24 lines.txt | nl -ba | sed 's/^[[:space:]]*//' | tr '\t\n' '\0\0' \
    | xargs -0 -P 8 -n 2 sh -c '"$CANS" say "$2" -o "$XOUT/$1.wav" --json' _ > x.jsonl 2> x.err
```

### Attempt 4 — 2026-08-21 16:24–17:45 — **all three runs completed on a quiet box**

Preconditions passed at **16:27:01 with 1-minute load 9.22** and no cans worker resident. The box stayed quiet for the whole 78 minutes: across the three runs' own 30-second load samples, **every single sample was below 16** (mean 7.1–8.4, max 14.97). This is the sustained window attempt 3 never got.

Binary: `projects/worktrees/cans/cans-v2/bin/cans` at `c443de9` (plus uncommitted README/tape changes) — **not rebuilt** for this attempt. Scratch dir outside the repo: `/private/tmp/…/scratchpad/measure4.bBtMcq`.

#### The sampler trap that produced a false "2 workers" reading

The sampler this task file prescribes — `pgrep -f 'cans/native/bin/qwen3-tts-worker' | wc -l` — **cross-matches other samplers**. A concurrent sampler's own `pgrep` process has the pattern in its argv, so sampler A counts sampler B's `pgrep` as a worker. Mid-run (a) this drove the reading to 2 and then 3 while exactly one worker was resident (393 samples read 2, 63 read 3).

It is not the self-match the earlier note warns about (an inline `while` loop) — a separate script file does not fix it. Proof it was an artifact, taken at the same instants:

```
$ ./diag2.sh              # logs `pgrep -fl` (identities, not a count) once a second, 40 blocks
$ grep -v '^===' diag2.txt | grep -v '\.cans/native/bin/qwen3-tts-worker …/models$'
(nothing — every block contained the one worker and nothing else)
$ tail -20 a.pgrep   ; # old counter, same seconds
2 2 3 3 3 3 3 3 3 3 2 2 2 3 3 2 3 2 3 3
$ cat a2.pgrep       ; # anchored ps counter, same seconds
1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
```

The fix used for (b), (c) and the parallel check on (a) — anchored on the absolute path, so neither its own `grep` nor any concurrent sampler's `pgrep`/`grep` can be counted, and the foreign `~/.cache/festival-voice/…/qwen3-tts-worker` is excluded by path:

```sh
ps -Ao command= | grep -c '^~/\.cans/native/bin/qwen3-tts-worker '
```

**Every "worker max" below is from that counter.** It never read above 1 in any run.

#### The three runs

| | (a) stream, 50 lines | (b) loop, 50 calls | (c) `xargs -P 8`, 24 lines |
|---|---|---|---|
| Started | 16:27:01 | 16:52:16 | 17:29:02 |
| Launch load | 9.22 | 9.11 | 8.13 |
| **Wall** | **1 512.6 s** (`real 25m12.556s`) | **2 109.2 s** | **943.2 s** (`real 15m43.217s`) |
| `user` / `sys` | 27m28.811s / 2m40.714s | — (see note) | 16m40.733s / 1m41.119s |
| Exit code | 0 | 0 | 0 |
| Records | **50 / 50** | **50 / 50** | **24 / 24** |
| Error records | 0 | 0 | 0 |
| Wavs written | 50 | 50 | 24 (`001`…`024`) |
| `line` dense 1…50 | **true** | n/a | n/a |
| `ttfa_ms` min | 8 637 | 10 297 | 9 146 |
| `ttfa_ms` q1 / **median** / q3 | 23 099 / **35 730** / 36 057 | 35 638 / **35 842** / 36 067 | — / **35 736** / — |
| `ttfa_ms` max | **37 614** | **77 932** | **36 714** |
| Σ `ttfa_ms` | 1 495.8 s | 1 770.6 s | 782.7 s |
| Lines over 20 s | 39 / 50 | 45 / 50 | 21 / 24 |
| **Worker max** | **1** (430 samples, `{0,1}`) | **1** (1 937 + 44 samples, `{0,1}`) | **1** (878 samples, `{0,1}`) |
| Load samples | n=51 mean **8.36** min 5.55 max 14.97 — **100 % < 16** | n=70 mean **7.06** min 5.19 max 9.70 — **100 % < 16** | n=32 mean **7.45** min 5.75 max 11.10 — **100 % < 16** |
| Pageouts | 872 087 → 872 724 (delta 637) | 872 724 → 875 724 (delta 3 000) | 875 724 → 875 724 — **delta 0** |
| stderr | empty | see note | 20 × `waiting for the mouth…`, nothing else |
| Worker after | gone | gone | gone |

**No run is contaminated.** Every load sample of all three runs was under the D013 bar of 16.

#### The margin — measured, clean

```
margin        = loop wall − stream wall = 2 109.2 − 1 512.6 = 596.7 s
per-line cost = 596.7 / 50               = 11.93 s per line
```

**`--stream` saved 596.7 s on a 50-line document — 11.93 s per line — and the number is clean** (both runs ran with 100 % of load samples below 16, mean 8.4 and 7.1).

#### What the margin is made of

The margin splits cleanly into the part `--stream` removes structurally and the part that is the mouth being the mouth:

| Component | (a) stream | (b) loop | Difference |
|---|---|---|---|
| Wall | 1 512.6 s | 2 109.2 s | 596.7 s |
| Σ reported synthesis (`ttfa_ms`) | 1 495.8 s | 1 770.6 s | 274.9 s |
| **Everything else** (process start, `doctor.Prepare`, `keep.Load`, lock, **GGUF load**, writes, records) | **16.8 s total — 0.34 s/line** | **338.6 s total — 6.77 s/line** | **321.8 s — 6.44 s/line** |

- **98.9 % of the stream run was the mouth.** For all 50 lines, everything that is not synthesis cost **16.8 seconds**.
- The loop paid **6.77 s per line** of that same overhead, and run (c) independently paid **6.69 s per line** (160.5 s of non-synthesis over 24 calls). Both reproduce the **~6.6 s GGUF load** measured in the baseline at the top of this file, from two different directions.
- **The defensible structural claim is therefore ~6.4 s per line** (321.8 s over 50), and that is the number the README should use if it quotes one. The remaining 274.9 s of the 596.7 s margin is the loop's *higher reported synthesis time* — a fresh worker per call, subject to the mouth's end-of-speech variance — which is real but is not something `--stream` engineered away.

#### (c) `xargs -P 8` — the lock does its job

Twenty-four `cans say` calls were launched eight at a time. Sampled mid-run:

```
$ ps -Ao command= | grep -c '^…/cans-v2/bin/cans say '
8
$ head -5 x.err
waiting for the mouth…
waiting for the mouth…
waiting for the mouth…
waiting for the mouth…
waiting for the mouth…
$ sort -n c.pgrep | tail -1
1
```

**Eight concurrent `cans` processes, seven of them blocked on the lock, one worker resident, 24/24 wavs, zero errors, and a pageouts delta of exactly 0.** `-P 8` cannot start a second worker and cannot make the machine swap — which is the whole reason D001–D003 exist. The cost is that `-P 8` buys nothing: 943 s for 24 lines is 39.3 s/line, the same serial rate as the loop.

#### Near-silent outputs (the known mouth fault, not v2's)

Wavs under 2 000 bytes — 1 484 bytes each, ~0.03 s of near-silence, reported as success:

| Run | Indices | Count |
|---|---|---|
| (a) stream | **002, 033, 036, 048** | 4 / 50 (8 %) |
| (b) loop | **001, 006, 011, 031, 044** | 5 / 50 (10 %) |
| (c) xargs | **014** | 1 / 24 (4 %) |

```
$ jq -c 'select(.line==2 or .line==33 or .line==36 or .line==48)|{line,ttfa_ms}' stream.jsonl
{"line":2,"ttfa_ms":35489}
{"line":33,"ttfa_ms":35730}
{"line":36,"ttfa_ms":35882}
{"line":48,"ttfa_ms":36197}
```

Same shape as attempt 3: **35–36 s of synthesis to produce 0.03 s of audio**, at the run's median cost, with the worker reporting success. It appears in all three modes at ~4–10 % of lines, so it is independent of `--stream` — a **mouth** fault, already recorded in `CONTEXT.md` for the operator. The largest wav in the same stream run is `016.wav` at 1 158 614 bytes (~24 s).

#### Exact commands

```bash
D=$(mktemp -d …/scratchpad/measure4.XXXXXX); mkdir -p "$D"/out/{stream,loop,x}; cd "$D"
CANS=~/Dev/AI/veronica-campaign/projects/worktrees/cans/cans-v2/bin/cans

for i in $(seq 1 50); do echo "Measurement line $i. The worker stays warm."; done > lines.txt

# sampler (separate script file; anchored ps counter — see the trap above)
cat > pgrepsample.sh <<'EOF'
#!/bin/sh
while :; do
  ps -Ao command= | grep -c '^~/\.cans/native/bin/qwen3-tts-worker ' >> "$1"
  sleep 1
done
EOF
# load sampler, every 30 s
cat > loadsample.sh <<'EOF'
#!/bin/sh
while :; do
  printf '%s %s\n' "$(date +%H:%M:%S)" \
    "$(uptime | sed 's/.*load averages*: *//' | awk '{print $1}' | tr -d ',')" >> "$1"
  sleep 30
done
EOF

# before each run: uptime (1-min < 16) and no cans worker resident
uptime; pgrep -fl 'cans/native/bin/qwen3-tts-worker'; vm_stat | grep Pageouts

# (a) stream
./pgrepsample.sh "$D/a.pgrep" & ./loadsample.sh "$D/a.load" &
{ time ( "$CANS" say --stream -o "$D/out/stream/%03d.wav" --json \
    < "$D/lines.txt" > "$D/stream.jsonl" 2> "$D/stream.err" ) ; } 2> time_a.txt

# (b) loop
./pgrepsample.sh "$D/b.pgrep" & ./loadsample.sh "$D/b.load" &
{ time ( i=0; while IFS= read -r l; do i=$((i+1));
    "$CANS" say "$l" -o "$D/out/loop/$(printf %03d $i).wav" --json;
  done < "$D/lines.txt" > "$D/loop.jsonl" 2> "$D/loop.err" ) ; } 2> time_b.txt

# (c) xargs -P 8 over the first 24 lines (BSD xargs has no -d, so null-delimited)
export CANS XOUT="$D/out/x"
{ time ( head -24 "$D/lines.txt" | nl -ba | sed 's/^[[:space:]]*//' | tr '\t\n' '\0\0' \
    | xargs -0 -P 8 -n 2 sh -c '"$CANS" say "$2" -o "$XOUT/$(printf %03d "$1").wav" --json' _ \
    > "$D/x.jsonl" 2> "$D/x.err" ) ; } 2> time_c.txt

# analysis
jq -s 'map(.ttfa_ms) | sort | .[length/2|floor], max' stream.jsonl
jq -s 'map(.ttfa_ms) | add / 1000'                    stream.jsonl
jq -s 'map(select(.error!=null)) | length'            stream.jsonl
jq -s 'map(.line) == [range(1;51)]'                   stream.jsonl
find out/stream -name '*.wav' -size -2000c -exec basename {} \; | sort
sort -n a2.pgrep | tail -1
```

#### One honest caveat about run (b)'s wall

The agent harness killed the background shell holding run (b) after 60 minutes of wall time, **during call 50** — `loop.err` contains one `say: interrupted` from that kill, and no worker leaked. Calls 1–49 had already completed and were timed exactly; call 50 was re-run alone 90 seconds later on the same idle box and its wall added:

```
calls 1–49 : 2 062 s   (b.t0 epoch → mtime of loop.jsonl at record 49; ±1 s)
call 50    :    47.229 s (`real 0m47.229s`, load 6.90 → 8.54, worker max 1)
loop wall  : 2 109.2 s
```

The loop is 50 independent one-shot invocations, so the sum is the same quantity a single uninterrupted `time` would have printed, to within the ±1 s of the mtime read. This is why (b) has no `user`/`sys` figures. Every other number in run (b) — records, errors, wavs, `ttfa_ms`, worker max, load — is from the 50 completed calls.

Total attempt time was 78 minutes against a 75-minute budget; run (c) started at 17:29:02, inside the budget, and was allowed to finish.

### Finding for the reviewer — intermittent near-silent wavs

**3 of the 50 stream wavs are 1 484 bytes**, about **0.03 s** of near-silence, for ordinary inputs — `026.wav`, `034.wav`, `048.wav` — while `004.wav` in the same run is 1 145 808 bytes (~24 s). The same thing appeared in the cancel check (`002.wav`, 1 484 bytes, input `Line 2`).

The important part is what those three lines *cost*:

```
$ find out/stream -name '*.wav' -size -5000c -exec basename {} \; | sort
026.wav  034.wav  048.wav
$ jq -c 'select(.line==26 or .line==34 or .line==48) | {line, ttfa_ms}' stream.jsonl
{"line":26,"ttfa_ms":36501}
{"line":34,"ttfa_ms":40391}
{"line":48,"ttfa_ms":54272}
```

**36.5 s, 40.4 s and 54.3 s of synthesis to produce 0.03 s of audio each.** So this is *not* an early end-of-speech — the worker ran a long generation and emitted almost nothing at the end of it. Line 48 was the third most expensive line in the whole run and returned silence.

The files themselves are well-formed WAV — PCM, mono, 24 000 Hz, 16-bit, `data` chunk 1 440 bytes, leading samples zero:

```
$ xxd -l 48 out/stream/048.wav
00000000: 5249 4646 c405 0000 5741 5645 666d 7420  RIFF....WAVEfmt
00000010: 1000 0000 0100 0100 c05d 0000 80bb 0000  .........]......
00000020: 0200 1000 6461 7461 a005 0000 0000 0000  ....data........
```

`--stream` did the right thing at every step: it wrote what the worker returned, emitted a success record with the true `ttfa_ms`, and continued. Nothing in the write path or the record path is wrong, and no error was available for it to report — the worker reported success.

This is a **mouth** fault, not a `03_stream` defect, and it is outside this festival's scope to fix. But it is worth an operator decision, because at roughly **6 % of lines** a 200-line render would silently lose about a dozen lines *and pay full synthesis time for each one*. A caller cannot currently tell these apart from real output without inspecting wav length — which suggests a cheap future guard: warn when a returned wav is under some floor.

### Not cans' worker (recorded once, per D013)

```
$ ps -o pid,ppid,pcpu,rss,etime,command -p 62136
  PID  PPID  %CPU    RSS  ELAPSED COMMAND
62136 61770   0.1  24688 08:56:49 ~/.cache/festival-voice/models/qwen3-tts/bin/qwen3-tts-worker …/models
```

24 MB resident, 0.1 % CPU, idle 8 h 56 m — no model loaded. It is `festival-voice`, not `~/.cans/native/bin/qwen3-tts-worker`, and it must not be killed. Because a bare `pgrep -f qwen3-tts-worker` matches it, **every cans check must use the path-qualified pattern** `pgrep -f 'cans/native/bin/qwen3-tts-worker'`.

### Clean real-mouth numbers taken on the quiet window (functional checks, not the measurement)

From the `05_testing` gate at 14:54–14:57, load **9.9–14.7**, i.e. genuinely under the bar. Not a substitute for run (a)/(b)/(c), but the only uncontaminated mouth numbers this session produced:

| Check | Load | Result |
|-------|------|--------|
| One-shot `cans say "Put the cans on."` | 9.90 → 10.25 | `ttfa_ms=5652`, exit 0, 14 s wall, pgrep max 1 — matches the 5 839 ms / 13.1 s baseline, so one-shot is unchanged |
| 3-line `--stream --json` | 10.02 → 11.11 | exit 0, 65 s, 3 records, 3 wavs, **pgrep max 1 over 61 samples**, stderr empty; `ttfa_ms` 6 122 / 27 348 / 24 215 |
| 20-line stream, SIGINT after 002.wav | 10.48 → 14.66 | exit **130**, `interrupted after line 2`, 001+002 kept, no 003, worker gone |
| one-shot straight after the interrupt | 14.66 | exit 0, **14 s**, no `waiting for the mouth…` — lock already released |
