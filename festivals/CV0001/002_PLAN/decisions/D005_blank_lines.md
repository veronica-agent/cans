# D005 — Stream: blank lines skipped; other failures continue

**Decision:** In `--stream`, a stdin line that is empty after trimming is skipped and does not consume an output index. Any other per-line failure is reported on stderr (`line N: <error>`), emitted as `{"line":N,"error":"…"}` under `--json`, and the stream continues. At EOF the exit code is 1 if any line failed, else 0.

**Why:** Speaking nothing is meaningless; skipping matches `xargs` and keeps `%03d` indices dense. A 200-line render must not lose 199 lines to one bad one.

**Not:** Aborting on first error. Treating blank lines as errors.
