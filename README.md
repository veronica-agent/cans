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

macOS, Apple Silicon. No Python, no uv.

```bash
brew install --HEAD veronica-agent/tap/cans
```

The mouth is a native [Qwen3-TTS worker](https://github.com/Obedience-Corp/qwen3-tts-native). First `cans doctor` (or first `cans`) downloads it once, about 1.6 GB, into `~/.cans/native`.

```bash
cans doctor
cans
```

<details>
<summary>curl</summary>

```bash
curl -fsSL https://raw.githubusercontent.com/veronica-agent/cans/main/install.sh | sh
```

Until a tagged release exists, that script needs Go and installs `cans` onto PATH. Then `cans doctor`.

</details>

<details>
<summary>From source</summary>

```bash
git clone https://github.com/veronica-agent/cans.git
cd cans
just install
```

`just install` (same as `just run install`) is `go install` into `$(go env GOBIN)`, or `$(go env GOPATH)/bin` if GOBIN is empty. `cans doctor` fetches the mouth. `just uninstall` removes the binary, brew formula, and `~/.cans`.

</details>

## Commands

```bash
cans
cans say "Put the cans on."
cans keep take.wav -text "the words in the wav"
cans doctor
```

Keep needs the clip and the words spoken in it. Throat stays put for the session.

The CLI is Go. The mouth is a native [Qwen3-TTS](https://github.com/Obedience-Corp/qwen3-tts-native) worker that clones a wav. No Python.

---

Built with [Festival](https://fest.build)
