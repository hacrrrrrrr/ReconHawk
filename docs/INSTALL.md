# Install ReconHawk

## Build from source

Requires Go 1.22+.

    git clone https://github.com/hacrrrrrrr/ReconHawk.git
    cd ReconHawk
    go build -o reconhawk ./cmd/reconhawk
    ./reconhawk --help

## Install with Go

    go install github.com/hacrrrrrrr/ReconHawk/cmd/reconhawk@latest

Then ensure your Go bin directory is on PATH:

    export PATH="$PATH:$(go env GOPATH)/bin"

Verify:

    reconhawk version
    reconhawk --help

## Download releases

When a GitHub Release is published, users can download the platform-specific binary from the Releases page instead of compiling.

Only use ReconHawk against assets you own or are explicitly authorized to test.
