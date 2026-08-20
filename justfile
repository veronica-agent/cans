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

[doc('Run, say, keep, install Python sidecar')]
mod run '.justfiles/run.just'

[doc('Stamp embed, goreleaser snapshot')]
mod dist '.justfiles/dist.just'

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
