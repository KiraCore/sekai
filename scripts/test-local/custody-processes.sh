#!/usr/bin/env bash
# To run test locally: make network-start && ./scripts/test-local/token-transfers.sh
set -e
set -x
. /etc/profile
. ./scripts/sekai-env.sh

TEST_NAME="CUSTODY" && timerStart $TEST_NAME
echoInfo "INFO: $TEST_NAME - Integration Test - START"

echoInfo "INFO: Creating, deleting and testing accounts recovery..."
addAccount tester1
addAccount tester2
addAccount tester3

ACCOUNT1_ADDRESS=$(showAddress tester1)
ACCOUNT2_ADDRESS=$(showAddress tester2)
ACCOUNT3_ADDRESS=$(showAddress tester3)
INPUT_CUSTODIANS="$ACCOUNT2_ADDRESS,$ACCOUNT3_ADDRESS"

KEY1="6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b"
KEY11="e0bc614e4fd035a488619799853b075143deea596c477b8dc077e309c0fe42e9"
KEY2="d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35"
KEY21="d8bdf9a0cb27a193a1127de2924b6e5a9e4c2d3b3fe42e935e160c011f3df1fc"
KEY3="4e07408562bedb8b60ce05c1decfe3ad16b72230967de01f640b7e4729b49fce"
KEY31="5b65712d565c1551340998102d418ceccb35db8dbfb45f9041c4cae483d8717b"

sendTokens validator $(showAddress tester1) 1000000000000 ukex 100 ukex

enableCustody tester1 75 0 0 0 150 ukex 0 $KEY11
CUSTODY_KEY1=$(getCustodyKey $ACCOUNT1_ADDRESS)

[ "$CUSTODY_KEY1" != $KEY11 ] && echoErr "ERROR: Expected key to be $KEY11, but got '$CUSTODY_KEY1'" && exit 1

enableCustody tester1 50 0 0 0 150 ukex $KEY1 $KEY21
CUSTODY_KEY2=$(getCustodyKey $ACCOUNT1_ADDRESS)

[ "$CUSTODY_KEY2" != $KEY21 ] && echoErr "ERROR: Expected key to be $KEY21, but got '$CUSTODY_KEY2'" && exit 1

addCustodians tester1 "$INPUT_CUSTODIANS" 150 ukex $KEY2 $KEY31
RECEIVED_CUSTODIANS=$(getCustodians $ACCOUNT1_ADDRESS)
CUSTODY_KEY3=$(getCustodyKey $ACCOUNT1_ADDRESS)

[ "$CUSTODY_KEY3" != $KEY31 ] && echoErr "ERROR: Expected key to be $KEY31, but got '$CUSTODY_KEY3'" && exit 1
#[ "$INPUT_CUSTODIANS" != "$RECEIVED_CUSTODIANS" ] && echoErr "ERROR: Expected custodians to be '$INPUT_CUSTODIANS', but got '$RECEIVED_CUSTODIANS'" && exit 1

echoInfo "INFO: $TEST_NAME - Integration Test - END, elapsed: $(prettyTime $(timerSpan $TEST_NAME))"
