# Implementation plan — CV0001 cans-v2

Branch `cans-v2` in `projects/worktrees/cans/cans-v2`, based on `1e8cea2`. One PR at the end. Decisions are in `../decisions/`; requirements in `../../001_INGEST/output_specs/requirements.md` (cited as P0-n / P1-n).

## Shape of the code after this festival

```
cmd/cans/main.go          thin: dispatch, parseKeep, parseSay → say.Run
cmd/cans/say_args.go      parseSay: interleaved flags + text (D004)
internal/say/             Run(ctx, Options, stdin, stdout, stderr) int — one-shot, -o, stdin, --json, --stream, cancel (D010)
internal/mouth/           Acquire / Release — flock on CANS_HOME/mouth.lock (D011)
internal/tts/session.go   OpenWith(ctx, Options) — takes the lock before StartWorker; SayTo(ctx, text, cur, out)
internal/tts/synth.go     SayWith unchanged; SayTo added
internal/booth/booth.go   Run uses OpenWith (waits, stderr line) — holds the lock for the session (D001)
tapes/pipe.tape           second tape; just vhs pipe
festivals/CV0001/         snapshot (D009)
```

Nothing else moves. `internal/tts/worker.go` stays at 196 lines.

## 003_IMPLEMENT

### 01_out — `-o`, stdin, `--json`, exit codes (P0-1…4, 9, 10, 15–18, 20)

| # | Task | Files | Done when |
|---|------|-------|-----------|
| 01 | `parse_say` — `parseSay(args) (say.Options, error)`: interleaved flags/text, `-o/--out`, `--json`, `--play`, bare `-`; unknown flag → error (exit 2). `--stream`, `--nowait`, `--wait` are **parsed** here too so the grammar is settled once, but `Run` rejects `--stream` with exit 2 "not yet" until `03_stream` — no: they are accepted and wired in their own sequences; this task only adds the three flags to the parser with tests. | `cmd/cans/say_args.go`, `say_args_test.go` | table-driven tests: both orders, `-` handling, unknown flag, `--wait` duration parse |
| 02 | `internal_say` — create `internal/say` with `Options`, `Run`; move the `case "say"` body into it; `tts.SayTo(ctx, text, cur, out)` + `(*Session).SayTo`; `out == ""` keeps the temp path; `-o` → MkdirAll parent (D012), write there, never `RemoveTemp`; `--play` plays after writing. `CANS_SAY_BIN` path: copy the script's wav to `out`. | `internal/say/say.go`, `internal/tts/session.go`, `internal/tts/synth.go`, `cmd/cans/main.go` | `cans say "x"` output/behavior byte-identical (P0-17); `-o` writes and leaves the file; tests on fake worker + `CANS_SAY_BIN` |
| 03 | `stdin_json_exit` — stdin as one utterance when no text (`-` or empty argv with non-TTY stdin); TTY + empty → exit 2; `--json` record via `json.Encoder` + flush; exit-code mapping 0/1/2; stdout/stderr discipline; `main_test.go` coverage through `run()` with injected stdin. | `internal/say/say.go`, `cmd/cans/main.go`, tests | `echo x \| cans say -o t.wav` works; `cans say < /dev/tty`-style TTY case is exit 2 (inject `isTTY`); JSON record shape matches D006 |

### 02_lock — the mouth lock (P0-11…14, 15 (75), 19 (kill -9), 20)

| # | Task | Files | Done when |
|---|------|-------|-----------|
| 01 | `flock` — `internal/mouth`: `Acquire`, `ErrBusy`, `Release`, per D003/D011. Tests: second acquire `wait=0` → `ErrBusy`; release → reacquire; ctx cancel while waiting → `ctx.Err()`; bounded wait → `ErrBusy` after deadline; `onWait` fires once; cross-process: re-exec test helper holds the lock, `Process.Kill()`, next acquire succeeds. | `internal/mouth/lock.go`, `lock_test.go` | all of the above green; file never deleted |
| 02 | `session_lock` — `tts.OpenWith(ctx, Options{Wait, OnWait})` acquires `ship.Home()/mouth.lock` **before** `StartWorker`, stores it on `Session`, `Close` releases **after** `Client.Close`. `Open` = `OpenWith` defaults (wait forever, stderr `waiting for the mouth…`). `SayWith` unchanged in behavior; `SayTo` gains an options variant. Test: hold Session A (fake worker); `OpenWith(wait=0)` with `CANS_WORKER_BIN` pointing at a **missing** file returns `ErrBusy`, proving the lock precedes the worker start. | `internal/tts/session.go`, `session_test.go` | ordering test green; existing tts tests unchanged |
| 03 | `flags_booth` — `--nowait` (wait 0) and `--wait <dur>` wired from `Options` to `OpenWith`; `mouth.ErrBusy` → exit 75 with `mouth busy` on stderr; booth `Run` uses `OpenWith` (wait forever, stderr line before the TUI). | `internal/say/say.go`, `internal/booth/booth.go`, tests | `cans say --nowait x` while another cans holds the mouth → 75; booth waits |

