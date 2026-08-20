# cans

<p align="center">
  <img src="docs/01-western-cartoon.jpg" width="380" alt="Veronica" />
</p>

<p align="center"><b>Put the cans on.</b></p>

<p align="center">
  <img src="docs/booth.gif" width="640" alt="The cans booth" />
</p>

Type a line. She speaks it on the machine, in a kept voice.

## Install

```bash
git clone git@github.com:veronica-agent/cans.git
cd cans
just build quick
just run install
```

`just run booth` opens the room. `just run say "Put the cans on."` speaks one line.

## Usage

| | |
|---|---|
| `cans` | Booth. Type, enter, she talks. Throat stays put for the session. |
| `cans say "…"` | One line, then exit. Prints `ttfa_ms`. |
| `cans keep take.wav -text "the words in the wav"` | Pin that throat until the next keep. |

Default voice is Veronica (`voices/veronica/ref.wav`). Keep needs both the clip and the transcript.

```bash
cans keep take.wav -text "Just like that, feel the rhythm of my voice."
```

From just, flags after `--`:

```bash
just run keep take.wav -- -text "Just like that, feel the rhythm of my voice."
```

## Make

| | |
|---|---|
| `just build quick` | `bin/cans` |
| `just test unit` | Tests (no MLX) |
| `just vhs booth` | Re-record `docs/booth.gif` |

## Machine

macOS on Apple Silicon. The mouth is [Qwen3-TTS 0.6B Base](https://huggingface.co/mlx-community/Qwen3-TTS-12Hz-0.6B-Base-bf16) through [mlx-audio](https://github.com/Blaizzy/mlx-audio). `uv` installs the sidecar.

---

[Festival](https://fest.build)
