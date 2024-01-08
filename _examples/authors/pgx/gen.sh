#!/bin/sh
set -u
set -e
set -x

go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

rm -rf internal proto api go.mod go.sum main.go registry.go buf*

sqlc generate
sqlc-grpc -m booktest -tracing
