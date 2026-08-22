---
fest_type: task
fest_id: 02_out_template.md
fest_name: out_template
fest_parent: 03_stream
fest_order: 2
fest_status: completed
fest_autonomy: medium
fest_created: 2026-08-21T05:04:56.824893-06:00
fest_updated: 2026-08-21T07:00:07.01412-06:00
fest_tracking: true
---


# Task: out_template

## Objective

`-o 'out/%03d.wav'` is a template in stream mode, validated up front; anything else is exit 2.

## Requirements

- [x] `internal/say/template.go`: `func checkOut(out string, stream bool) error` and `func outPath(out string, idx int) string`.
- [x] Stream mode with `-o`: the template must contain **exactly one** verb and it must be integer-formatting — `%d`, `%3d`, `%03d`, `%-3d`. `%%` is a literal percent and does not count. Any other verb (`%s`, `%v`, `%f`), two verbs, or no verb → `say: -o needs one %d in --stream` (exit 2).
- [x] One-shot: `-o` is a literal path. A verb in it → `say: -o template needs --stream` (exit 2). `%%` is written as `%`.
- [x] `outPath` = `fmt.Sprintf(out, idx)` (one-shot: `strings.ReplaceAll(out, "%%", "%")`); parent dir created by `SayTo` (D012). Validation runs in `Run` **before** `doctor.Prepare` and before the lock is taken, so a typo never waits on the mouth.

## Implementation

1. Parse with a small scanner over the string rather than a regexp: walk runes; on `%` look at the next rune (`%` → literal; digits or `-` → consume, then require `d`; anything else → error); count verbs.
2. Wire `checkOut` into `Run`; replace the per-idx path in `runStream` with `outPath`.
3. Tests (error first) in `template_test.go`, table-driven: `%s`, `%03d-%d`, `out.wav` in stream (no verb), `%03d` in one-shot, `%q`; then `%03d` → `001`, `%d` → `1`, `out/%02d.wav` → `out/01.wav`, `%%` literal in one-shot, `x%%y%03d` in stream.

## Done when

- [x] `./bin/cans say --stream -o out.wav < lines.txt` → exit 2 with `say: -o needs one %d in --stream`, and the mouth was **not** started (no `waiting` line, instant)
- [x] `CANS_NOPLAY=1 go test ./...` green; `gofmt -l .` empty; `go vet` clean