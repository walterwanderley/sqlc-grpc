#!/bin/sh
set -u
set -e
set -x

go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

rm -rf internal proto go.mod go.sum main.go registry.go

sqlc generate
sqlc-grpc -m booktest
