#!/usr/bin/env bash
set -eo pipefail

# Upload the Dgraph database schema and generate Dgraph Go clients.
#
# Usage: ./tools/upload-db-schema.sh

tools_dir="$(dirname "$0")"
repo_root_dir="$(dirname "$0")/.."
cd "$repo_root_dir"

# upload the schema to the GraphQL server
curl -X POST localhost:8080/admin/schema --data-binary "@internal/models/.schema.graphql"
echo

# regenerate the GraphQL clients
. "$tools_dir/generate-dgraph-clients.sh"
