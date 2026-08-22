# Festival Rules: cans-v2

## Cans-specific (do not violate)

- Git author: `Veronica` / `318153306+veronica-agent@users.noreply.github.com`. `gh auth switch -u veronica-agent` before any push or PR.
- Remote: `git@github-veronica-agent:veronica-agent/cans.git`. Display name stays the one the campaign identity lock already sets — do not change it.
- Public surface is **professional** — an engineer's repo, not a product page: no pitch, no suggestive example text. Every README / tape / fixture / festival line passes the campaign phrase lock (`docs/phrases/NEVER.md`, campaign-private) and the professional-surface grep (pattern in `CONTEXT.md §Professional grep`, campaign-private). Example text is boring and technical — `"Put the cans on."` is fine.
- `cans say "x"` behaves exactly as at `1e8cea2`. One-shot stays one-shot.
- The script owns the document. No `read`, no `-f`, no stripping, no chunking. If a task adds a flag that is not in `design-pipes.md`, the task is wrong.
- No daemon, no queue file, no persisted lock state. The lock is `flock`; the kernel cleans up.
- Exactly one `qwen3-tts-worker` resident, ever. Lock lifetime equals `Session` lifetime.
- stdout is data, stderr is prose. They never mix.
- No new module dependencies. The lock uses stdlib `syscall.Flock`.
- No Python in the shipped payload. `tapes/render-demo-tape.py` is dev tooling and stays.
- Do not import `projects/veronica-voice` or `qwen3-tts-native` into `go.mod`.

## Code

- Files under 500 lines, functions under 50. `internal/tts/worker.go` is at 196 — add files, do not grow it.
- `context.Context` first on anything that does I/O; check `ctx.Err()` before long work; honor cancellation through to the worker.
- Wrap errors with the failing operation (`fmt.Errorf("say: %w", err)` is the project's established style).
- Tests: error cases first, table-driven where there are several shapes, run on the fake worker (`CANS_WORKER_BIN=internal/tts/testdata/fakeworker`) so `go test ./...` needs no real mouth. `CANS_NOPLAY=1` in tests.
- `gofmt -l .` empty, `go vet ./...` clean, `just test unit` green before every gate.

## Process

- `fest next` → do the task → `fest task completed --yes` → at the commit gate `fest commit -m "<type>: <summary>"`.
- All implementation on branch `cans-v2` in `projects/worktrees/cans/cans-v2` (linked to `WI-a2e393`). Never edit `projects/cans` directly.
- Never raw `git commit`. Never "Co-authored-by" or AI attribution in a message.
- Sequence order is load-bearing: `01_out` → `02_lock` → `03_stream` → `04_tape` → `05_snapshot`. Do not start `03_stream` before `02_lock` is committed.
- Implementation and the review gate are done by **different** agents. The reviewer reads the diff cold.
- Record every number you measure (load time, RSS, stream vs loop) in the task file or `002_PLAN/inputs/measurements.md`, with the command that produced it.
- Update `CONTEXT.md` when a decision is made or a blocker is hit.
