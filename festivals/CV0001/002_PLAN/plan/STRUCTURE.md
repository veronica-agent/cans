# Structure — CV0001 cans-v2

```
cans-v2-CV0001
├── 001_INGEST          (workflow)  design pack + operator direction → output_specs
├── 002_PLAN            (workflow)  this phase: gaps, decisions D001–D013, measurements, plan, scaffold
├── 003_IMPLEMENT       (sequences, on branch cans-v2)
│   ├── 01_out          say args, internal/say, -o, stdin, --json, exit codes
│   ├── 02_lock         internal/mouth flock, OpenWith, --nowait/--wait, booth holds it
│   ├── 03_stream       --stream loop, %03d template, Ctrl-C → 130, measurements
│   ├── 04_tape         tapes/pipe.tape + just vhs pipe, README scripting section
│   └── 05_snapshot     festivals/CV0001 snapshot, professional recheck, fresh-home doctor
└── 004_REVIEW          (freeform)  the ship bar, identity, fest validate, the PR
```

Dependencies: `01 → 02 → 03 → 04 → 05`, strictly. `02_lock` before `03_stream` so the naive path is safe before the fast path is fast. `04_tape` needs `-o` and `--stream`. `05_snapshot` is last so it captures the finished tree. Review last; the PR opens from review.

Each implementation sequence ends with the four gates: testing, review (a different agent), iterate, fest_commit.
