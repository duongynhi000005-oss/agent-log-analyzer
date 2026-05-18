#!/usr/bin/env sh
set -eu

cd "$(dirname "$0")/.."
mkdir -p .data
go run ./cmd/signature-research \
  -config docs/signature-research-sources.json \
  -out .data/signature-candidates.json \
  "$@"
