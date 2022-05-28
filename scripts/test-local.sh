#!/usr/bin/env bash
set -e
set -x
. /etc/profile

echo "INFO: Started local tests in '$PWD'..."
timerStart

echoInfo "INFO: Cleanup system resources..."
kill -9 $(lsof -t -i:9090) || echoWarn "WARNING: Nothing running on port 9090, or failed to kill processes"
kill -9 $(lsof -t -i:6060) || echoWarn "WARNING: Nothing running on port 6060, or failed to kill processes"
kill -9 $(lsof -t -i:26656) || echoWarn "WARNING: Nothing running on port 26656, or failed to kill processes"
kill -9 $(lsof -t -i:26657) || echoWarn "WARNING: Nothing running on port 26657, or failed to kill processes"
kill -9 $(lsof -t -i:26658) || echoWarn "WARNING: Nothing running on port 26658, or failed to kill processes"
kill -9 $(lsof -t -i:11000) || echoWarn "WARNING: Nothing running on port 11000, or failed to kill processes"

echoInfo "INFO: Installing latest sekai-utils release..."
./scripts/sekai-utils.sh sekaiUtilsSetup
loadGlobEnvs

echoInfo "INFO: Ensuring correct sekaid version is installed..."
SEKAID_VERSION=$(sekaid version)
SEKAID_EXPECTED_VERSION=$(./scripts/version.sh)

[ "$SEKAID_VERSION" != "$SEKAID_EXPECTED_VERSION" ] && \
    echoErr "ERROR: Expected installed sekaid version to be $SEKAID_EXPECTED_VERSION, but got $SEKAID_VERSION, try to make build & install first" && exit 1

echoInfo "INFO: Stopping local network..."
./scripts/test-local/network-stop.sh || ( systemctl2 stop sekai && exit 1 )

echoInfo "INFO: Launching local network..."
./scripts/test-local/network-start.sh || ( systemctl2 stop sekai && exit 1 )

echoInfo "INFO: Testing wallets & transfers..."
./scripts/test-local/token-transfers.sh || ( systemctl2 stop sekai && exit 1 )

echoInfo "INFO: Testing account permissions whitelist, blacklist & clear..."
./scripts/test-local/account-permissions.sh || ( systemctl2 stop sekai && exit 1 )

echoInfo "INFO: Stopping local network..."
./scripts/test-local/network-stop.sh || ( systemctl2 stop sekai && exit 1 )

echoInfo "INFO: Success, all local tests passed, elapsed: $(prettyTime $(timerSpan))"