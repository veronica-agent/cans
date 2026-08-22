# D013 — Measurements: median and max, N=50

**Decision:** `inputs/measurements.md` records, with the command for each: worker GGUF load time (ready), one-shot synthesis time, worker max RSS, cold `cans say` wall — each over several runs, as median and max. `03_stream` adds: a 50-line stream vs the same 50 lines as a loop of `cans say -o`, and `xargs -P 8` over 24 lines with `pgrep -fc qwen3-tts-worker` sampled every second (max must be 1). N=50 rather than the pack's 200 because the worker's per-line time varies ~5× (it sometimes runs to its token budget instead of stopping at end-of-speech); the margin, not the magnitude, is the claim, and the per-line load cost being removed is constant in N.

**Why:** The pack's estimates predate the native mouth, and a margin nobody can reproduce is marketing. The variance is pre-existing mouth behavior — recorded and flagged, not fixed here.

**Not:** Quoting any number without its command. Running measurements concurrently with anything else (two workers on Metal at once inflated an early probe 8×).
