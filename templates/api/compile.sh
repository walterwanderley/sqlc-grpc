#!/bin/sh
set -u
set -e
echo "Generating protocol buffer..."
protoc -I. -Ideps --go_out . --go_opt paths=source_relative --go-grpc_out . --go-grpc_opt paths=source_relative service.proto
