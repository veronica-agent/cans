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

# Real mouth + VHS terminal, muxed to docs/booth.mp4 (gif stays silent).
demo:
    #!/usr/bin/env bash
    set -euo pipefail
    just vhs doctor
    command -v uv >/dev/null 2>&1 || { echo "vhs demo: need uv" >&2; exit 1; }
    just build quick
    repo_root="$(git rev-parse --show-toplevel)"
    cd "$repo_root"
    chmod +x tapes/demo-say

    line="Put the cans on."
    wav="$repo_root/tapes/booth-line.wav"
    export UV_PROJECT_ENVIRONMENT="${CANS_VENV:-$HOME/.cans/venv}"
    echo "baking \"$line\"…"
    say_out="$(mktemp -t cans-say)"
    uv run --project "$repo_root" --locked python sidecar/say.py \
        --text "$line" \
        --ref "$repo_root/voices/veronica/ref.wav" \
        --ref-text "Just like that, feel the rhythm of my voice." \
        --out "$wav" | tee "$say_out"

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

