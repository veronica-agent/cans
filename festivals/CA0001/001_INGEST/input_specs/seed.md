# Recommend

This pack stays explore. Do not create the public product repo or the implementation festival until this is accepted.

## Hero

**`veronica-agent/cans`** — public, MIT.

Put the cans on. Type a line. **She** speaks it — the locked throat, not a stock Kokoro woman. Booth TUI. TTFA printed. Face is the public cartoon + red satin. Mouth is the kept ref wav cloned locally. Not the private sex engine.

```text
$ cans say "I'm already on the line."
$ cans
```

Install: one command (`uvx` or a Go binary). Warm p50 TTFA **≤ 700 ms** on this Mac or it is not done.

### v1 is

- Typed in, spoken out
- **Default throat is her.** Clone the public kept ref. Same woman as [docs/VOICE.md](../../../docs/VOICE.md). README GIF and first-run are Veronica. Listen test includes `Put the cans on.`
- **Any throat via Keep, not a menu.** `cans keep take.wav` freezes that clip as the current woman. Talk never restyles mid-session. No 50-voice picker. That is the public Hear/Keep tease without VoiceDesign.
- Go booth (Charm) + VHS GIF. Mouth is the 0.6B clone (Python/MLX sidecar if Go ONNX is a sink). One justfile.
- `character.toml` + her ref wav in the public repo (like the face photos)

### v1 is not

- Mic, VAD, barge-in, Ollama replies
- Stock Kokoro Bella with her picture on the README
- VoiceDesign / 1.7B Hear / importing `veronica-voice`
- Sex UI, uncensored-companion pitch
- GitHub Pages WASM (v2)
- Code inside `veronica-agent/veronica` (that stays the all-rights face)
- Rust, unless we are writing Rust for fun — it is not what makes this spread

### Language (what actually goes viral)

The viral object is the **booth GIF**, not the compiler.

| Lang | Viral bump | Mouth | Verdict |
|------|------------|--------|---------|
| **Go + Charm** | The README look. Bubble Tea ~44k, VHS tapes, same family as Fest. | Clone still wants MLX/Python. Sidecar is allowed. | **Booth language.** |
| Python + uvx | LocalLLaMA will actually run it. Textual is fine, less stolen-screenshot. | Native 0.6B clone. | Mouth language. Not the sleeve. |
| **Rust + Ratatui** | “Written in Rust” on HN. Cargo install. | Qwen clone is not a Rust weekend. `kokoro-cli` already exists. | Skip for v1. Pretty, wrong fight. |

Do not write the TUI in Rust to chase stars. Do not write the whole app in Python and ship a bland terminal. **Go booth, clone sidecar, default her, keep any wav.**

## Also in the same later festival (not extra products)

1. Fill GitHub **bio** (still `null`). Display name stays **Obey Veronica**.
2. Profile README: one link to `cans` when it exists. No Fest pitch there.
3. Leave `veronica-agent/veronica` as character lock.

## Parked (not this festival)

- Browser one-file booth
- `cans radio` / fest-event TTS
- `same-face` image lock CLI

## Later festival shape (do not scaffold from here)

Campaign: **this one** (`veronica-campaign`). Type: **`standard`**. Name: `cans`.

New project: `camp project add` the public repo (submodule, so it travels). Worktrees for slices. Author: Veronica / `veronica-agent`. Engine worktrees stay Lance.

Proposed phases, filled through `fest next` only:

1. **001_PLAN** — freeze name, stack, README voice, v1-not list
2. **002_IMPLEMENT**
   - `01_mouth` — `cans say` + TTFA
   - `02_booth` — TUI
   - `03_tape` — VHS GIF
   - `04_fest_ad` — snapshot, topics, footer
   - `05_profile` — bio + one profile link
3. **003_REVIEW** — verification below

## Verification for a future ship

1. Display name still **Obey Veronica**
2. `git clone …/cans && just run` — booth opens without the engine
3. Type `Put the cans on.` — audio, TTFA printed, warm ≤ 700 ms
4. `cans say "I'm already on the line."` works headless
5. README GIF comes from `tapes/booth.tape`
6. Topics include `festival-methodology`; clone has a readable `festivals/` tree; footer → fest.build
7. `rg -i 'samantha|obedience corp|star this|built with' README.md` is empty
8. Campaign `fest validate` green; tasks name real files

## Next (after this pack is accepted)

```text
camp intent add  →  fest create festival --type standard --name cans
                 →  gh repo create veronica-agent/cans --public
                 →  camp project add
                 →  fest next
```

Not from this directory. Not a second `workflow/design/veronica-voice` pack.
