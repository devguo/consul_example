#!/usr/bin/env bash

PROJECT_ROOT=${GOPATH}"src/github.com/devguo/consul_example"

SRC=$PROJECT_ROOT/proto/
TAR=$PROJECT_ROOT/pkg/pb/
protoc -I $SRC --go_out=plugins=grpc:$TAR $SRC/*.proto