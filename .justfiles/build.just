set positional-arguments

root := `git rev-parse --show-toplevel`

[private]
default:
    @just --list --justfile {{ source_file() }}

version := `git describe --tags --always --dirty`

# Build ./bin/cans
quick:
    mkdir -p {{ root }}/bin
    cd {{ root }} && go build -trimpath -ldflags "-s -w -X github.com/veronica-agent/cans/internal/ship.Version={{ version }}" -o bin/cans ./cmd/cans
