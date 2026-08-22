---
fest_type: task
fest_id: 01_stream_loop.md
fest_name: stream_loop
fest_parent: 03_stream
fest_order: 1
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:56.823666-06:00
fest_updated: 2026-08-21T06:58:31.18615-06:00
fest_tracking: true
---


# Task: stream_loop

## Objective

`--stream`: one `Session` for the run, one utterance per stdin line, a flushed record per line, blank lines skipped, bad lines reported and skipped, exit 1 at EOF if any failed.

## Requirements

- [x] `internal/say/stream.go`: `func runStream(ctx context.Context, o Options, stdin io.Reader, stdout, stderr io.Writer) int`; `Run` dispatches on `o.Stream`.
- [x] `doctor.Prepare` → `keep.Load()` → `tts.OpenWith(ctx, lockOpts(...))` once (busy → 75; other → 1) → `defer sess.Close()`.
- [x] `bufio.Scanner` over stdin with `sc.Buffer(make([]byte, 0, 64*1024), 1<<20)`. Two counters: `lineNo` (every stdin line, 1-based) and `idx` (spoken lines, 1-based). `strings.TrimSpace` each line; blank → `continue` without touching `idx` (D005).
- [x] Per spoken line: `out := ""` when `o.Out == ""` else the path for `idx` (task 02 adds the template; until then use `o.Out` only when it has no `%`); `r, err := sess.SayTo(ctx, line, cur, out)`.
- [x] Failure: stderr `line N: <err>`; under `--json` write `{"line":N,"error":"…"}`; `failed++`; continue. Success: record per D006 — JSON `{"line":N,"wav":…,"ttfa_ms":…,"sample_rate":…}`; else `-o` → the path; else `ttfa_ms=N`. **Flush after every record** (wrap stdout in a `bufio.Writer`, `Flush()` each time).
- [x] No `-o`: `play.File` then `RemoveTemp` per line (D007). `-o` + `--play`: play after writing.
- [x] At EOF: `sc.Err()` → stderr, 1. `failed > 0` → 1. Else 0. Backpressure is inherent: the next `SayTo` is sent only after the previous returns — do not add a buffer or goroutine.
- [x] Two record structs (`okRecord`, `errRecord`) so `ttfa_ms: 0` is never dropped by `omitempty`.

## Implementation

1. Keep `runStream` under 50 lines by splitting: `speakLine(...)` (one line → record or error) and `emitStream(...)`.
2. Extend `internal/tts/testdata/fakeworker/main.go` by one branch: `if req.Text == "fail" { print an error record; continue }`. That is the only change to testdata.
3. Tests (error first) in `internal/say/stream_test.go`, fake worker, `CANS_HOME` temp, `CANS_NOPLAY=1`:
   - input `"a\nfail\n\nb\n"`, `Out` = a per-idx path (the template lands in task 02): `001.wav` and `002.wav` exist, no `003.wav`; stdout has two paths; stderr has `line 2: …`; exit 1.
   - same with `JSON`: three records, `{"line":2,"error":…}` in the middle, `line` values 1, 2, 4.
   - **one worker**: point `CANS_WORKER_BIN` at a tiny wrapper script that appends a line to a counter file then `exec`s the built fake worker; after a 5-line stream the counter file has exactly one line.
   - busy: hold the lock, `Wait: 0` → 75 before any stdin is read.
   - empty stdin → 0 with nothing on stdout.

## Done when

- [x] Tests green; `CANS_NOPLAY=1 go test ./...` green; `gofmt -l .` empty; `go vet` clean
- [x] Manual, recorded in the testing gate: `printf 'Put the cans on.\nOne worker, one load.\nFiles land here.\n' | ./bin/cans say --stream -o '/tmp/cans-out/%03d.wav' --json` → three records; `pgrep -f qwen3-tts-worker | wc -l` during the run → 1