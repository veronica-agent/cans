---
fest_type: task
fest_id: 02_recheck.md
fest_name: recheck
fest_parent: 05_snapshot
fest_order: 2
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:57.277515-06:00
fest_updated: 2026-08-21T17:58:18.492058-06:00
fest_tracking: true
---


# Task: recheck

## Objective

The whole surface, once more, before review: grep, footer, tests, build hygiene, fresh-home doctor, and `cans` without `fest`.

## Requirements

- [x] Professional-surface grep (`CONTEXT.md §Professional grep`) over `README.md`, `docs/`, `tapes/`, `festivals/` — empty. `rg -c 'fest.build' README.md` → `1`.
- [x] `gofmt -l .` empty; `go vet ./...`; `CANS_NOPLAY=1 go test ./...` green; `git diff origin/main -- go.mod go.sum` empty; every `.go` file < 500 lines (`wc -l $(git ls-files '*.go') | sort -n | tail -3`).
- [x] Fresh home: `tmp=$(mktemp -d); cp bin/cans $tmp/; cd $tmp; CANS_HOME=$tmp/home CANS_WORKER_BIN=$HOME/.cans/native/bin/qwen3-tts-worker CANS_WORKER_MODELS=$HOME/.cans/native/models ./cans doctor` → all ok; then `./cans say -o $tmp/take.wav "Put the cans on."` → prints the path, file has a valid header (`ffprobe` or `afplay`); `ls $tmp/home` shows `mouth.lock` and `shipped/`, nothing else.
- [x] `cans` without `fest`: `PATH=/usr/bin:/bin ./bin/cans version` works (no runtime dependency on the `fest` binary — trivially true, but recorded).
- [x] `just --list` still lists `vhs pipe`; `just dist check` passes (`goreleaser check`).

## Implementation

Run each block, paste the output into this task file under **Results**, fix anything red in place (small fixes only — anything larger goes back to the owning sequence as an iterate finding).

## Results

All from the `cans-v2` worktree. Box was quiet for every real-mouth step: load
`7.39 8.72 8.36` at 17:56 (all < 16) and `pgrep -fl 'cans/native/bin/qwen3-tts-worker'` empty
before each one; steps run one at a time.

### 1. Professional-surface grep

```
$ rg -i '<pattern 1 — CONTEXT.md §Professional grep>' README.md | rg -v '<allowed footer line>'
(no output — exit 1)

$ rg -i '<pattern 2 — CONTEXT.md §Professional grep>' README.md docs/ tapes/ festivals/
… 24 hits, 11 files, EVERY ONE under festivals/CA0001/ …
exit=0

$ rg -i '<pattern 2>' festivals/CV0001/
(no output — exit 1)          # this festival's snapshot is clean

$ rg -i -l '<pattern 2>' README.md docs/ tapes/ festivals/ | sed 's|/.*||' | sort | uniq -c
  11 festivals               # zero in README.md, docs/, tapes/

$ rg -c 'fest.build' README.md
1
```

**`README.md`, `docs/`, `tapes/` and `festivals/CV0001/` are clean.** The only hits in the whole
public tree are pre-existing, committed content in **`festivals/CA0001/`** — the previous
festival's snapshot — and they are **not this sequence's to fix** (see the finding below).

### 2. Build hygiene

```
$ gofmt -l .
(no output — exit 0)

$ go vet ./...
(no output — exit 0)

$ CANS_NOPLAY=1 go test ./...
ok  	github.com/veronica-agent/cans/cmd/cans	0.610s
ok  	github.com/veronica-agent/cans/internal/audio	(cached)
ok  	github.com/veronica-agent/cans/internal/booth	(cached)
ok  	github.com/veronica-agent/cans/internal/doctor	(cached)
ok  	github.com/veronica-agent/cans/internal/keep	(cached)
ok  	github.com/veronica-agent/cans/internal/mouth	(cached)
ok  	github.com/veronica-agent/cans/internal/play	(cached)
ok  	github.com/veronica-agent/cans/internal/say	(cached)
ok  	github.com/veronica-agent/cans/internal/ship	(cached)
ok  	github.com/veronica-agent/cans/internal/tts	(cached)
exit=0                        # all ten packages

$ git diff origin/main -- go.mod go.sum
(empty — no new dependencies)

$ wc -l $(git ls-files '*.go') | sort -n | tail -3
     352 internal/say/stream_test.go
     418 internal/say/say_test.go
    5360 total
# five largest single files: 201 internal/tts/worker.go, 217 internal/ship/ship_test.go,
# 236 internal/tts/synth_test.go, 352 internal/say/stream_test.go, 418 internal/say/say_test.go
# every file < 500. worker.go is 201, still 5 over the 196 the rules pin — a reviewer's note
# from 03_stream, unchanged by this sequence.
```

