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

# Remove go/brew/./bin cans. `just uninstall --all` also wipes ~/.cans (native mouth).
uninstall *args:
    #!/usr/bin/env bash
    set -euo pipefail
    # Shebang recipes do not get *args on $@ unless they are interpolated.
    set -- {{ args }}
    all=0
    for a in "$@"; do
        case "$a" in
            --all|all) all=1 ;;
        esac
    done
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
    if [ "$all" -eq 1 ]; then
        if [ -e "$home" ]; then
            rm -rf "$home"
            echo "removed ${home}"
        fi
    else
        echo "left ${home} (native mouth). wipe it with: just uninstall --all"
    fi

uninstall-home:
    just uninstall --all
