# D001 — Lock lifetime equals Session lifetime; the booth holds it

**Decision:** The mouth lock is acquired before `StartWorker` and released after `Client.Close` returns — its lifetime is exactly one `tts.Session`. The booth opens one `Session` for its whole run, so it holds the lock for its whole run. A script started alongside waits (stderr: `waiting for the mouth…`) or exits 75 with `--nowait`. A booth started while a script holds the mouth waits the same way, before the TUI opens.

**Why:** Releasing per line while keeping the worker resident would let a second worker load — two resident workers is the exact failure the lock exists to prevent. Closing the booth's session per line throws away the warm path that is the booth's point.

**Not:** A lock the booth skips. A lock released between utterances.
