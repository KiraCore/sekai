#!/usr/bin/env bash
# To run test locally: make network-start && ./scripts/test-local/token-transfers.sh
set -e
set -x
. /etc/profile
. ./scripts/sekai-env.sh

TEST_NAME="TOKEN-TRANSFERS" && timerStart $TEST_NAME
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

echoInfo "INFO: Reading validator, faucet & tester-1 account ukex balances before token transfers"
VALIDATOR_BALANCE_START=$(showBalance validator ukex)
FAUCET_BALANCE_START=$(showBalance faucet ukex)
TESTER1_BALANCE_START=$(showBalance tester1 ukex)

echoInfo "INFO: Sending 1M KEX from validator to faucet"
sendTokens validator $(showAddress faucet) 1000000000000 ukex 100 ukex

echoInfo "INFO: Sending 7 ukex from faucet to tester1"
sendTokens faucet tester1 7 ukex 150 ukex

echoInfo "INFO: Reading validator, faucet & tester-1 account ukex balances after token transfers"
VALIDATOR_BALANCE_END=$(showBalance validator ukex)
FAUCET_BALANCE_END=$(showBalance faucet ukex)
TESTER1_BALANCE_END=$(showBalance tester1 ukex)

VALIDATOR_BALANCE_EXPECTED=$(($VALIDATOR_BALANCE_START - 1000000000000 - 100))
FAUCET_BALANCE_EXPECTED=$(($FAUCET_BALANCE_START + 1000000000000 - 7 - 150))
TESTER1_BALANCE_EXPECTED=7

[ "$VALIDATOR_BALANCE_EXPECTED" != "$VALIDATOR_BALANCE_END" ] && echoErr "ERROR: Expected validator account balance to be '$VALIDATOR_BALANCE_EXPECTED', but got '$VALIDATOR_BALANCE_END'" && exit 1
[ "$FAUCET_BALANCE_EXPECTED" != "$FAUCET_BALANCE_END" ] && echoErr "ERROR: Expected faucet account balance to be '$FAUCET_BALANCE_EXPECTED', but got '$FAUCET_BALANCE_END'" && exit 1
[ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_END" ] && echoErr "ERROR: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_END'" && exit 1

echoInfo "INFO: $TEST_NAME - Integration Test - END, elapsed: $(prettyTime $(timerSpan $TEST_NAME))"