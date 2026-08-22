---
fest_type: gate
fest_id: 04_testing.md
fest_name: Testing and Verification
fest_parent: 01_out
fest_order: 4
fest_status: completed
fest_autonomy: medium
fest_gate_id: testing
fest_gate_type: testing
fest_managed: true
fest_created: 2026-08-21T05:04:57.42617-06:00
fest_updated: 2026-08-21T06:09:57.116443-06:00
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

Worktree: `projects/worktrees/cans/cans-v2`. Recorded 2026-08-21.

```
$ gofmt -l .
(empty)

$ go vet ./...
(empty, exit 0)

$ CANS_NOPLAY=1 go test ./...
ok  github.com/veronica-agent/cans/cmd/cans
ok  github.com/veronica-agent/cans/internal/audio
ok  github.com/veronica-agent/cans/internal/booth
ok  github.com/veronica-agent/cans/internal/doctor
ok  github.com/veronica-agent/cans/internal/keep
ok  github.com/veronica-agent/cans/internal/play
ok  github.com/veronica-agent/cans/internal/say
ok  github.com/veronica-agent/cans/internal/ship
ok  github.com/veronica-agent/cans/internal/tts

$ git diff origin/main -- go.mod go.sum
(empty)

$ wc -l $(git diff --name-only origin/main -- '*.go') | sort -n
      76 internal/tts/synth.go
      83 internal/tts/session.go
      98 internal/tts/synth_bin.go
     130 cmd/cans/main.go
     142 cmd/cans/main_test.go
     529 total
```

`git diff --name-only origin/main` does not list untracked files. New files in this sequence, all < 500:

```
      11 internal/say/exit.go
      24 cmd/cans/tty.go
      29 internal/say/options.go
      39 internal/say/input.go
      40 cmd/cans/tty_test.go
      79 internal/say/say.go
      90 internal/say/input_test.go
      98 cmd/cans/say_args.go
     118 cmd/cans/say_args_test.go
     242 internal/say/say_test.go
```

Largest functions (non-test): `parseSay` 48 lines, `run` 44 lines. None over 50.

```
$ ./bin/cans say "Put the cans on." ; echo "exit=$?"
ttfa_ms=36796
exit=0
```

No `cans-*.wav` in `$TMPDIR` before or after. `ttfa_ms` is total synth wall time (CONTEXT deferred item); behavior matches v1: print, play, delete temp.

Sequence 03 Done-when:

```
$ ./bin/cans say < /dev/null ; echo exit=$?
say: empty text
exit=2

$ ./bin/cans say   # under a pty
say: missing text
exit=2

$ echo "Put the cans on." | CANS_NOPLAY=1 CANS_SAY_BIN=<fake> ./bin/cans say -o /tmp/cans-out/a.wav --json
{"wav":"/tmp/cans-out/a.wav","ttfa_ms":12,"sample_rate":24000}
```

`/dev/null` is a char device; `stdinIsTTY` uses `TIOCGETA`, not `Stat` `ModeCharDevice`, so the empty-text case is this one.

`just test unit` green. `just build quick` green.