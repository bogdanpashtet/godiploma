ARG GOLANG_VERSION=1.24

FROM golang:${GOLANG_VERSION}

ARG PROTOC_GEN_GO_VERSION=v1.30.0
ARG PROTOC_GEN_GO_GRPC_VERSION=v1.3.0
ARG BUF_VERSION=v1.17.0
ARG PROTOC_GATEWAY_VERSION=v2.15.2
ARG SWAGGER_VERSION=v0.30.4
ARG PROTOC_GEN_DOC_VERSION=v1.5.1

ENV GOPATH=/go

RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@${PROTOC_GATEWAY_VERSION}
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@${PROTOC_GATEWAY_VERSION}

RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION}
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION}

RUN go install github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@${BUF_VERSION}
RUN go install github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@${BUF_VERSION}
RUN go install github.com/bufbuild/buf/cmd/buf@${BUF_VERSION}
RUN go install github.com/go-swagger/go-swagger/cmd/swagger@${SWAGGER_VERSION}
RUN go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@${PROTOC_GEN_DOC_VERSION}

# Setup for python
ARG PYTHON_MYPY_PROTOBUF_VERSION=3.6.0
ARG PYTHON_GRPCIO_VERSION=1.60.0

RUN apt-get update && \
    apt-get -y install git python3 python3-pip protobuf-compiler && \
    python3 -m pip install --break-system-packages \
        grpcio==${PYTHON_GRPCIO_VERSION} \
        grpcio-tools==${PYTHON_GRPCIO_VERSION} \
        mypy-protobuf==${PYTHON_MYPY_PROTOBUF_VERSION}

WORKDIR "/worksp"

ENV BUF_CACHE_DIR="/worksp/.cache/buf"
