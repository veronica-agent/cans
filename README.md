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

macOS on Apple Silicon.

```bash
brew install --HEAD veronica-agent/tap/cans
cans doctor
cans
```

or:

```bash
brew install uv
curl -fsSL https://raw.githubusercontent.com/veronica-agent/cans/main/install.sh | sh
cans doctor
```

`cans doctor` extracts her clip, syncs the sidecar, and checks the machine. The first spoken line downloads [Qwen3-TTS 0.6B Base](https://huggingface.co/mlx-community/Qwen3-TTS-12Hz-0.6B-Base-bf16) through [mlx-audio](https://github.com/Blaizzy/mlx-audio).

## Usage

| Command | |
|---|---|
| `cans` | Booth. Type, enter, she talks. Throat stays put for the session. |
| `cans say "…"` | One line, then exit. Prints `ttfa_ms`. |
| `cans keep take.wav -text "the words in the wav"` | Pin that throat until the next keep. |
| `cans doctor` | Sidecar + machine check. |
| `cans version` | Print version. |

Default voice is Veronica. Keep needs both the clip and the transcript.

```bash
cans keep take.wav -text "Just like that, feel the rhythm of my voice."
```

## From source

```bash
git clone git@github.com:veronica-agent/cans.git
cd cans
just build quick
just run install
./bin/cans doctor
```

| Recipe | |
|---|---|
| `just build quick` | `bin/cans` |
| `just test unit` | Tests (no MLX) |
| `just dist snapshot` | Local goreleaser tarball |
| `just vhs booth` | Re-record `docs/booth.gif` |

---

[Festival](https://fest.build)
