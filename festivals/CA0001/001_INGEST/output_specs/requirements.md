# Requirements

## P0

1. `cans say <text>` synthesizes with the current kept throat and plays audio. Prints TTFA.
2. `cans` (no args) opens a Charm booth: prompt, waveform/status, TTFA. Typed enter speaks.
3. First run uses her kept ref wav.
4. `cans keep <wav>` sets the current throat. Subsequent say/booth use it. No restyle mid-session.
5. Git author Veronica. Remote `veronica-agent/cans`.
6. Self-contained. No other voice-engine import.

## P1

1. VHS GIF of the booth in README (`tapes/booth.tape`).
2. `character.toml` + her cartoon in the repo.
3. Festival snapshot under `festivals/` in the cans repo (readable copy).
4. README in her register. Footer may link https://fest.build
5. `justfile`: install, run, say, keep, test, tape, lint
6. Leave the professional GitHub profile README alone.

## P2

1. `cans radio` (speak festival events)
2. Browser one-file booth
3. Extra surfaces (browser booth, radio) stay parked
