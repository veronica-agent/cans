#!/usr/bin/env python3
"""Fill tapes/demo.tape from demo.tape.in using the baked line."""

from __future__ import annotations

import json
import pathlib
import subprocess
import sys


def last_json(path: pathlib.Path) -> dict:
    last = None
    for raw in path.read_text().splitlines():
        raw = raw.strip()
        try:
            obj = json.loads(raw)
        except json.JSONDecodeError:
            continue
        if isinstance(obj, dict) and obj.get("wav"):
            last = obj
    if last is None:
        sys.exit("bake printed no json")
    return last


def main() -> int:
    if len(sys.argv) != 7:
        sys.exit("usage: render-demo-tape.py ROOT WAV SAY_OUT SILENT TAPE LINE")
    root = pathlib.Path(sys.argv[1])
    wav = sys.argv[2]
    say_out = pathlib.Path(sys.argv[3])
    silent = sys.argv[4]
    tape = pathlib.Path(sys.argv[5])
    line = sys.argv[6]

    ttfa = int(last_json(say_out)["ttfa_ms"])
    dur = float(
        subprocess.check_output(
            [
                "ffprobe",
                "-v",
                "error",
                "-show_entries",
                "format=duration",
                "-of",
                "csv=p=0",
                wav,
            ],
            text=True,
        ).strip()
    )
    wav_ms = int(round(dur * 1000))
    type_ms, pre_ms, pause_ms = 55, 500, 350
    # VHS Type is slower than TypingSpeed * len; this lands audio on "speaking".
    sync_ms = 400
    offset_ms = pre_ms + len(line) * type_ms + pause_ms + sync_ms
    after_ms = wav_ms + 1500

    text = (root / "tapes" / "demo.tape.in").read_text()
    text = text.replace("Output docs/booth.mp4", f'Output "{silent}"')
    text = text.replace("__AFTER_ENTER__", f"{after_ms}ms")
    text = text.replace(
        'Env CANS_SAY_BIN "tapes/demo-say"',
        "\n".join(
            [
                'Env CANS_SAY_BIN "tapes/demo-say"',
                f'Env CANS_DEMO_TTFA "{ttfa}"',
                f'Env CANS_DEMO_HOLD_MS "{wav_ms}"',
            ]
        ),
    )
    tape.write_text(text)
    (root / "tapes" / "demo.offset").write_text(str(offset_ms))
    print(f"ttfa {ttfa}ms  wav {wav_ms}ms  enter at {offset_ms}ms  hold {after_ms}ms")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
