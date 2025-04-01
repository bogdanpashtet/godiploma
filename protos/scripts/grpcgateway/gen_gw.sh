#!/usr/bin/env bash
set \
  -o errexit \
  -o nounset \
  -o pipefail

gen_gw() {
  buf generate \
    --template ./protos/proto/client/godiploma/file/v1/file_api.buf.gen.yaml \
    --path ./protos/proto/client/godiploma/file/v1/file_api.proto
}

merge_gw(){
  swagger mixin -q \
    ./protos/third_party/OpenAPI/client/godiploma/file/v1/*.json \
    ./protos/third_party/OpenAPI/client/godiploma/file/v1/*.json \
    -o ./protos/third_party/OpenAPI/client/godiploma/swagger.json || true
}

gen_gw
merge_gw
