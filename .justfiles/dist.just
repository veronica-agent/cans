set positional-arguments

root := `git rev-parse --show-toplevel`

[private]
default:
    @just --list --justfile {{ source_file() }}

# Copy checkout payload into the embed tree
stamp:
    mkdir -p {{ root }}/internal/ship/fs/sidecar {{ root }}/internal/ship/fs/voices/veronica
    cp -f {{ root }}/character.toml {{ root }}/pyproject.toml {{ root }}/uv.lock {{ root }}/internal/ship/fs/
    cp -f {{ root }}/sidecar/say.py {{ root }}/internal/ship/fs/sidecar/
    cp -f {{ root }}/voices/veronica/ref.wav {{ root }}/voices/veronica/meta.json {{ root }}/internal/ship/fs/voices/veronica/

# Validate .goreleaser.yaml
check:
    cd {{ root }} && goreleaser check

# Local snapshot tarball (no GitHub upload)
snapshot: stamp
    cd {{ root }} && goreleaser release --snapshot --clean --skip=publish
