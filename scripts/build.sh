#!/usr/bin/env bash
set -e
set -x
. /etc/profile


CONSTANS_FILE=./types/constants.go
SEKAI_VER=$(grep -Fn -m 1 'SekaiVersion ' $CONSTANS_FILE | rev | cut -d "=" -f1 | rev | xargs | tr -dc '[:alnum:]\-\.' || echo '')
($(isNullOrEmpty "$SEKAI_VER")) && ( echoErr "ERROR: SekaiVersion was NOT found in contants '$CONSTANS_FILE' !" && sleep 5 && exit 1 )

COMMIT=$(git log -1 --format='%H') && \
ldfName="-X github.com/cosmos/cosmos-sdk/version.Name=sekai" && \
ldfAppName="-X github.com/cosmos/cosmos-sdk/version.AppName=sekaid" && \
ldfVersion="-X github.com/cosmos/cosmos-sdk/version.Version=$SEKAI_VER" && \
ldfCommit="-X github.com/cosmos/cosmos-sdk/version.Name=$COMMIT"

rm -rfv "${GOBIN}/sekaid" || echo "ERROR: Failed to wipe old sekaid binary"

go mod tidy
GO111MODULE=on go mod verify
go build -ldflags "${ldfName} ${ldfAppName} ${ldfVersion} ${ldfCommit}" -o "${GOBIN}/sekaid" ./cmd/sekaid
echoInfo "INFO: Sucessfully intalled SEKAI $(sekaid version)"
