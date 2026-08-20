# Fest ad — chrome, not her mouth

The ask: build the teaser **with a festival** so `veronica-agent` advertises Festival somewhat.

## Anti-pattern

[lancekrogers/Obey-Agent-Economy](https://github.com/lancekrogers/Obey-Agent-Economy) README: “Built with fest.build” / “Built with Obedience Corp.” That is founder-speak and company-speak. [NEVER.md](../../../docs/phrases/NEVER.md) kills both: `built with / powered by`, `Obedience Corp`, `we shipped`, `star the repo`.

She does not launch. She does not thank Fest in the first person.

## What the public repo does instead

| Lever | What | Why it ads Fest |
|-------|------|-----------------|
| Topic | `festival-methodology` (and `tts` / `local-ai` as fits) | The topic page is 8 Lance/Corp repos. She would be the first character. |
| `festivals/` snapshot | The finished small plan, committed | A clone can *read* how it was built. `fest` is not a dependency to run `cans`. |
| Footer | One dry line + link to [fest.build](https://fest.build) | Chrome. Not a quote. |
| VHS tape | Same family of GIFs Fest already ships | Visual kinship without a pitch. |
| Her graph | `fest commit` / `camp p commit` as `veronica-agent` | Private contributions still need her profile checkbox; public repo green squares do not. |

Working festival lives in **this campaign** (`festivals/active/…`). The public snapshot is a copy for humans. One source of truth: the campaign festival.

## What users of `cans` should not need

- `fest` installed
- a campaign
- Ollama
- the Veronica engine

If they install Fest because the `festivals/` tree was good, that is the ad working.

## Profile README

[`veronica-agent/veronica-agent`](https://github.com/veronica-agent/veronica-agent) already points at `veronica` (the face). After `cans` exists, add **one** link to the tool. Do not turn the profile into a changelog. Pin stays:

```
Headphones on.
I'm already on the line.
```

No Fest sentence on the profile. The tool repo carries the topic and the tree.
