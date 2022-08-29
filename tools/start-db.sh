#!/usr/bin/env bash
set -eo pipefail

# Starts Docker daemon and the Dgraph database via Docker-Compose.
#
# Usage: ./tools/start-db.sh

cd "$(dirname "$0")/../deployments"
sudo systemctl start docker
docker-compose -f docker-compose.dev.yml up -d
