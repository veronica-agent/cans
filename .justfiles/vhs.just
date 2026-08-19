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
    just vhs record tapes/booth.tape
