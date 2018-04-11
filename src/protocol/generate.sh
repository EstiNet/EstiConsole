#!/usr/bin/env bash
export PATH=$PATH:$GOPATH/bin
protoc -I protocol protocol/rpcserver.proto --go_out=plugins=grpc:protocol
