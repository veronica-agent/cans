# GitHub identity

All Veronica work is **Veronica** — campaign and engine. Lance stays on My_Tools and org admin.

| | Campaign (`veronica-campaign`) | Engine (`veronica-voice`) |
|--|--|--|
| Author | `Veronica` / `318153306+veronica-agent@users.noreply.github.com` | same |
| SSH host | `github-veronica-agent` | `github-veronica-agent` |
| Remote | `Obedience-Corp/veronica-campaign` | `lancekrogers/veronica-voice` (she has write; repo not transferred) |
| `gh` while committing | `gh auth switch -u veronica-agent` | `gh auth switch -u veronica-agent` |

`gh auth switch` only changes the API token. The green squares come from the **commit email** plus access to the repo.

She is a write collaborator on the campaign repo and the engine repo. Do **not** add her to the Obedience-Corp org.

On a new machine: `just identity`, then confirm `ssh -T git@github-veronica-agent` says `Hi veronica-agent`.
