# D012 — -o creates parent directories

**Decision:** `SayTo` runs `os.MkdirAll(filepath.Dir(out), 0o755)` before writing. `-o out/%03d.wav` works on a fresh checkout without a `mkdir out`.

**Why:** Every loop example in `design-pipes.md` writes into `out/`. Failing on a missing directory is a papercut, not a safeguard.

**Not:** Overwrite protection. The script owns its files.
