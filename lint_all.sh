#!/bin/bash
set -euo pipefail

projects=("ntpnow" "calendar")

for dir in "${projects[@]}"; do
  echo "→ Running golangci-lint for project: $dir"
  cd "$dir"
  golangci-lint run ./...
  cd -
done
