#!/bin/bash
set -e

echo "Building Clingy Client..."

echo "Building UI (OpenTUI)..."
bun build ../ui/src/index.tsx --compile --outfile ./clingy-ui

echo "Building API (Go server)..."
cd ../api && go build -o ../build/clingy-api . && cd ../build

echo "Building wrapper..."
go build -o ./clingy .

echo "✓ Build complete! Run: ./build/clingy"

