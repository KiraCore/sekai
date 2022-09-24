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
addAccount tester4
addAccount tester5

ACCOUNT1_ADDRESS=$(showAddress tester1)
ACCOUNT2_ADDRESS=$(showAddress tester2)
ACCOUNT3_ADDRESS=$(showAddress tester3)
ACCOUNT4_ADDRESS=$(showAddress tester4)
ACCOUNT5_ADDRESS=$(showAddress tester5)
INPUT_CUSTODIANS="$ACCOUNT2_ADDRESS,$ACCOUNT3_ADDRESS,$ACCOUNT4_ADDRESS"

PASSWORD="test_password"

KEY1="6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b"
KEY11="e0bc614e4fd035a488619799853b075143deea596c477b8dc077e309c0fe42e9"
KEY2="d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35"
KEY21="d8bdf9a0cb27a193a1127de2924b6e5a9e4c2d3b3fe42e935e160c011f3df1fc"
KEY3="4e07408562bedb8b60ce05c1decfe3ad16b72230967de01f640b7e4729b49fce"
KEY31="5b65712d565c1551340998102d418ceccb35db8dbfb45f9041c4cae483d8717b"
KEY4="4b227777d4dd1fc61c6f884f48641d02b4d121d3fd328cb08b5531fcacdabf8a"
KEY41="033c339a7975542785be7423a5b32fa8047813689726214143cdd7939747709c"
KEY5="ef2d127de37b942baad06145e54b0c619a1f22327b2ebbcfbec78f5564afe39d"
KEY51="c81d40dbeed369f1476086cf882dd36bf1c3dc35e07006f0bec588b983055487"
KEY6="e7f6c011776e8db7cd330b54174fd76f7d0216b612387a5ffcfb81e6f0919683"
KEY61="9e259b7f6b4c741937a96a9617b3e6b84e166ff6e925e414e7b72936f5a2a51f"

sendTokens validator "$ACCOUNT1_ADDRESS" 1000000000000 ukex 100 ukex
sendTokens validator "$ACCOUNT2_ADDRESS" 1000000000000 ukex 100 ukex
sendTokens validator "$ACCOUNT3_ADDRESS" 1000000000000 ukex 100 ukex
sendTokens validator "$ACCOUNT4_ADDRESS" 1000000000000 ukex 100 ukex
sendTokens validator "$ACCOUNT5_ADDRESS" 1000000000000 ukex 100 ukex

# Send tokens without enabled custodians and password protection
TESTER5_BALANCE_EXPECTED=1000000000007
custodySendTokens tester1 tester5 7 ukex 150 ukex $PASSWORD
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)

[ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR: Send tokens without enabled custodians: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
# ------------

# Enable custodians
enableCustody tester1 60 0 0 0 150 ukex 0 $KEY11
CUSTODY_KEY1=$(getCustodyKey tester1)

[ "$CUSTODY_KEY1" != $KEY11 ] && echoErr "ERROR: Expected key to be $KEY11, but got '$CUSTODY_KEY1'" && exit 1
# ------------

# Add custodians test double SHA protection
addCustodians tester1 "$INPUT_CUSTODIANS" 150 ukex $KEY1 $KEY21
CUSTODY_KEY2=$(getCustodyKey tester1)

[ "$CUSTODY_KEY2" != $KEY21 ] && echoErr "ERROR: Expected key to be $KEY21, but got '$CUSTODY_KEY2'" && exit 1
# ------------

# Send tokens with enabled custodians and not password protection
TESTER5_BALANCE_EXPECTED=1000000000007
TXHASH=$(custodySendTokens tester1 tester5 7 ukex 150 ukex $PASSWORD)
HASH=${TXHASH,,}
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
POOL=$(getCustodyPool tester1)

[ "$POOL" == "null" ] || [ "$POOL" == "{}" ] && echoErr "ERROR: Expected not empty pool" && exit 1
[ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR: Send tokens with enabled custodians: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
# ------------

# Approve transaction
approveTransaction tester2 tester1 "$HASH" 150 ukex
VOTES=$(getCustodyPoolVotes tester1 "$HASH")

[ "$VOTES" != 1 ] && echoErr "ERROR: Expected votes is '1', but got '$VOTES'" && exit 1
# ------------

# Decline transaction
declineTransaction tester3 tester1 "$HASH" 150 ukex
VOTES=$(getCustodyPoolVotes tester1 "$HASH")

[ "$VOTES" != 1 ] && echoErr "ERROR: Expected votes is '1', but got '$VOTES'" && exit 1
# ------------

# Approve transaction
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
approveTransaction tester4 tester1 "$HASH" 150 ukex
TESTER5_BALANCE_EXPECTED=1000000000014
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
POOL=$(getCustodyPool tester1)

[ "$POOL" != "null" ] || [ "$POOL" != "{}" ] && echoErr "ERROR: Expected empty pool" && exit 1
[ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR: Send tokens with enabled and approved custodians: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
# ------------

# Disable custodians and enable password protection
enableCustody tester1 0 1 0 0 150 ukex $KEY2 $KEY31
CUSTODY_KEY=$(getCustodyKey tester1)

[ "$CUSTODY_KEY" != $KEY31 ] && echoErr "ERROR: Expected key to be $KEY31, but got '$CUSTODY_KEY'" && exit 1
# ------------

# Send tokens with disabled custodians and password protection
TESTER5_BALANCE_EXPECTED=1000000000014
TXHASH=$(custodySendTokens tester1 tester5 7 ukex 150 ukex $PASSWORD)
HASH=${TXHASH,,}
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
POOL=$(getCustodyPool tester1)

[ "$POOL" == "null" ] || [ "$POOL" == "{}" ] && echoErr "ERROR: Expected not empty pool" && exit 1
[ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR: Send tokens with enabled custodians: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
# ------------

# Password confirm transaction
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
passwordConfirmTransaction tester4 tester1 "$HASH" "$PASSWORD" 150 ukex
TESTER5_BALANCE_EXPECTED=1000000000021
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
POOL=$(getCustodyPool tester1)

[ "$POOL" == "null" ] || [ "$POOL" == "{}" ] && echoErr "ERROR: Expected empty pool" && exit 1
[ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR: Send tokens with enabled and approved custodians: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
# ------------

# Drop custodians
dropCustodians tester1 150 ukex $KEY3 $KEY41
CUSTODY_KEY=$(getCustodyKey tester1)

[ "$CUSTODY_KEY" != $KEY41 ] && echoErr "ERROR: Expected key to be $KEY41, but got '$CUSTODY_KEY'" && exit 1
# ------------

# Enable custodians and password protection
enableCustody tester1 30 1 0 0 150 ukex $KEY4 $KEY51
CUSTODY_KEY=$(getCustodyKey tester1)

[ "$CUSTODY_KEY" != $KEY51 ] && echoErr "ERROR: Expected key to be $KEY51, but got '$CUSTODY_KEY'" && exit 1
# ------------

# Send tokens with enabled custodians and password protection
TESTER5_BALANCE_EXPECTED=1000000000021
TXHASH=$(custodySendTokens tester1 tester5 7 ukex 150 ukex $PASSWORD)
HASH=${TXHASH,,}
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
POOL=$(getCustodyPool tester1)

[ "$POOL" == "null" ] || [ "$POOL" == "{}" ]  && echoErr "ERROR: Expected not empty pool" && exit 1
[ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR: Send tokens with enabled custodians: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
# ------------

# Add custodians
addCustodians tester1 "$INPUT_CUSTODIANS" 150 ukex $KEY5 $KEY61
CUSTODY_KEY=$(getCustodyKey tester1)

[ "$CUSTODY_KEY" != $KEY61 ] && echoErr "ERROR: Expected key to be $KEY61, but got '$CUSTODY_KEY6'" && exit 1
# ------------

# Send tokens with enabled custodians and password protection
TESTER5_BALANCE_EXPECTED=1000000000021
TXHASH=$(custodySendTokens tester1 tester5 7 ukex 150 ukex $PASSWORD)
HASH=${TXHASH,,}
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
POOL=$(getCustodyPool tester1)

[ "$POOL" == "null" ] || [ "$POOL" == "{}" ] && echoErr "ERROR: Expected not empty pool" && exit 1
[ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR: Send tokens with enabled custodians: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
# ------------

# Approve transaction
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
approveTransaction tester4 tester1 "$HASH" 150 ukex
TESTER5_BALANCE_EXPECTED=1000000000021
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
POOL=$(getCustodyPool tester1)

[ "$POOL" != "null" ] || [ "$POOL" != "{}" ] && echoErr "ERROR: Expected empty pool" && exit 1
[ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR: Send tokens with enabled and approved custodians: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
# ------------

# Password confirm transaction
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
passwordConfirmTransaction tester4 tester1 "$HASH" "$PASSWORD" 150 ukex
TESTER5_BALANCE_EXPECTED=1000000000028
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
POOL=$(getCustodyPool tester1)

[ "$POOL" == "null" ] || [ "$POOL" == "{}" ] && echoErr "ERROR: Expected empty pool" && exit 1
[ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR: Send tokens with enabled and approved custodians: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
# ------------

# Send tokens with enabled custodians and password protection
TESTER5_BALANCE_EXPECTED=1000000000028
TXHASH=$(custodySendTokens tester1 tester5 7 ukex 150 ukex $PASSWORD)
HASH=${TXHASH,,}
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
POOL=$(getCustodyPool tester1)

[ "$POOL" == "null" ] || [ "$POOL" == "{}" ] && echoErr "ERROR: Expected not empty pool" && exit 1
[ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR: Send tokens with enabled custodians: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
# ------------

# Password confirm transaction
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
passwordConfirmTransaction tester4 tester1 "$HASH" "$PASSWORD" 150 ukex
TESTER5_BALANCE_EXPECTED=1000000000028
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
POOL=$(getCustodyPool tester1)

[ "$POOL" == "null" ] || [ "$POOL" == "{}}" ] && echoErr "ERROR: Expected not empty pool" && exit 1
[ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR: Send tokens with enabled and approved custodians: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
# ------------

# Approve transaction
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
approveTransaction tester4 tester1 "$HASH" 150 ukex
TESTER5_BALANCE_EXPECTED=1000000000035
TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
POOL=$(getCustodyPool tester1)

[ "$POOL" != "null" ] || [ "$POOL" != "{}" ] && echoErr "ERROR: Expected empty pool" && exit 1
[ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR: Send tokens with enabled and approved custodians: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
# ------------

echoInfo "INFO: $TEST_NAME - Integration Test - END, elapsed: $(prettyTime "$(timerSpan $TEST_NAME)")"
