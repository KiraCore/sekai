#!/usr/bin/env bash
set -e
set -x
. /etc/profile

echo "INFO: Started local tests in '$PWD'..."
timerStart

echoInfo "INFO: Installing latest sekai-utils release..."
./scripts/sekai-utils.sh sekaiUtilsSetup
loadGlobEnvs

echoInfo "INFO: Ensuring correct sekaid version is installed..."
SEKAID_VERSION=$(sekaid version)
SEKAID_EXPECTED_VERSION=$(grep -Fn -m 1 'SekaiVersion ' ./types/constants.go | rev | cut -d "=" -f1 | rev | xargs | tr -dc '[:alnum:]\-\.' || echo '')
RELEASE_EXPECTED_VERSION=$(grep -Fn -m 1 'Release: ' ./RELEASE.md | rev | cut -d ":" -f1 | rev | xargs | tr -dc '[:alnum:]\-\.' || echo '')

[ "$SEKAID_VERSION" != "$SEKAID_EXPECTED_VERSION" ] && \
    echoErr "ERROR: Expected installed sekaid version to be $SEKAID_EXPECTED_VERSION, but got $SEKAID_VERSION, try to make build & install first" && exit 1
[ "$SEKAID_EXPECTED_VERSION" != "$RELEASE_EXPECTED_VERSION" ] && \
    echoErr "ERROR: Expected release version to be $SEKAID_EXPECTED_VERSION, but got $RELEASE_EXPECTED_VERSION, try to make build & install first" && exit 1

echoInfo "INFO: Launching local network..."
./scripts/test-local/network-setup.sh || ( systemctl2 stop sekai && exit 1 )

echoInfo "INFO: Testing wallets & transfers..."
./scripts/test-local/token-transfers.sh || ( systemctl2 stop sekai && exit 1 )

echoInfo "INFO: Stopping local network..."
systemctl2 stop sekai

echoInfo "INFO: Success, all local tests passed, elapsed: $(prettyTime $(timerSpan))"