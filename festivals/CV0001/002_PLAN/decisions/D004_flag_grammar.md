# D004 — say flags interleave with text, both orders

**Decision:** `cans say` parses its arguments the way `keep` already does (`cmd/cans/main.go` `parseKeep`): flags and text may interleave, both orders work. The flag set is exactly `-o/--out <path>`, `--json`, `--stream`, `--play`, `--nowait`, `--wait <duration>`, and a bare `-` meaning stdin. Positionals are joined with single spaces into the text. Unknown flag → usage on stderr, exit 2.

**Why:** Stdlib `flag` stops at the first positional; the loop examples in `design-pipes.md` put text first (`cans say "$line" -o out.wav`).

**Not:** A config file. A `--voice`. Any flag not in `design-pipes.md`.
