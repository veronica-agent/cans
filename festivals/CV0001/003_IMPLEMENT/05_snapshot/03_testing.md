---
fest_type: gate
fest_id: 03_testing.md
fest_name: Testing and Verification
fest_parent: 05_snapshot
fest_order: 3
fest_status: completed
fest_autonomy: medium
fest_gate_id: testing
fest_gate_type: testing
fest_managed: true
fest_created: 2026-08-21T05:04:57.560643-06:00
fest_updated: 2026-08-21T18:00:17.972843-06:00
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

`05_snapshot` adds **no Go code** — it is the `festivals/CV0001/` snapshot plus the surface
recheck. So the code categories above are the full-suite regression bar: the whole test suite on
the fake worker, plus the one-shot on the real mouth to prove `cans say` still behaves as at
`1e8cea2`. Nothing under `internal/` or `cmd/` changed in this sequence.

Box quiet for the real-mouth step: load `6.76 8.13 8.17` at 17:58 (all < 16) and
`pgrep -fl 'cans/native/bin/qwen3-tts-worker'` empty before it. Run alone.

```
$ gofmt -l .
(no output — exit 0)

$ go vet ./...
(no output — exit 0)

$ CANS_NOPLAY=1 go test -count=1 ./...
ok  	github.com/veronica-agent/cans/cmd/cans	1.940s
ok  	github.com/veronica-agent/cans/internal/audio	0.344s
ok  	github.com/veronica-agent/cans/internal/booth	0.679s
ok  	github.com/veronica-agent/cans/internal/doctor	0.519s
ok  	github.com/veronica-agent/cans/internal/keep	0.989s
ok  	github.com/veronica-agent/cans/internal/mouth	1.398s
ok  	github.com/veronica-agent/cans/internal/play	1.139s
ok  	github.com/veronica-agent/cans/internal/say	7.092s
ok  	github.com/veronica-agent/cans/internal/ship	1.507s
ok  	github.com/veronica-agent/cans/internal/tts	3.570s
exit=0
# -count=1, no cache: all ten packages green on the fake worker, no real mouth involved.
# Error handling, invalid input and recovery are covered there — say/say_test.go 418 lines and
# say/stream_test.go 352 lines are error-cases-first tables from 01_out and 03_stream.

$ git diff origin/main -- go.mod go.sum
(empty — no new dependencies)

$ wc -l $(git diff --name-only origin/main -- '*.go') | sort -n | tail -5
     200 cmd/cans/main_test.go
     201 internal/tts/worker.go
     352 internal/say/stream_test.go
     418 internal/say/say_test.go
    3147 total
# every changed file < 500. internal/tts/worker.go is 201 against the 196 the rules pin
# ("add files, do not grow it") — carried over from 03_stream, not touched here.

$ just build quick
go build -trimpath -ldflags "-s -w -X …/internal/ship.Version=v0.1.0-27-g5e8123d" -o bin/cans ./cmd/cans
# no warnings, clean version string

$ ./bin/cans say "Put the cans on." ; echo "exit=$?"
ttfa_ms=5669
exit=0
# wall 14s; played through the speakers; one-shot output shape unchanged (bare ttfa_ms=N,
# no path, no JSON). Against the 5 839 ms / 13.1 s baseline and the 5 652 ms recorded in
# 03_stream/05_testing this is the same run — v1 behaviour is intact.

$ pgrep -fl 'cans/native/bin/qwen3-tts-worker'     # after the run
(no output — exit 1)                                # worker gone, lock released by the kernel

$ ls -ld /var/folders/…/T/cans-say.*
-rw-------@ 1 user  staff  355 Aug 19 22:09 …/T/cans-say.F26P2oxvoX
# The run's own temp wav was removed: the only cans-say.* file in TMPDIR is a 355-byte
# leftover dated Aug 19 22:09, present before this run and unchanged by it — debris from
# before this branch existed, not a regression. No new temp file appeared.
```

### Sequence checks from the task files

Both re-verified after the last edit to the task documents:

```
$ rg -i '<pattern 1 — CONTEXT.md §Professional grep>' README.md | rg -v '<allowed footer line>'
(no output)
$ rg -i '<pattern 2 — CONTEXT.md §Professional grep>' README.md docs/ tapes/ festivals/CV0001/
(no output)
$ rg -c 'fest.build' README.md
1
$ fest validate festivals/CV0001
VALIDATION PASSED WITH WARNINGS — 90/100
# two warnings, both the cold reviewer's hidden .review-*.md scratch files carried by the rsync
$ git check-ignore -v festivals/CV0001/fest.yaml
(no output — not ignored)
$ find festivals/CV0001 -type f | wc -l
      88

# NOTE, added in 05_iterate: the review gate made those two warnings a Critical (D009 was
# amended to exclude `.review-*`). After the fix the same two commands read 100/100 and 86
# files. Everything else in this gate is unchanged — no Go code moved.
$ git status --short
?? festivals/CV0001/
```

`02_recheck.md` carries the fresh-home doctor / `say -o` / `ffprobe` / `ls $CANS_HOME` output,
`just dist check`, `just vhs pipe`, and the **`festivals/CA0001/` phrase-lock finding** raised to
`004_REVIEW` (pre-existing committed content in another festival's snapshot, untouched here).

**No regressions. Gate green.**