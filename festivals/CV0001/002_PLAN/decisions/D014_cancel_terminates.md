# D014 — Ctrl-C terminates an in-flight utterance (amends D008)

**Decision:** On the first SIGINT/SIGTERM, `say` stops reading stdin and closes the session. If the worker is idle, `Close` sends `shutdown` and waits. If the worker is mid-synthesis it cannot abort, so the process is sent SIGTERM and, if still alive 2 s later, SIGKILL (`exec.Cmd.Cancel` + `WaitDelay`). The in-flight line is dropped; finished wavs stay; the lock is released; exit 130; stderr `interrupted after line N` where N is the last stdin line fully processed (spoken or reported), or `interrupted before the first line`. After the first signal the handler is removed (`stop()`), so a second Ctrl-C falls through to the default disposition and ends the process at once.

**Why:** D008 said "Not: killing the worker mid-utterance" and the README was about to promise that Ctrl-C lets the line finish. The review of `03_stream` found `exec.CommandContext` already kills the worker on cancel, and the baseline measurements show a single line can run 17–30 s when the mouth misses end-of-speech. Making a user wait that long after Ctrl-C is worse than a clean terminate: nothing is lost (the in-flight wav is never written on cancel anyway), the kernel frees the worker's memory and the `flock` on exit, and the next `cans say` starts immediately.

**Not:** Waiting for the in-flight line. Pretending cancel is graceful when it is not. Changing the protocol (the worker still has no abort).

**Supersedes:** D008's "Not: killing the worker mid-utterance" and the README line it implied. The README now says: *Ctrl-C stops the stream: the line being spoken is dropped, finished wavs stay, exit 130.*
