#!/usr/bin/env just --justfile
# cans — local mouth, Charm booth

set dotenv-load := false

binary_name := "cans"
bin_dir := "bin"

export CANS_ROOT := justfile_directory()

[doc('Build the cans binary')]
mod build '.justfiles/build.just'

[doc('Tests')]
mod test '.justfiles/test.just'

[doc('Vet / fmt')]
mod lint '.justfiles/lint.just'

[doc('Record the booth with VHS')]
mod vhs '.justfiles/vhs.just'

[doc('Run booth / say / keep')]
mod run '.justfiles/run.just'

[doc('Stamp embed, goreleaser snapshot')]
mod dist '.justfiles/dist.just'

[doc('Build bin/cans and copy it onto PATH')]
install:
    #!/usr/bin/env bash
    set -euo pipefail
    just build quick
    dest="${DEST:-}"
    if [ -z "$dest" ]; then
        if [ -d /opt/homebrew/bin ] && [ -w /opt/homebrew/bin ]; then
            dest=/opt/homebrew/bin
        else
            dest="${HOME}/.local/bin"
        fi
    fi
    mkdir -p "$dest"
    install -m 755 "{{ bin_dir }}/cans" "$dest/cans"
    echo "installed $dest/cans"

[private]
default:
    @echo "cans — put the cans on."
    @echo ""
    @just --list --unsorted

fmt:
    gofmt -w .

vet:
    go vet ./...

clean:
    rm -rf {{ bin_dir }}
    @echo "cleaned {{ bin_dir }}"
