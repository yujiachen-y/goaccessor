#!/bin/bash

# Set the script to fail if any command fails
set -e

# Step 1: run `go generate ./...` to generate code.
go generate ./...

# Step 2: run all test under 'test' folder and its sub folders, the test cases are written as `_test.go` files.
go test ./...

# Step 3: clear all generate code whose file suffix is `_goaccessor.go`.
find ./ -name '*_goaccessor.go' -delete
