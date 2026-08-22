# D010 — internal/say owns the say and stream flow

**Decision:** A new package `internal/say` exposes `Run(ctx, Options, io.Reader, io.Writer, io.Writer) int` — options in, stdin/stdout/stderr in, exit code out. `cmd/cans` parses flags into `say.Options` and returns `Run`'s code. One-shot, `-o`, stdin, `--json`, `--stream`, the lock flags and cancellation all live there. `internal/tts` gains `SayTo` (caller-supplied output path) and `OpenWith` (lock options); nothing else.

**Why:** `cmd/cans/main.go` is already the flag parser for `keep`; putting a 200-line stream loop in it breaks the file limit and makes the flow untestable without the binary. `say.Run` is tested directly against the fake worker.

**Not:** Growing `internal/tts/worker.go` (196 lines). A second binary.
