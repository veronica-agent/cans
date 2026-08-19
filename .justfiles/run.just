set positional-arguments

root := `git rev-parse --show-toplevel`

[private]
default:
    @just --list --justfile {{ source_file() }}

install:
    cd {{ root }} && uv sync

booth:
    cd {{ root }} && go run ./cmd/cans

say *text:
    cd {{ root }} && go run ./cmd/cans say {{ text }}

keep *args:
    cd {{ root }} && go run ./cmd/cans keep {{ args }}
