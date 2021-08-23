#!/bin/bash

COMMIT_SHA=$(git rev-parse HEAD)
VERSION=$(git describe --tags)
BUILD_TIME=$(date +'%Y-%m-%d %T')

if which systeminfo >/dev/null; then
  BUILD_OS="$(systeminfo | grep "OS Name:" | sed -e "s/OS Name://" -e "s/  //g" -e "s/ //")"
elif which lsb_release >/dev/null; then
  BUILD_OS="$(lsb_release -i -s) $(lsb_release -r -s)"
else
  BUILD_OS="null"
fi

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
    -X 'main.VERSION=${VERSION}' \
    -X 'main.COMMIT_SHA=${COMMIT_SHA}' \
    -X 'main.BUILD_TIME=${BUILD_TIME}' \
    -X 'main.VERSION=${VERSION}' \
    -X 'main.BUILD_OS=${BUILD_OS}' \
"

go build ${OUTPUT_ARG} -trimpath -ldflags "${LDFlags}"
