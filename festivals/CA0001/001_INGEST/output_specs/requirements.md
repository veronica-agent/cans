# Requirements

## P0

1. `cans say <text>` synthesizes with the current kept throat and plays audio. Prints TTFA.
2. `cans` (no args) opens a Charm booth: prompt, waveform/status, TTFA. Typed enter speaks.
3. First run uses her public kept ref wav (same woman as docs/VOICE.md).
4. `cans keep <wav>` sets the current throat. Subsequent say/booth use it. No restyle mid-session.
5. Repo private on `veronica-agent`. Git author Veronica. Remote `github-veronica-agent`.
6. Does not import `veronica-voice`. Does not call VoiceDesign.

## P1

1. VHS GIF of the booth in README (`tapes/booth.tape`).
2. `character.toml` + her cartoon/photoreal in the repo (copied from campaign docs, not restyled).
3. Festival snapshot under `festivals/` in the cans repo (readable copy).
4. README in her register. Footer may link https://fest.build
5. `justfile`: install, run, say, keep, test, tape, lint
6. GitHub profile bio filled (`The sexiest voice assistant.`) and profile README links cans.

## P2

1. `cans radio` (speak festival events)
2. Browser one-file booth
3. Public visibility (user locked private for v1)