### 3. Fresh home (real mouth)

Built first, then copied out of the worktree so nothing in `~/.cans` could be picked up
implicitly. `$tmp` is a scratch dir; `CANS_HOME` is `$tmp/home`, which does not exist at the
start.

```
$ just build quick
go build -trimpath -ldflags "-s -w -X …/internal/ship.Version=v0.1.0-27-g5e8123d" -o bin/cans ./cmd/cans
# clean tag, no modified-tree suffix

$ tmp=<scratch>/freshhome; rm -rf $tmp; mkdir -p $tmp; cp bin/cans $tmp/; cd $tmp
$ CANS_HOME=$tmp/home \
  CANS_WORKER_BIN=$HOME/.cans/native/bin/qwen3-tts-worker \
  CANS_WORKER_MODELS=$HOME/.cans/native/models ./cans doctor
  machine  ok  darwin/arm64
  worker   ok  ~/.cans/native/bin/qwen3-tts-worker
  payload  ok  <scratch>/freshhome/home/shipped
  throat   ok  <scratch>/freshhome/home/shipped/voices/veronica/ref.wav
  play     ok  /usr/bin/afplay
put the cans on.
exit=0                        # five rows, all ok; payload and throat unpacked into the new home

$ ./cans say -o $tmp/take.wav "Put the cans on."      # same env as above
<scratch>/freshhome/take.wav
exit=0  wall=14s              # prints the path, nothing else on stdout

$ ffprobe -v error -show_entries format=format_name,duration,size \
      -show_entries stream=codec_name,sample_rate,channels -of default=nw=1 take.wav
codec_name=pcm_s16le
sample_rate=24000
channels=1
format_name=wav
duration=1.087542
size=52246
exit=0                        # valid header, real audio — not the 1 484-byte near-silent fault

$ ls -A $CANS_HOME
mouth.lock
shipped
                              # exactly the two, nothing else
```

### 4. `cans` without `fest`

```
$ env -i PATH=/usr/bin:/bin HOME=$HOME ./bin/cans version
cans v0.1.0-27-g5e8123d
exit=0

$ env -i PATH=/usr/bin:/bin sh -c 'command -v fest || echo "fest not on PATH"'
fest not on PATH
```

### 5. `just`

```
$ just --list | grep -i 'vhs\|pipe'
    vhs ...         # Record the booth with VHS
$ just --list vhs
Available recipes:
    booth
    demo        # the worker (cans doctor puts both in ~/.cans/native/bin).
    doctor
    pipe        # Real mouth, no audio: a script pipes lines in, wavs land in out/.
    record tape
                              # `just vhs pipe` is still there

$ just dist check
cd … && goreleaser check
  • checking                                  path=.goreleaser.yaml
  • 1 configuration file(s) validated
  • thanks for using GoReleaser!
exit=0
```

### Finding for the operator — `festivals/CA0001/` fails the phrase lock in a public repo

`veronica-agent/cans` is **public** (`gh repo view --json visibility` → `PUBLIC`). The committed
`festivals/CA0001/` snapshot still carries 24 phrase-lock hits across 11 files, 10 of them the
tracked `festivals/CA0001/001_INGEST/input_specs/` explore pack — the exact directory **D009
excludes** from this festival's snapshot. Shapes: the account display name spelled out; a market
table rejecting a chat-partner product framing; the retired v1 `cans say` example line; and, in
`explore-recommend.md` / `seed.md`, an **earlier form of the professional-grep pattern itself**,
which `CONTEXT.md` marks campaign-private.

`e5a5197` ("strip private campaign guts from the public festival tree") already scrubbed that tree
once and deliberately left these, so some are intentional public facts rather than leaks. Either
way it is a decision about **another festival's shipped output**, larger than this task's
"small fixes only", and outside the `cans-v2` diff. Recorded here and raised to `004_REVIEW`;
nothing in `CA0001` was touched.

## Done when

- [x] Every block above is green and its output is recorded here