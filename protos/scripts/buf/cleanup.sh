#!/usr/bin/env bash

cleanup() {
  # when moving contracts to other folders, delete the previous auto-generated code
  rm -rf protos/gen/go
  rm -rf protos/third_party/OpenAPI/client
}

cleanup
