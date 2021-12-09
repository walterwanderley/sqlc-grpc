#!/bin/sh
set -u
set -e
set -x

go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

rm -rf internal proto go.mod go.sum main.go registry.go

sqlc generate
sqlc-grpc -m booktest
