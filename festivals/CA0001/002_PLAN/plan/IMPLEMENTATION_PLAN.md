# Implementation plan — cans

Repo: `veronica-agent/cans` (and the worktree once added).

## Layout to create

```
cmd/cans/main.go          CLI: say | keep | (default) booth
internal/tts/             start sidecar, request synth, parse TTFA
internal/keep/            current.lock → wav + ref_text
internal/booth/           bubbletea UI
sidecar/say.py            mlx_audio Base 0.6B clone from ref wav
voices/veronica/ref.wav   kept ref (same woman)
voices/veronica/meta.json ref_text from the clip
character.toml
justfile
```

## Mouth

Clone with `mlx-community/Qwen3-TTS-12Hz-0.6B-Base-bf16`, `generate(text=, ref_audio=, ref_text=)`. Sidecar is local to this repo.

Sidecar prints one JSON line `{path, ttfa_ms, sample_rate}` then exits. Go plays the wav.

Tests mock the sidecar.

## Keep

`cans keep path.wav` writes `state/current.json` `{wav, ref_text}`. Default current is `voices/veronica`. Optional `--text` for the transcript of the keep clip; default empty string if unknown (clone still gets the wav).

## Booth

Full-screen Charm. Input. Status: idle / speaking / ttfa. No mic.

## Tape

`tapes/booth.tape` via VHS. If VHS cannot record audio, the GIF is visual only; `just say` proves the mouth.

## Identity

Every project commit: Veronica email. `fest commit` from the linked worktree.
