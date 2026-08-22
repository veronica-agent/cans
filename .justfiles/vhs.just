set positional-arguments

root := `git rev-parse --show-toplevel`

[private]
default:
    @just --list --justfile {{ source_file() }}

doctor:
    #!/usr/bin/env bash
    set -euo pipefail
    for tool in vhs ttyd ffmpeg; do
        command -v "$tool" >/dev/null 2>&1 || {
            echo "vhs: missing required tool: $tool" >&2
            exit 1
        }
    done
    vhs --version

record tape:
    #!/usr/bin/env bash
    set -euo pipefail
    just vhs doctor
    repo_root="$(git rev-parse --show-toplevel)"
    tape="{{ tape }}"
    [[ "$tape" = /* ]] || tape="$repo_root/$tape"
    cd "$repo_root"
    exec vhs "$tape"

booth:
    just build quick
    chmod +x {{ root }}/tapes/fake-say
    just vhs record tapes/booth.tape

# Real mouth, no audio: a script pipes lines in, wavs land in out/.
pipe:
    just build quick
    just vhs record tapes/pipe.tape

# Real mouth + VHS terminal, muxed to docs/booth.mp4 (gif stays silent).
# The line is baked with qwen3-tts-cli, the one-shot tool that ships next to
# the worker (cans doctor puts both in ~/.cans/native/bin).
demo:
    #!/usr/bin/env bash
    set -euo pipefail
    just vhs doctor
    just build quick
    repo_root="$(git rev-parse --show-toplevel)"
    cd "$repo_root"
    chmod +x tapes/demo-say

    cans_home="${CANS_HOME:-$HOME/.cans}"
    worker="${CANS_WORKER_BIN:-$cans_home/native/bin/qwen3-tts-worker}"
    models="${CANS_WORKER_MODELS:-$cans_home/native/models}"
    cli="$(dirname "$worker")/qwen3-tts-cli"
    [[ -x "$cli" ]] || { echo "vhs demo: $cli missing — run cans doctor" >&2; exit 1; }

    line="Put the cans on."
    wav="$repo_root/tapes/booth-line.wav"
    echo "baking \"$line\"…"
    bake_log="$(mktemp -t cans-bake)"
    "$cli" -m "$models" -t "$line" -r "$repo_root/voices/veronica/ref.wav" -o "$wav" >/dev/null 2>"$bake_log"
    ms="$(sed -n 's/^ *Total: *\([0-9][0-9]*\) ms.*/\1/p' "$bake_log" | tail -n 1)"
    rm -f "$bake_log"
    say_out="$(mktemp -t cans-say)"
    printf '{"wav":"%s","ttfa_ms":%s,"sample_rate":24000}\n' "$wav" "${ms:-240}" | tee "$say_out"

    silent="$repo_root/docs/booth-silent.mp4"
    tape="$repo_root/tapes/demo.tape"
    python3 "$repo_root/tapes/render-demo-tape.py" \
        "$repo_root" "$wav" "$say_out" "$silent" "$tape" "$line"

    vhs "$tape"
    offset_ms="$(cat "$repo_root/tapes/demo.offset")"
    ffmpeg -y -i "$silent" -i "$wav" \
        -filter_complex "[1:a]aformat=channel_layouts=stereo,adelay=${offset_ms}|${offset_ms},apad[a]" \
        -map 0:v -map "[a]" \
        -c:v libx264 -pix_fmt yuv420p -c:a aac -b:a 192k \
        -movflags +faststart -shortest \
        "$repo_root/docs/booth.mp4"
    rm -f "$silent" "$say_out"
    echo "wrote docs/booth.mp4  docs/booth.gif"

