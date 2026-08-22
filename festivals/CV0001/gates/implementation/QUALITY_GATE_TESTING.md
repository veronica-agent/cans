---
# Template metadata (for fest CLI discovery)
id: QUALITY_GATE_TESTING
aliases:
  - testing-verify
  - qg-test
description: Standard quality gate task for testing and verification

# Fest document metadata (becomes document frontmatter)
fest_type: gate
fest_id: <no value>
fest_name: Testing and Verification
fest_parent: <no value>
fest_order: <no value>
fest_gate_type: testing
fest_autonomy: medium
fest_status: pending
fest_tracking: true
fest_created: 2026-08-21T04:32:49-06:00
---

# Gate: Testing and Verification

Verify all functionality implemented in this sequence works correctly.

## Test Categories

### Unit Tests

- [ ] All unit tests pass
- [ ] New/modified code has test coverage
- [ ] Tests are meaningful (not just coverage padding)

### Integration Tests

- [ ] Integration tests pass
- [ ] Components work together correctly

### Error Handling

- [ ] Invalid inputs are rejected gracefully
- [ ] Error messages are clear and actionable
- [ ] Recovery paths work correctly

## Verification

- [ ] Build completes without warnings
- [ ] No regressions introduced
- [ ] Coverage meets project requirements

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

_(paste command output here)_
