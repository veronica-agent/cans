# D009 — Snapshot exclusions and public wording

**Decision:** The festival snapshot copied into `projects/cans/festivals/CV0001/` excludes `CONTEXT.md`, `001_INGEST/input_specs/`, `.fest/`, `.workflow/`, `.festival-checksums.json`, and the reviewers' hidden working notes `.review-*` (amended in `004_REVIEW`: they are scratch, and `fest validate` warns on their filename shape). Every other festival document is written as if public: it names the campaign phrase lock and the professional-surface grep, never quotes them, and refers to "the operator," not a person. `05_snapshot` runs the grep from `CONTEXT.md §Professional grep` over the snapshot directory and fails if anything matches.

**Why:** The tree is the fest-ad; a stranger reads it. v1's snapshot already exposed more than the current lock allows; v2 does not add to that.

**Not:** Scrubbing after the fact. Copying `input_specs/` (the design pack stays in the campaign).
