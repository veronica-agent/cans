---
fest_type: festival
fest_id: CV0001
fest_name: cans-v2
fest_status: active
fest_created: 2026-08-21T04:32:49.124639-06:00
fest_updated: 2026-08-21T05:14:26.058302-06:00
fest_tracking: true
---



# cans-v2

**Status:** Planning | **Created:** 2026-08-21

## Festival Objective

**Primary Goal:** Cans becomes a unix primitive a script can drive over a document: text in from argv or stdin, wav out where you point it, and one mouth at a time no matter how the script loops.

**Vision:** A shell loop over a chapter renders one wav per line with a single model load, and `xargs -P 8` cannot put two workers in memory. `cans say "x"` is byte-for-byte what it was. The festival tree lands in the public repo as the second readable plan, and nothing on that surface stops being professional.

## Success Criteria

### Functional Success

- [ ] `cans say "x" -o take.wav` writes the wav, does not play, does not delete; stdout is the path
- [ ] `echo x | cans say` and `cans say -` read one utterance from stdin; empty argv on a TTY is still exit 2
- [ ] `cat lines | cans say --stream -o 'out/%03d.wav'` writes one wav per line over **one** `Session`
- [ ] `--json` emits one record per utterance on stdout, flushed as each finishes; prose is on stderr
- [ ] A mouth lock on `CANS_HOME/mouth.lock` keeps exactly one `qwen3-tts-worker` resident across any loop, `xargs -P`, or second terminal; `--nowait` exits 75, `--wait` bounds the block
- [ ] The booth holds the lock for its session; a script started alongside it waits (or gets 75), never talks over it
- [ ] Ctrl-C mid-stream keeps completed wavs, leaves no worker and no held lock; `kill -9` leaves the next run unblocked
- [ ] A second VHS tape shows a script piping lines in and wavs landing on disk

### Quality Success

- [ ] `cans say "x"` behaves exactly as at `1e8cea2`; every existing test passes
- [ ] Stream and lock paths are tested with the fake worker (`internal/tts/testdata/fakeworker`) — CI needs no real mouth
- [ ] 200-line stream beats the 200-call loop by a recorded margin; both numbers are in this festival
- [ ] `xargs -P 8` over 50 lines: one worker at every `pgrep` sample, no swap
- [ ] `just test unit`, `go vet`, `gofmt -l` clean; files under 500 lines, functions under 50
- [ ] README, docs, tapes and the festival snapshot pass the professional-surface grep; exactly one Festival footer
- [ ] `fest validate` green on this festival

## Progress Tracking

### Phase Completion

- [ ] 001_INGEST: design pack + user direction structured into output_specs
- [ ] 002_PLAN: decisions (lock lifetime, one PR, flag grammar), measurements, STRUCTURE, IMPLEMENTATION_PLAN, scaffolded sequences
- [ ] 003_IMPLEMENT: 01_out, 02_lock, 03_stream, 04_tape, 05_snapshot — on branch `cans-v2`, committed through `fest commit`
- [ ] 004_REVIEW: the ship bar, identity, professional surface, `fest validate`

## Complete When

- [ ] All phases completed
- [ ] One PR from `cans-v2` to `main` on `veronica-agent/cans` is open with CI green and the bar recorded in its body
- [ ] `projects/cans/festivals/CV0001/` holds the readable snapshot of this festival
- [ ] Git log authors on the branch are Veronica / `318153306+veronica-agent@users.noreply.github.com`