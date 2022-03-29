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

RELEASE_FILE=./RELEASE.md
RELEASE_VERSION=$(grep -Fn -m 1 'Release: ' $RELEASE_FILE | rev | cut -d ":" -f1 | rev | xargs | tr -dc '[:alnum:]\-\.' || echo '')
RELEASE_LINE_NR=$(getFirstLineByPrefix "Release:" $RELEASE_FILE)

# If release file is not present or release version is NOT defined then create RELEASE.md or append the Release version
if ($(isNullOrEmpty "$RELEASE_VERSION")) || [ ! -f $RELEASE_FILE ] || [ $RELEASE_LINE_NR -le 0 ] ; then
    touch $RELEASE_FILE
    echo -e "\n\rRelease: \`$VERSION\`" >> $RELEASE_FILE
# Otherwsie replace release with the number defined by the constants file
else
    RELEASE_LINE_NR=$(getFirstLineByPrefix "Release:" $RELEASE_FILE)
    setLineByNumber $RELEASE_LINE_NR "Release: \`$VERSION\`" $RELEASE_FILE
fi

RELEASE_VERSION=$(grep -Fn -m 1 'Release: ' $RELEASE_FILE | rev | cut -d ":" -f1 | rev | xargs | tr -dc '[:alnum:]\-\.' || echo '')
[ "${VERSION}" != "${RELEASE_VERSION}" ] && echoErr "ERROR: Inconsistency between RELEASE.md version and SekaiVersion, expected '${VERSION}', but got '${RELEASE_VERSION}'" && exit 1

COMMIT=$(git log -1 --format='%H') && \
ldfName="-X github.com/cosmos/cosmos-sdk/version.Name=sekai" && \
ldfAppName="-X github.com/cosmos/cosmos-sdk/version.AppName=sekaid" && \
ldfVersion="-X github.com/cosmos/cosmos-sdk/version.Version=$VERSION" && \
ldfCommit="-X github.com/cosmos/cosmos-sdk/version.Commit=$COMMIT" && \
ldfBuildTags="-X github.com/cosmos/cosmos-sdk/version.BuildTags=${PLATFORM},${ARCH}"

rm -fv "$OUTPUT" || echo "ERROR: Failed to wipe old sekaid binary"

go mod tidy
GO111MODULE=on go mod verify
env GOOS=$PLATFORM GOARCH=$ARCH go build -ldflags "${ldfName} ${ldfAppName} ${ldfVersion} ${ldfCommit} ${ldfBuildTags}" -o "$OUTPUT" ./cmd/sekaid

( [ "$PLATFORM" == "$LOCAL_PLATFORM" ] && [ "$ARCH" == "$LOCAL_ARCH" ] && [ -f $OUTPUT ] ) && \
    echoInfo "INFO: Sucessfully built SEKAI $($OUTPUT version)" || echoInfo "INFO: Sucessfully built SEKAI to '$OUTPUT'"
