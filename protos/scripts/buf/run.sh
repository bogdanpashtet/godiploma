#!/usr/bin/env bash
set \
  -o errexit \
  -o nounset \
  -o pipefail

docker run --rm \
  -v "$(pwd)":/worksp \
  -e GOCACHE=/worksp/.cache/go-build \
  -e GOMODCACHE=/worksp/.cache/go/pkg/mod \
  -u "$(id -u):$(id -g)" \
  protos/buf "$@"
