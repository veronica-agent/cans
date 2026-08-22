# D003 — flock via stdlib syscall.Flock, polled under ctx

**Decision:** `syscall.Flock(fd, LOCK_EX|LOCK_NB)` on `CANS_HOME/mouth.lock` (created `0644` if missing, never deleted). Acquire loops on `LOCK_NB` with a 100 ms sleep, checking `ctx.Done()` and the `--wait` deadline each turn. Release = `LOCK_UN` + close the fd.

**Why:** No new dependency (`golang.org/x/sys` stays indirect). The kernel drops a `flock` when the holder dies, so Ctrl-C, a panic, or `kill -9` cannot wedge the next run — the reason to prefer it over a PID file with staleness heuristics. Polling keeps cancellation and bounded wait trivial; a blocking `LOCK_EX` in a goroutine cannot be cancelled.

**Not:** A PID file. A socket. FIFO fairness.
