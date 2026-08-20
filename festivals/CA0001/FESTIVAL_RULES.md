# Festival Rules: cans

## Cans-specific (do not violate)

- Git author: `Veronica` / `318153306+veronica-agent@users.noreply.github.com`
- `gh` while pushing: `veronica-agent`. Remote: `git@github-veronica-agent:veronica-agent/cans.git`
- Display name stays **Obey Veronica**
- Default throat is her kept ref. Never a stock speaker pretending to be her.
- `keep` changes the woman. Talk does not restyle mid-session.
- Do not import `projects/veronica-voice`
- Do not ship VoiceDesign or 1.7B Hear
- README pull-quote is her ([docs/phrases/NEVER.md](../../../docs/phrases/NEVER.md)): no star-beg, no “built with”, no Obedience Corp, no Samantha, no Lance
- Fest advertising is chrome: topic if we ever go public, `festivals/` snapshot, footer link
- Repo is **private**
- Go TUI. Python sidecar only for the clone model.

## Code

- Keep it small. Files under 500 lines, functions under 50.
- Tests: error cases first. Context if we grow I/O.
- `justfile`: install, run, say, keep, test, tape, lint. No extra help menu.

## Process

- `fest next` → do the task → `fest task completed` → `fest commit`
- Implementation in a worktree: `camp project worktree add <slice> --start-point origin/main --workitem WI-9608d7`
- Relink with `fest link` if the worktree path is the execution dir
- Do not bulk-scaffold empty sequences
