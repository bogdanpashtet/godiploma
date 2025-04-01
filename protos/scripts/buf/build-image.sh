#!/usr/bin/env bash
set \
  -o errexit \
  -o nounset \
  -o pipefail

docker build \
  -t protos/buf \
  -f protos/build/buf/Dockerfile \
  --load \
  protos/build/buf
