set positional-arguments

root := `git rev-parse --show-toplevel`

[private]
default:
    @just --list --justfile {{ source_file() }}

# Build ./bin/cans
quick:
    mkdir -p {{ root }}/bin
    cd {{ root }} && go build -o bin/cans ./cmd/cans
