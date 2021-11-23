#!/bin/sh
set -u
set -e

Compile () {
    rm -rf $1
    mkdir -p $1
    echo "Compiling $1.proto..." 
    buf generate --path $1.proto -o $1
}

for i in *.proto; do
    pkg=$(echo "$i" | cut -f 1 -d '.')
    Compile $pkg
done

echo "Generating OpenAPIv2 specs"
buf generate --template buf.doc.yaml

echo "Finished!"
