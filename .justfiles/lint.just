set positional-arguments

root := `git rev-parse --show-toplevel`

[private]
default:
    @just --list --justfile {{ source_file() }}

vet:
    cd {{ root }} && go vet ./...

all: vet
    cd {{ root }} && gofmt -l .
