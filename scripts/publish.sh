#!/usr/bin/env bash
set -e
set -x
. /etc/profile

PKG_CONFIG_FILE=./nfpm.yaml
VERSION=$(./scripts/version.sh)

function pcgConfigure() {
    local ARCH="$1"
    local VERSION="$2"
    local PLATFORM="$3"
    local SOURCE="$4"
    local CONFIG="$5"
    SOURCE=${SOURCE//"/"/"\/"}
    sed -i"" "s/\${ARCH}/$ARCH/" $CONFIG
    sed -i"" "s/\${VERSION}/$VERSION/" $CONFIG
    sed -i"" "s/\${PLATFORM}/$PLATFORM/" $CONFIG
    sed -i"" "s/\${SOURCE}/$SOURCE/" $CONFIG
}

function pcgRelease() {
    local ARCH="$1" && ARCH=$(echo "$ARCH" |  tr '[:upper:]' '[:lower:]' )
    local VERSION="$2" && VERSION=$(echo "$VERSION" |  tr '[:upper:]' '[:lower:]' )
    local PLATFORM="$3" && PLATFORM=$(echo "$PLATFORM" |  tr '[:upper:]' '[:lower:]' )

    local BIN_PATH=./bin/$ARCH/$PLATFORM
    local RELEASE_DIR=./bin/deb/$PLATFORM
    
    local ldfBuildTags="-X github.com/cosmos/cosmos-sdk/version.BuildTags=${PLATFORM} ${ARCH}"
    mkdir -p $BIN_PATH $RELEASE_DIR

    echoInfo "INFO: Building $ARCH package for $PLATFORM..."
    
    TMP_PKG_CONFIG_FILE=./nfpm_${ARCH}_${PLATFORM}.yaml
    rm -rfv $TMP_PKG_CONFIG_FILE && cp -v $PKG_CONFIG_FILE $TMP_PKG_CONFIG_FILE

    if [ "$PLATFORM" != "windows" ] ; then
        local RELEASE_PATH="${RELEASE_DIR}/sekai_${VERSION}_${ARCH}.deb"
        ./scripts/build.sh "${PLATFORM}" "${ARCH}" "$BIN_PATH/sekaid"
        pcgConfigure "$ARCH" "$VERSION" "$PLATFORM" "$BIN_PATH" $TMP_PKG_CONFIG_FILE
        nfpm pkg --packager deb --target "$RELEASE_PATH" -f $TMP_PKG_CONFIG_FILE
        cp -fv "$RELEASE_PATH" ./bin/sekai-${PLATFORM}-${ARCH}.deb
    else
        ./scripts/build.sh "${PLATFORM}" "${ARCH}" "$BIN_PATH/sekaid.exe"
        # deb is not supported on windows, simply copy the executables
        cp -fv $BIN_PATH/sekaid.exe ./bin/sekai-${PLATFORM}-${ARCH}.exe
    fi
}

rm -rfv ./bin

go mod tidy
GO111MODULE=on go mod verify

# NOTE: To see available build architectures, run: go tool dist list
pcgRelease "amd64" "$VERSION" "linux"
pcgRelease "amd64" "$VERSION" "darwin"
pcgRelease "amd64" "$VERSION" "windows"
pcgRelease "arm64" "$VERSION" "linux"
pcgRelease "arm64" "$VERSION" "darwin"
pcgRelease "arm64" "$VERSION" "windows"

cp -fv ./scripts/sekai-utils.sh ./bin/sekai-utils.sh
cp -fv ./scripts/sekai-env.sh ./bin/sekai-env.sh

rm -rfv ./bin/amd64 ./bin/arm64 ./bin/deb
echoInfo "INFO: Sucessfully published SEKAI deb packages into ./bin"
