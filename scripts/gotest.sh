#!/bin/bash

set -e

echo "Running tests..."
PKGS=`go list github.com/ewangplay/serval/...`
go test -cover -p 1 -timeout=20m $PKGS
