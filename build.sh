#!/bin/bash

COMMIT_SHA=$(git rev-parse HEAD)
VERSION=$(git describe --tags)
BUILD_TIME=$(date +'%Y-%m-%d %T')

CUSTOM_GOOS=$1
CUSTOM_GOARCH=$2
OUTPUT_ARG=""

if [[ "$CUSTOM_GOOS" != "" ]]; then
  export GOOS="$CUSTOM_GOOS"
fi

if [[ "$CUSTOM_GOARCH" != "" ]]; then
  export GOARCH="$CUSTOM_GOARCH"
fi

if [[ "$3" == "-o" ]]; then
  export OUTPUT_ARG="-o $4"
fi

LDFlags="\
    -s -w
    -X 'main.version=${VERSION}' \
    -X 'main.commitSHA=${COMMIT_SHA}' \
    -X 'main.buildTime=${BUILD_TIME}' \
"

CGO_ENABLED=0 go build ${OUTPUT_ARG} -trimpath -ldflags "${LDFlags}"
