#!/bin/bash
set -e

echo "Building Clingy Client..."

mkdir -p build

echo "Building UI (OpenTUI)..."
bun build ./ui/src/index.tsx --compile --outfile build/clingy-ui

echo "Building API (Go server)..."
cd api && go build -o ../build/clingy-api . && cd ..

echo "Building wrapper..."
go build -o build/clingy .

echo "✓ Build complete! Run: ./build/clingy"

