#!/bin/bash
set -e

readonly env_file="$1"
env $(cat "./.env" "./$env_file" | grep -Ev '^#' | xargs) true

# If CI environment → use plain go test
if [[ -n "$CI" ]]; then
    echo "CI detected → using go test"
    TEST_CMD="go test"
else
    # Local: use richgo if installed
    if command -v richgo &> /dev/null; then
        echo "Local environment → using richgo"
        TEST_CMD="richgo test"
    else
        echo "Local environment → richgo not installed, using go test"
        TEST_CMD="go test"
    fi
fi

# Run tests
env $(cat "./.env" "./$env_file" | grep -Ev '^#' | xargs) \
    $TEST_CMD -count=1 -p=8 -parallel=8 -race -v ./...
