#!/usr/bin/env bash

set -euo pipefail

if [[ -d ../go-cache ]]; then
  GOPATH=$(realpath ../go-cache)
  export GOPATH
fi

GOOS="linux" GOARCH="amd64" go build -ldflags='-s -w' -o bin/main github.com/nncdevel-io/buildpack-application-config-environment-variable-labels/cmd/main
ln -fs main bin/build
ln -fs main bin/detect
chmod +x bin/*

echo build succeeded.
