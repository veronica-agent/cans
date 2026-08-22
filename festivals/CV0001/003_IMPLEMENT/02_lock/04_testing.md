---
fest_type: gate
fest_id: 04_testing.md
fest_name: Testing and Verification
fest_parent: 02_lock
fest_order: 4
fest_status: completed
fest_autonomy: medium
fest_gate_id: testing
fest_gate_type: testing
fest_managed: true
fest_created: 2026-08-21T05:04:57.435166-06:00
fest_updated: 2026-08-21T06:40:57.097595-06:00
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
ok  github.com/veronica-agent/cans/internal/mouth
ok  github.com/veronica-agent/cans/internal/play
ok  github.com/veronica-agent/cans/internal/say
ok  github.com/veronica-agent/cans/internal/ship
ok  github.com/veronica-agent/cans/internal/tts

$ git diff origin/main -- go.mod go.sum
(empty)

$ wc -l $(git diff --name-only origin/main -- '*.go') | sort -n | tail -5
     131 internal/tts/session.go
     132 cmd/cans/main.go
     164 internal/booth/booth.go
     200 cmd/cans/main_test.go
     339 internal/say/say_test.go
```

New untracked: `internal/mouth/lock.go` 115, `lock_test.go` 151, `internal/tts/session_test.go` 76. `worker.go` still 196.

```
$ ./bin/cans say "Put the cans on." ; echo exit=$?
ttfa_ms=5904
exit=0
$ ls -l ~/.cans/mouth.lock
-rw-r--r--  0 Aug 21 06:32 .../mouth.lock
```

Lock file remains after the process exits (never deleted).

Sequence checks (real binary; lock held with `fcntl.flock` on `~/.cans/mouth.lock`, same as a live booth):

```
$ ./bin/cans say --nowait "Put the cans on."
say: mouth busy
exit=75

$ ./bin/cans say --wait 200ms "Put the cans on."
waiting for the mouth…
say: mouth busy
exit=75   elapsed 0.207s
```

Booth under a pty, then `--nowait`: stderr `say: mouth busy`, exit 75, stdout empty. `--nowait` does not print the waiting line.