#!/bin/sh
set -u
set -e
set -x

rm -rf internal proto api go.mod go.sum main.go registry.go buf*

sqlc-grpc -m booktest -tracing -metric
