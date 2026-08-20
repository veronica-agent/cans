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

# Remove go-installed / brew / ./bin cans, plus leftover uv venv.
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
    if [ -d "${home}/venv" ]; then
        rm -rf "${home}/venv"
        echo "removed ${home}/venv"
    fi
    echo "keep/shipped still in ${home} — just uninstall-home"

# Remove ~/.cans (keep, shipped, native, venv).
uninstall-home:
    #!/usr/bin/env bash
    set -euo pipefail
    home="${CANS_HOME:-$HOME/.cans}"
    if [ ! -e "$home" ]; then
        echo "no ${home}"
        exit 0
    fi
    rm -rf "$home"
    echo "removed ${home}"
