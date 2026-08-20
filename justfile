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

[doc('Install, run booth / say / keep')]
mod run '.justfiles/run.just'

[doc('Stamp embed, goreleaser snapshot')]
mod dist '.justfiles/dist.just'

[doc('go install cans onto GOBIN / GOPATH/bin')]
install:
    just build install

[doc('remove cans binaries; --all also wipes ~/.cans')]
uninstall *args:
    just build uninstall {{ args }}

[doc('same as just uninstall --all')]
uninstall-home:
    just uninstall --all

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
