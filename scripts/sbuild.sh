#!/usr/bin/env bash
set -e
set -x
. /etc/profile

LOCAL_PLATFORM="$(uname)" && LOCAL_PLATFORM="$(echo "$LOCAL_PLATFORM" |  tr '[:upper:]' '[:lower:]' )"
LOCAL_ARCH=$(([[ "$(uname -m)" == *"arm"* ]] || [[ "$(uname -m)" == *"aarch"* ]]) && echo "arm64" || echo "amd64")
LOCAL_OUT="${GOBIN}/sekaid"

PLATFORM="$1" && [ -z "$PLATFORM" ] && PLATFORM="$LOCAL_PLATFORM"
ARCH="$2" && [ -z "$ARCH" ] && ARCH="$LOCAL_ARCH"
OUTPUT="$3" && [ -z "$OUTPUT" ] && OUTPUT="$LOCAL_OUT"

VERSION=$(./scripts/version.sh)
COMMIT=$(git log -1 --format='%H')
ldfName="-X github.com/cosmos/cosmos-sdk/version.Name=sekai"
ldfAppName="-X github.com/cosmos/cosmos-sdk/version.AppName=sekaid"
ldfVersion="-X github.com/cosmos/cosmos-sdk/version.Version=$VERSION"
ldfCommit="-X github.com/cosmos/cosmos-sdk/version.Commit=$COMMIT"
ldfBuildTags="-X github.com/cosmos/cosmos-sdk/version.BuildTags=${PLATFORM},${ARCH}"

rm -fv "$OUTPUT" || echo "ERROR: Failed to wipe old sekaid binary"

go mod tidy
GO111MODULE=on go mod verify

# Buld staic binary
env CGO_ENABLED=0 GOOS=$PLATFORM GOARCH=$ARCH go build -ldflags "-extldflags '-static' ${ldfName} ${ldfAppName} ${ldfVersion} ${ldfCommit} ${ldfBuildTags}" -o "$OUTPUT" ./cmd/sekaid

( [ "$PLATFORM" == "$LOCAL_PLATFORM" ] && [ "$ARCH" == "$LOCAL_ARCH" ] && [ -f $OUTPUT ] ) && \
    echo "INFO: Sucessfully built SEKAI $($OUTPUT version)" || echo "INFO: Sucessfully built SEKAI to '$OUTPUT'"
