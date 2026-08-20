# Structure — CA0001 cans

```
cans-CA0001
├── 001_INGEST          (workflow)  specs from explore + this session
├── 002_PLAN            (workflow)  this phase
├── 003_IMPLEMENT       (sequences)
│   ├── 01_scaffold     go module, justfile, identity, voices/veronica
│   ├── 02_mouth        cans say + clone sidecar + TTFA
│   ├── 03_keep         cans keep take.wav
│   ├── 04_booth        Charm TUI
│   ├── 05_tape         VHS GIF
│   └── 06_snapshot     festival copy + README chrome
└── 004_REVIEW          listen test, identity, validate
```

Dependencies: 01 → 02 → 03 → 04 → 05. 06 after 04 (README needs a booth to describe). Review last.
