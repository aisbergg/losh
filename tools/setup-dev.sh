#!/usr/bin/env bash
set -eo pipefail

# Install all dev dependencies and build dev tools.
#
# Usage: ./tools/setup-dev.sh

repo_root_dir="$(dirname "$0")/.."
tools_dir="$(dirname "$0")"
bin_dir="$repo_root_dir/bin/tools"
mkdir -p "$bin_dir"
pushd "$tools_dir" >/dev/null;

# download go dependencies
go mod download

# build custom GraphQL client generator
go build -o ../bin/tools/gqlgenc ./gqlgenc/main.go
go build -o ../bin/tools/codegen ./codegen/main.go

# build dev tools
export GOBIN"=$bin_dir"
# go install github.com/estesp/manifest-tool/v2/cmd/manifest-tool;

popd >/dev/null
