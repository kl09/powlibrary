#!/bin/sh

set -e
set -u

if [ $# -ne 0 ]; then
	echo "Usage: $0" >&2
	exit 2
fi

go install github.com/bufbuild/buf/cmd/buf@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

buf format -w
buf generate
