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

`just install` (same as `just run install`) is `go install` into `$(go env GOBIN)`, or `$(go env GOPATH)/bin` if GOBIN is empty. `cans doctor` fetches the mouth. `just uninstall` removes the binary. `just uninstall --all` also wipes `~/.cans` (the native mouth).

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

## Scripting

The script owns the document. cans speaks what it is handed.

<p align="center">
  <img src="docs/pipe.gif" width="680" alt="A script piping lines into cans" />
</p>

`--stream` reads stdin a line at a time and speaks each line through one worker: one model load for the whole document. Blank lines are skipped and do not take an index.

```bash
# one wav per paragraph, script decides what a paragraph is
awk -v RS='' '{gsub(/\n/," "); print}' chapter.md | cans say --stream -o 'out/%03d.wav'

# a line at a time, with the script's own naming
i=0
while IFS= read -r line; do
  i=$((i+1))
  cans say "$line" -o "$(printf 'out/%03d.wav' "$i")"
done < lines.txt

# metadata for a build step
cans say --stream --json -o 'out/%03d.wav' < lines.txt | jq -r 'select(.error==null) | .wav' > manifest.txt
```

| Flag | Effect |
|------|--------|
| `-o`, `--out` | Write the wav here. Under `--stream` the path takes one `%d`, as in `out/%03d.wav`. |
| `--stream` | Read stdin line by line, one utterance per line, one worker for all of them. Text on argv is a usage error. |
| `--json` | One JSON record per utterance on stdout; `--stream` adds `"line"`. `wav` is present only with `-o`. |
| `--play` | Play the wav as well as writing it. Needs `-o`. |
| `--nowait` | Do not queue behind another cans: give up at once. |
| `--wait <dur>` | Queue for at most that long. Without either flag, cans waits. |
| `-` | Read one utterance from stdin. Text on argv is a usage error. |

| Exit | Meaning |
|------|---------|
| `0` | Spoke it. |
| `1` | Runtime failure, or a line in a stream failed. |
| `2` | Usage error. |
| `75` | Another cans holds the mouth and the wait was refused or ran out. |
| `130` | Interrupted. |

stdout carries data — a JSON record, a wav path, or `ttfa_ms=N`, the worker's total synthesis time for that line. Prose goes to stderr.

Ctrl-C stops the stream: the line being spoken is dropped, finished wavs stay, exit 130. A second Ctrl-C stops at once.

Without -o the wav is a temp file removed after playback. `--json` without `-o` omits `wav` so a script is not handed a path that will not survive. Pass `-o` if the file is needed.

---

Built with [Festival](https://fest.build)
