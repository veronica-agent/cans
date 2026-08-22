# BAR — commands and recorded output

Fill each block with the exact command and its output when `004_REVIEW` runs. Real mouth, idle machine (`uptime` first). Do not paraphrase output.

## 0. Preconditions

```bash
uptime                       # 1-min load must be below 16
pgrep -fl qwen3-tts-worker   # must print nothing
just build quick
```

## 1. Stream — one worker, one load

## 2. xargs -P 8 — one worker, no pageouts

## 3. Ctrl-C mid-stream → 130, wavs kept, no orphan

## 4. kill -9 → next run unblocked

## 5. Booth holds the lock

## 6. One-shot unchanged; tests green on the fake worker

## 7. Margin (from measurements.md)

## 8. Professional-surface grep + one footer

## 9. Fresh CANS_HOME doctor

## 10. Identity

## 11. fest validate (festival + snapshot)

## 12. Snapshot re-sync

## 13. PR
