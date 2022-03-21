#!/usr/bin/env bash
set -e
set -x
. /etc/profile

TEST_NAME="TOKEN-TRANSFERS"
timerStart $TEST_NAME
echoInfo "INFO: $TEST_NAME - Integration Test - START"

echoInfo "INFO: Creating, deleting and testing accounts recovery..."
addAccount tester1
addAccount tester2
ACCOUNT3_MNEMONIC=$(addAccount tester3 | jq .mnemonic | xargs)
ACCOUNT3_ADDRESS=$(showAddress tester3)

[ "$(isAccount tester3)" != "true" ] && echoErr "ERROR: Expected account 'tester3' to exist, but it was NOT found" && exit 1
deleteAccount tester3
[ "$(isAccount tester3)" != "false" ] && echoErr "ERROR: Expected account 'tester3' to NOT exist, but it was found" && exit 1

recoverAccount tester3 "$ACCOUNT3_MNEMONIC"

ACCOUNT3_RECOVERED_ADDRESS=$(showAddress tester3)
[ "$ACCOUNT3_ADDRESS" != "$ACCOUNT3_RECOVERED_ADDRESS" ] && echoErr "ERROR: Expected account 'tester3' to be '$ACCOUNT3_ADDRESS' after recovery but got '$ACCOUNT3_RECOVERED_ADDRESS'" && exit 1

# TODO: Add token transfer tests

echoInfo "INFO: $TEST_NAME - Integration Test - END, elapsed: $(prettyTime $(timerSpan $TEST_NAME))"