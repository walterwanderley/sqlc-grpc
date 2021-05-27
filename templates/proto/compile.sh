#!/bin/sh
set -u
set -e

Compile () {
    rm -rf $1
    mkdir -p $1
    echo "Compiling $1.proto..."
    protoc -I. -Ivendor --go_out $1 --go_opt paths=source_relative --go-grpc_out $1 --go-grpc_opt paths=source_relative $1.proto
}

for i in *.proto; do
    pkg=$(echo "$i" | cut -f 1 -d '.')
    Compile $pkg
done    

echo "Finished!"
