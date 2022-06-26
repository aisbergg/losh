#!/bin/sh

# Activate the development environment.
#
# Usage: source ./tools/activate.sh

script_path=$(dirname $0)/
bin_path="$(readlink -m "$script_path/../bin")"
tools_bin_path="$bin_path/tools"
export GOBIN="$bin_path"
export PATH="$tools_bin_path:$bin_path:$PATH"
