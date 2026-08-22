# 001_INGEST — presentation

## What we are building

1. `cans say` gains `-o` (write the wav where told, don't play, don't delete) and stdin (`echo x | cans say`, `cans say -`).
2. `--stream` speaks one utterance per stdin line over **one** warm `Session` — one model load instead of N.
3. `--json` puts JSONL records on stdout, flushed per utterance; prose stays on stderr.
4. A mouth lock (`flock` on `CANS_HOME/mouth.lock`) keeps exactly one worker resident under any loop, `xargs -P`, or second terminal. `--nowait` → 75, `--wait` bounds the block.
5. A pipe tape, a README scripting section, and the CV0001 festival tree in the public repo.

## Locked — do not re-ask

| ID | Decision |
|----|----------|
| D001 | Lock lifetime = `Session` lifetime; the booth holds it for its whole run |
| D002 | One worktree, one branch (`cans-v2`), one PR |
| D003 | stdlib `syscall.Flock`, non-blocking acquire polled under `ctx`; lock never deleted |
| D004 | `say` flags interleave with text both orders, like `parseKeep`; exactly 7 flags, nothing else |
| D005 | Blank stream lines skipped (no index consumed); other failures continue, exit 1 at EOF |
| D006 | `{"wav","ttfa_ms","sample_rate"}`; stream adds `"line":N`; bare `-o` prints the path; plain `say` keeps `ttfa_ms=N` |
| D007 | `--stream` without `-o` plays each line in order |

## P0 (20 items, all cited in `requirements.md`)

`-o` / `--play` · stdin one-shot / `-` / TTY-empty = exit 2 · `--stream` over one `Session` · `%03d` template · per-line failure policy · stream-to-speakers · `--json` records · exit codes 0/1/2/75 · **stdout data, stderr prose** · **`cans say "x"` unchanged from `1e8cea2`** · interleaved flag grammar · mouth lock + wait line + `--nowait`/`--wait` · booth holds the lock · ctrl-C / `kill -9` safety · fake-worker testable.

P1: pipe tape, README scripting section, CV0001 snapshot, measurements with the command that produced them.
P2 (**do not build**): `cansd`, mid-utterance cancel, `ttfa_ms` fix, ref-text change, re-cut booth GIF, FIFO fairness.

## Interpretations I made

- Freeze point for "unchanged" is **`1e8cea2`**, not the pack's `07ddf48` — the pack predates PR #13. Same behavior, newer base.
- Phase names follow the real tree (`002_PLAN`, `004_REVIEW`), not the pack's `001_PLAN` / `003_REVIEW`; worktree is `cans-v2`, not the pack's `cans-v2-out`.
- Ctrl-C safety and fake-worker testability are **P0**, not quality-nice-to-have: they are items 3, 4 and 6 of the ship bar.
- `--play` is folded into the `-o` surface rather than counted as a sixth "v2 is" item.
- The pack left the booth-vs-lock call open for planning; `CONTEXT.md D001` already closed it, so it is not re-opened.

## The ask

Approve `purpose.md` / `requirements.md` / `constraints.md` / `context.md` so `002_PLAN` can size the measurements and scaffold `01_out` → `05_snapshot`. Confirm P2 stays parked and that the two low open questions (default ref text, non-infinite `--wait`) stay with the operator.
