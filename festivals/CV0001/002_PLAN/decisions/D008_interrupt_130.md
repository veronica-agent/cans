# D008 — Interrupted stream exits 130

**Decision:** `cans say` installs `signal.NotifyContext` for SIGINT/SIGTERM. In `--stream`, cancellation stops reading stdin, finishes the in-flight line (the worker has no mid-synth abort), closes the `Session` (`shutdown` + wait), releases the lock, prints `interrupted after line N` on stderr, and exits **130** (128 + SIGINT). Completed wavs stay on disk. One-shot interrupted while waiting for the lock exits 130 too.

**Why:** The pack said only "non-zero." 130 is what every shell user expects from Ctrl-C, and it is distinct from 1 (a line failed) so a rerun script can tell the two apart.

**Not:** Killing the worker mid-utterance. Deleting partial output.