### 03_stream — `--stream` (P0-5…9, 19, 20; P1-4 numbers)

| # | Task | Files | Done when |
|---|------|-------|-----------|
| 01 | `stream_loop` — in `say.Run`: `OpenWith` once, `bufio.Scanner` (1 MiB buffer) over stdin, skip blank lines (D005), `idx++`, `SayTo` per line, record per line (D006) + flush, play when no `-o` (D007), errors continue, exit 1 at EOF if any failed. | `internal/say/stream.go`, tests | 5-line stream on the fake worker → 5 wavs, 5 records, one worker start |
| 02 | `out_template` — `-o` with `%d`-family verb required in stream mode (else exit 2); exactly one verb; `fmt.Sprintf(tmpl, idx)`; MkdirAll. | `internal/say/template.go`, tests | `out/%03d.wav` → `out/001.wav`…; `%s` rejected; no verb rejected |
| 03 | `cancel` — `signal.NotifyContext` in `cmd/cans`; loop checks `ctx` between lines; on cancel: stop reading, close session, release lock, stderr `interrupted after line N`, exit 130 (D008). Test: cancel ctx after the first record on a never-EOF pipe → first wav present, `Run` returns 130, a fresh `mouth.Acquire(wait=0)` succeeds. README line per D014: the line being spoken is dropped; a second Ctrl-C stops at once. | `internal/say/stream.go`, `cmd/cans/main.go`, tests | test green; no orphaned fake worker (`cmd.Wait` returned) |
| 04 | `measure` — real mouth, nothing else running: (a) 50 lines as `--stream -o 'out/%03d.wav'` vs the same 50 as a `while read` loop of `cans say -o`; wall, median/max per line; `pgrep -fc qwen3-tts-worker` sampled each second (max 1 in both); (b) `xargs -P 8` over 24 lines of `cans say -o` with `pgrep` sampling (max 1), no swap (`vm_stat` pageouts delta 0). Record in `002_PLAN/inputs/measurements.md §Stream`, with every command. | measurements.md | numbers + commands recorded; `pgrep` max = 1 in every run |

### 04_tape — the pipe demo and the README (P1-1, P1-2, P1-5)

| # | Task | Files | Done when |
|---|------|-------|-----------|
| 01 | `pipe_tape` — `tapes/pipe.tape`: three boring technical lines into `lines.txt`, `cat lines.txt \| cans say --stream -o 'out/%03d.wav' --json`, `ls out/`. `just vhs pipe` → `docs/pipe.gif`. Real mouth; `Wait+Screen` if the installed vhs supports it, else generous `Sleep`. | `tapes/pipe.tape`, `.justfiles/vhs.just`, `docs/pipe.gif` | gif regenerates from `just vhs pipe`; frame checked |
| 02 | `readme_scripting` — README "Scripting" section: the three loops from `design-pipes.md`, the flag table, the exit-code table (0/1/2/75/130), the cancel-between-requests line, `docs/pipe.gif` next to the booth gif. `usage` const in `main.go` updated. Professional-surface grep clean. | `README.md`, `cmd/cans/main.go` | grep (CONTEXT §Professional grep) empty; one footer |

### 05_snapshot — the public tree (P1-3, P1-5)

| # | Task | Files | Done when |
|---|------|-------|-----------|
| 01 | `snapshot` — copy the festival into `festivals/CV0001/` with **exactly the exclusion list D009 carries** (see `002_PLAN/decisions/D009_public_snapshot.md`; amended in `004_REVIEW`, so read it there rather than restating it here); add the one README line pointing at the tree next to the CA0001 line, if such a line exists. | `festivals/CV0001/`, `README.md` | tree readable; grep over `festivals/CV0001/` empty |
| 02 | `recheck` — full professional grep over README, docs/, tapes/, festivals/; exactly one footer; `just test unit`, `go vet`, `gofmt -l`; fresh `CANS_HOME` doctor + `cans say -o` with the binary copied outside the checkout; `cans` runs with `fest` absent from PATH. | — | everything green, recorded in the task file |

## 004_REVIEW

`PHASE_GOAL.md` carries the nine-item bar from `design-recommend.md §Ship verification` plus identity (`git log --format=%an`), `fest validate`, and the PR. `BAR.md` holds the exact commands and their recorded output. Sign-off opens the PR from `cans-v2` to `main` under `veronica-agent`, body built from `BAR.md`.

## Verification at every gate

`gofmt -l .` empty · `go vet ./...` · `CANS_NOPLAY=1 go test ./...` (fake worker only) · `cans say "x"` unchanged · no new `go.mod` requires · files < 500 lines, functions < 50 · reviewer ≠ implementer.
