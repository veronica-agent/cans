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
