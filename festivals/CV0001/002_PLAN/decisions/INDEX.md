# Decisions

| ID | Title |
|----|-------|
| [D001](D001_lock_lifetime_booth.md) | Lock lifetime equals Session lifetime; the booth holds it |
| [D002](D002_one_pr.md) | One worktree, one branch, one PR |
| [D003](D003_flock_stdlib.md) | flock via stdlib syscall.Flock, polled under ctx |
| [D004](D004_flag_grammar.md) | say flags interleave with text, both orders |
| [D005](D005_blank_lines.md) | Stream: blank lines skipped; other failures continue |
| [D006](D006_records.md) | stdout records |
| [D007](D007_stream_plays.md) | --stream without -o plays each line |
| [D008](D008_interrupt_130.md) | Interrupted stream exits 130 |
| [D009](D009_public_snapshot.md) | Snapshot exclusions and public wording |
| [D010](D010_internal_say.md) | internal/say owns the say and stream flow |
| [D011](D011_internal_mouth.md) | internal/mouth is the lock |
| [D012](D012_mkdir_out.md) | -o creates parent directories |
| [D013](D013_measurements.md) | Measurements: median and max, N=50 |
| [D014](D014_cancel_terminates.md) | Ctrl-C terminates an in-flight utterance (amends D008) |
