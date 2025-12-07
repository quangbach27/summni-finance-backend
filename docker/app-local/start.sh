#!/usr/bin/env sh
set -e
SRC_DIR=/app/cmd/server 
BINARY=/app/debug_server

echo "⚙️ Starting application..."

if [ "$DEBUG" = "true" ]; then
    echo "🔧 Debug mode enabled — starting Delve..."

    go build \
        -gcflags="all=-N -l" \
        -o $BINARY $SRC_DIR/main.go

    # Run Delve in headless mode
    dlv exec "$BINARY" \
        --headless \
        --listen=:40000 \
        --api-version=2 \
        --accept-multiclient \
        --log

else
    echo "🚀 Running server normally..."
    go run $SRC_DIR/main.go
fi