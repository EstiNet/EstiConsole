#!/usr/bin/env bash
protoc -I protocol protocol/rpcserver.proto --go_out=plugins=grpc:protocol