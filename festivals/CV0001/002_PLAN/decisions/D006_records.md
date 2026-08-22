# D006 — stdout records

**Decision:**
- no `-o`, no `--json`: `ttfa_ms=N` — unchanged from `1e8cea2`
- `-o path`, no `--json`: the written path, one per line
- `--json`: `{"wav":"…","ttfa_ms":N,"sample_rate":24000}` per utterance, using the existing `tts.Result` tags; `--stream` adds `"line":N` (1-based stdin line number) as the first field; failures are `{"line":N,"error":"…"}`
- every record is followed by a flush

**Why:** Scripts need the path or the record; humans keep the v1 line. stdout is data, stderr is prose.

**Not:** Progress on stdout. Mixed formats.
