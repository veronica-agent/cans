set positional-arguments

root := `git rev-parse --show-toplevel`

[private]
default:
    @just --list --justfile {{ source_file() }}

version := `git describe --tags --always --dirty`
ldflags := "-s -w -X github.com/veronica-agent/cans/internal/ship.Version=" + version

# Build ./bin/cans
quick:
    mkdir -p {{ root }}/bin
    cd {{ root }} && go build -trimpath -ldflags "{{ ldflags }}" -o bin/cans ./cmd/cans

# go install onto GOBIN, or GOPATH/bin
install:
    #!/usr/bin/env bash
    set -euo pipefail
    cd "{{ root }}"
    go install -trimpath -ldflags "{{ ldflags }}" ./cmd/cans
    dest="$(go env GOBIN)"
    if [ -z "$dest" ]; then
        dest="$(go env GOPATH)/bin"
    fi
    echo "installed ${dest}/cans"

# Remove every cans install: go/brew bins, ./bin, and ~/.cans.
uninstall:
    #!/usr/bin/env bash
    set -euo pipefail
    dest="$(go env GOBIN)"
    if [ -z "$dest" ]; then
        dest="$(go env GOPATH)/bin"
    fi
    for p in "${dest}/cans" "{{ root }}/bin/cans"; do
        if [ -e "$p" ] || [ -L "$p" ]; then
            rm -f "$p"
            echo "removed $p"
        fi
    done
    if command -v brew >/dev/null 2>&1 && brew list --formula cans >/dev/null 2>&1; then
        brew uninstall --formula --force cans
    fi
    home="${CANS_HOME:-$HOME/.cans}"
    if [ -e "$home" ]; then
        rm -rf "$home"
        echo "removed ${home}"
    fi

uninstall-home:
    just uninstall

# Rebuild the quiet worker overlay against the installed v0.1.0 dylib.
overlay-worker:
    #!/usr/bin/env bash
    set -euo pipefail
    root="$(git rev-parse --show-toplevel)"
    lib="${CANS_HOME:-$HOME/.cans}/native/bin"
    src="$root/internal/ship/overlay"
    out="$src/qwen3-tts-worker-darwin-arm64"
    cc -O2 -Wall -Wextra -I"$src" -L"$lib" -Wl,-rpath,@loader_path \
      -o "$out" "$src/worker_main.c" -lqwen3tts -lm -lpthread
    install_name_tool -add_rpath @loader_path "$out" 2>/dev/null || true
    strip "$out"
    echo "wrote $out"
