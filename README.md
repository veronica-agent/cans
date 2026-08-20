<h1 align="center">cans</h1>

<p align="center">
  <img src="docs/01-western-cartoon.jpg" width="220" alt="Veronica" />
</p>

<p align="center"><b>Put the cans on.</b></p>

<p align="center">Type a line. She speaks it.</p>

<p align="center">
  <a href="docs/booth.mp4">
    <img src="docs/booth.gif" width="680" alt="The cans booth" />
  </a>
</p>

<p align="center"><a href="docs/booth.mp4">watch with sound</a></p>

## Install

macOS, Apple Silicon. Needs [uv](https://docs.astral.sh/uv/) before she speaks (`brew` pulls it in).

```bash
brew install --HEAD veronica-agent/tap/cans
cans
```

First run sets up the mouth. `cans doctor` if you want the checklist.

<details>
<summary>curl</summary>

```bash
brew install uv
curl -fsSL https://raw.githubusercontent.com/veronica-agent/cans/main/install.sh | sh
cans
```

Until a tagged release exists, that script needs Go and installs `cans` onto PATH.

</details>

<details>
<summary>From source</summary>

```bash
git clone https://github.com/veronica-agent/cans.git
cd cans
brew install uv
just install
cans
```

`just install` (same as `just run install`) is `go install` into `$(go env GOBIN)`, or `$(go env GOPATH)/bin` if GOBIN is empty.

</details>

## Commands

```bash
cans
cans say "Put the cans on."
cans keep take.wav -text "the words in the wav"
cans doctor
```

Keep needs the clip and the words spoken in it. Throat stays put for the session.

The CLI is Go. The mouth is [Qwen3-TTS 0.6B](https://huggingface.co/mlx-community/Qwen3-TTS-12Hz-0.6B-Base-bf16) on [MLX](https://github.com/Blaizzy/mlx-audio), which is Python, so `cans` starts that process for you.

---

Built with [Festival](https://fest.build)
