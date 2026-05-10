#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
UI_DIR="$(cd "$SCRIPT_DIR/../ui" && pwd)"
API_DIR="$(cd "$SCRIPT_DIR/../api" && pwd)"
BUILD_DIR="$SCRIPT_DIR"

echo "Building Clingy Client..."

echo "Building UI (OpenTUI)..."
bun build "$UI_DIR/src/index.tsx" --compile --outfile "$BUILD_DIR/clingy-ui"

echo "Building API (Go server)..."
cd "$API_DIR" && go build -o "$BUILD_DIR/clingy-api" .

echo "Building wrapper..."
cd "$BUILD_DIR" && go build -o clingy .

echo "Cleaning up..."
rm "$BUILD_DIR/clingy-ui"
rm "$BUILD_DIR/clingy-api"

echo "✓ Build complete! Run: $BUILD_DIR/clingy"

