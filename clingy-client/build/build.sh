#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
UI_DIR="$(cd "$SCRIPT_DIR/../ui" && pwd)"
API_DIR="$(cd "$SCRIPT_DIR/../api" && pwd)"
BUILD_DIR="$SCRIPT_DIR"

PORT="8888"
while [[ $# -gt 0 ]]; do
    case "$1" in
        -p|-port)
            if [ $# -lt 2 ] || [[ "$2" == -* ]]; then
                echo "Error: -p/-port requires a port number" >&2
                exit 1
            fi
            PORT="$2"
            shift 2 ;;
        *) echo "Error: Unknown option $1" >&2; exit 1 ;;
    esac
done

echo "Building Clingy Client on port $PORT..."

echo "Building UI (OpenTUI)..."
bun build "$UI_DIR/src/index.tsx" --compile --outfile "$BUILD_DIR/clingy-ui"

echo "Building API (Go server)..."
cd "$API_DIR" && go build -o "$BUILD_DIR/clingy-api" .

echo "Building wrapper with port $PORT..."
cd "$BUILD_DIR" && go build -ldflags="-X main.defaultPort=$PORT" -o clingy .

echo "Cleaning up..."
rm "$BUILD_DIR/clingy-ui"
rm "$BUILD_DIR/clingy-api"

echo "✓ Build complete! Run: $BUILD_DIR/clingy"

