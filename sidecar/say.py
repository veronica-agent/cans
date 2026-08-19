"""Clone a ref wav with Qwen3-TTS 0.6B Base. No veronica package import."""

from __future__ import annotations

import argparse
import json
import sys
import tempfile
import time
from pathlib import Path

BASE_MODEL = "mlx-community/Qwen3-TTS-12Hz-0.6B-Base-bf16"
CLONE_TEMPERATURE = 0.2
_STEPS_PER_S = 50
MAX_TOKENS = 360


def _token_budget(text: str) -> int:
    words = max(1, len((text or "").split()))
    return max(180, min(MAX_TOKENS, int(_STEPS_PER_S * (0.55 * words + 2.2))))


def main() -> int:
    p = argparse.ArgumentParser()
    p.add_argument("--text", required=True)
    p.add_argument("--ref", required=True)
    p.add_argument("--ref-text", default="")
    p.add_argument("--out", default="")
    args = p.parse_args()

    ref = Path(args.ref)
    if not ref.is_file():
        print(f"ref wav not found: {ref}", file=sys.stderr)
        return 2

    from mlx_audio.tts.utils import load_model  # type: ignore
    import numpy as np
    import soundfile as sf

    model = load_model(BASE_MODEL)
    started = time.perf_counter()
    first_chunk = None
    chunks: list = []
    rate = 24000
    for result in model.generate(
        text=args.text,
        ref_audio=str(ref),
        ref_text=args.ref_text or " ",
        lang_code="English",
        temperature=CLONE_TEMPERATURE,
        max_tokens=_token_budget(args.text),
    ):
        if first_chunk is None:
            first_chunk = time.perf_counter()
        audio = getattr(result, "audio", result)
        arr = np.array(audio, dtype=np.float32).reshape(-1)
        chunks.append(arr)
        rate = int(getattr(result, "sample_rate", rate) or rate)

    if not chunks:
        print("qwen3-tts produced no audio", file=sys.stderr)
        return 1

    pcm = np.concatenate(chunks)
    elapsed_first = (first_chunk or time.perf_counter()) - started
    out = Path(args.out) if args.out else Path(tempfile.mkstemp(suffix=".wav")[1])
    sf.write(str(out), pcm, rate)

    print(
        json.dumps(
            {
                "wav": str(out),
                "ttfa_ms": int(elapsed_first * 1000),
                "sample_rate": rate,
            }
        ),
        flush=True,
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
