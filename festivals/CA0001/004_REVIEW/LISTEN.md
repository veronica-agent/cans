# Review notes — 2026-08-19

Ran from the cans worktree.

| Check | Result |
|-------|--------|
| `git log -1` author | Veronica `<318153306+veronica-agent@users.noreply.github.com>` |
| `gh repo view veronica-agent/cans` | `veronica-agent/cans` |
| `fest validate` | 100 |
| `CANS_NOPLAY=1 go test ./...` | pass |
| `cans say "Put the cans on."` | **spoke.** `ttfa_ms=815` |
| `cans say "There you are."` | **spoke.** `ttfa_ms=2128` |

Sidecar reloads the 0.6B model every process, so warm p50 ≤ 700 ms is **not** met. First successful parse after fixing stdout JSON. Listen test is pass for “it is her throat running”; fail for lock-floor TTFA.

Deferred: keep a warm sidecar daemon. Not a festival blocker for v1 speak.
