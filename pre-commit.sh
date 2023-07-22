#!/bin/bash

# Run linter and formatter
golangci-lint run
go fmt ./...

# Run tests (optional)
#go test ./...

# Check for any changes (optional)
git diff --exit-code
