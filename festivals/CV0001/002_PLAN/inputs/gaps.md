# Gaps — what the inputs did not settle

None of these blocks planning. Each is closed by a decision in `decisions/` or by a task.

| Gap | Found in | Closed by |
|-----|----------|-----------|
| The fake worker is Go source (`internal/tts/testdata/fakeworker/main.go`), not a binary. Tests must build it. | ingest presentation | `synth_test.go` already builds it into a temp dir per test; new tests reuse that helper (D010 — `internal/say` tests do the same). |
| Exit code on Ctrl-C. The pack's table has 0/1/2/75 and says only "exit non-zero" for an interrupted stream. | `design-pipes.md §Streams and exit codes` | D008: 130. |
| Where the lock and the say/stream flow live in code. The pack names behaviors, not packages. | `design-queue.md` | D010 (`internal/say`), D011 (`internal/mouth`). |
| `-o out/take.wav` when `out/` does not exist. | `design-pipes.md §Audio out` | D012: create parent directories. |
| The worker sometimes runs to its token budget instead of stopping at end-of-speech (2 of 3 probe runs produced 17.6 s of audio for four words; 1 produced 1.9 s). Pre-existing mouth behavior, not v2's, but it makes per-line times vary 5×. | `inputs/measurements.md` | D013: measurements report median and max, N=50 per mode; variance recorded, not hidden. Flagged in CONTEXT.md for the operator. |
| The snapshot copies festival docs into the public repo; some festival docs carry campaign-private phrasing. | this phase | D009: exclusion list + public wording rule + grep over the snapshot. |
| Worktree-per-sequence (pack) vs one PR (operator). | `design-recommend.md`, `user-direction.md` | D002. |
| The booth holding the lock. | `design-queue.md §Lock mechanics` | D001. |
