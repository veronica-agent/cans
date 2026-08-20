set positional-arguments

root := `git rev-parse --show-toplevel`

[private]
default:
    @just --list --justfile {{ source_file() }}

# Unit tests (no audio, no MLX)
unit:
    cd {{ root }} && CANS_NOPLAY=1 go test ./...

v:
    cd {{ root }} && CANS_NOPLAY=1 go test -v ./...
