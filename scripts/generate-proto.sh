#!/bin/bash

# Exit on error
set -e

# Ensure output directory exists
mkdir -p api/gen/v1

# Generate the code
protoc \
    --proto_path=proto \
    --proto_path=proto/include \
    --go_out=api/gen \
    --go_opt=paths=source_relative \
    --go-grpc_out=api/gen \
    --go-grpc_opt=paths=source_relative \
    proto/v1/agent.proto

echo "Agent generation completed successfully"