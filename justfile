set dotenv-load := false

export CANS_ROOT := justfile_directory()
export HF_HOME := env("HF_HOME", "/Volumes/fast/huggingface")
export HF_HUB_CACHE := env("HF_HUB_CACHE", HF_HOME + "/hub")

default:
    @just --list --justfile {{ source_file() }}

install:
    go mod tidy
    uv sync

run:
    go run ./cmd/cans

say *text:
    go run ./cmd/cans say {{ text }}

keep wav:
    go run ./cmd/cans keep {{ wav }}

test:
    CANS_NOPLAY=1 go test ./...

lint:
    go vet ./...

tape:
    #!/usr/bin/env bash
    set -euo pipefail
    if ! command -v vhs >/dev/null; then
      echo "vhs not installed. brew install vhs  # GIF is optional"
      exit 0
    fi
    vhs tapes/booth.tape
