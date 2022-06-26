#!/usr/bin/env bash
set -eo pipefail

# Generate Dgraph Go clients. This requires the Dgraph database to be running
# and listening on localhost:8080.
#
# Usage: ./tools/generate-dgraph-clients.sh

repo_root_dir="$(dirname "$0")/.."
cd "$repo_root_dir"

find . -type f -iname ".gqlgenc-dgraph.yml" -print0 |
    while read -r -d $'\0' filename; do
        echo "Generating client for $filename"
        ./bin/tools/gqlgenc -d "$(dirname "$filename")" .gqlgenc-dgraph.yml
    done
