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

sendTokens validator "$ACCOUNT1_ADDRESS" 1000000000000 ukex 100 ukex
sendTokens validator "$ACCOUNT2_ADDRESS" 1000000000000 ukex 100 ukex
sendTokens validator "$ACCOUNT3_ADDRESS" 1000000000000 ukex 100 ukex
sendTokens validator "$ACCOUNT4_ADDRESS" 1000000000000 ukex 100 ukex
sendTokens validator "$ACCOUNT5_ADDRESS" 1000000000000 ukex 100 ukex

# TEST 1
# Send tokens with disabled custodians and disabled password protection --->

  # TEST 1/1
  # ---> Send tokens
  echoInfo "TEST 1/1"
  TESTER1_BALANCE_EXPECTED=$((1000000000000 - 150 - 7)) # -150 fee -7 sent
  TESTER5_BALANCE_EXPECTED=$((1000000000000 + 7)) # +7 received

  custodySendTokens tester1 tester5 7 ukex 150 ukex 600 ukex $PASSWORD

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 1/1: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 1/1: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
   # <--- Send tokens

# <--- Send tokens with disabled custodians and enabled password protection

# TEST 2
# Send tokens with enabled custodians and disabled password protection --->

  # TEST 2/1
  # Enable custodians --->
  echoInfo "TEST 2/1"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 150)) # -150 fee

  enableCustody tester1 60 0 0 0 150 ukex 0

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 2/1: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  # <--- Enable custodians

  # TEST 2/2
  # Add custodians test double SHA protection --->
  echoInfo "TEST 2/2"
  addCustodiansForce tester1 "$INPUT_CUSTODIANS" 1 ukex
  addCustodiansForce tester1 "$INPUT_CUSTODIANS" 2 ukex

  sleep 10

  TESTER1_BALANCE_EXPECTED=$(showBalance tester1 ukex)

  # <--- Add custodians test double SHA protection

  # TEST 2/3
  # Send tokens --->
  echoInfo "TEST 2/3"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 150)) # -150 fee
  TESTER5_BALANCE_EXPECTED=$TESTER5_BALANCE_EXPECTED # same

  HASH=$(custodySendTokens tester1 tester5 7 ukex 150 ukex 600 ukex $PASSWORD)

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
  POOL=$(getCustodyPool tester1)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 2/3: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 2/3: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$POOL" == "null" ] || [ "$POOL" == "{}" ] && echoErr "ERROR TEST 2/3: Expected not empty pool" && exit 1
  # <--- Send tokens

  # TEST 2/4
  # Approve transaction --->
  echoInfo "TEST 2/4"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 200)) # -200 reward
  TESTER2_BALANCE_EXPECTED=$((1000000000000 - 150 + 200))  # -150 fee + 200 reward
  TESTER5_BALANCE_EXPECTED=$TESTER5_BALANCE_EXPECTED # same

  approveTransaction tester2 tester1 "$HASH" 150 ukex

  VOTES=$(getCustodyPoolVotes tester1 "$HASH")
  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER2_BALANCE_REAL=$(showBalance tester2 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 2/4: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER2_BALANCE_EXPECTED" != "$TESTER2_BALANCE_REAL" ] && echoErr "ERROR TEST 2/4: Expected tester2 account balance to be '$TESTER2_BALANCE_EXPECTED', but got '$TESTER2_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 2/4: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$VOTES" != 1 ] && echoErr "ERROR TEST 2/4: Expected votes is '1', but got '$VOTES'" && exit 1
  # <--- Approve transaction

  # TEST 2/4/1
  # Repeat approve transaction --->
  echoInfo "TEST 2/4/1"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED)) # -200 reward
  TESTER2_BALANCE_EXPECTED=$(($TESTER2_BALANCE_EXPECTED - 150))  # -150 fee + 200 reward
  TESTER5_BALANCE_EXPECTED=$TESTER5_BALANCE_EXPECTED # same

  approveTransaction tester2 tester1 "$HASH" 150 ukex

  VOTES=$(getCustodyPoolVotes tester1 "$HASH")
  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER2_BALANCE_REAL=$(showBalance tester2 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 2/4/1: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER2_BALANCE_EXPECTED" != "$TESTER2_BALANCE_REAL" ] && echoErr "ERROR TEST 2/4/1: Expected tester2 account balance to be '$TESTER2_BALANCE_EXPECTED', but got '$TESTER2_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 2/4/1: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$VOTES" != 1 ] && echoErr "ERROR TEST 2/4/1: Expected votes is '1', but got '$VOTES'" && exit 1
  # <--- Repeat approve transaction

  # TEST 2/5
  # Decline transaction --->
  echoInfo "TEST 2/5"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 200)) # -200 reward
  TESTER3_BALANCE_EXPECTED=$((1000000000000 - 150 + 200)) # -150 fee +200 reward
  TESTER5_BALANCE_EXPECTED=$TESTER5_BALANCE_EXPECTED # same

  declineTransaction tester3 tester1 "$HASH" 150 ukex

  VOTES=$(getCustodyPoolVotes tester1 "$HASH")
  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER3_BALANCE_REAL=$(showBalance tester3 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 2/5: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER3_BALANCE_EXPECTED" != "$TESTER3_BALANCE_REAL" ] && echoErr "ERROR TEST 2/5: Expected tester3 account balance to be '$TESTER3_BALANCE_EXPECTED', but got '$TESTER3_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 2/5: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$VOTES" != 1 ] && echoErr "ERROR TEST 2/5: Expected votes is '1', but got '$VOTES'" && exit 1
  # <--- Decline transaction

  # TEST 2/5/1
  # Repeat decline transaction --->
  echoInfo "TEST 2/5/1"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED)) # same
  TESTER3_BALANCE_EXPECTED=$(($TESTER3_BALANCE_EXPECTED - 150)) # -150 fee
  TESTER5_BALANCE_EXPECTED=$TESTER5_BALANCE_EXPECTED # same

  declineTransaction tester3 tester1 "$HASH" 150 ukex

  VOTES=$(getCustodyPoolVotes tester1 "$HASH")
  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER3_BALANCE_REAL=$(showBalance tester3 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 2/5/1: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER3_BALANCE_EXPECTED" != "$TESTER3_BALANCE_REAL" ] && echoErr "ERROR TEST 2/5/1: Expected tester3 account balance to be '$TESTER3_BALANCE_EXPECTED', but got '$TESTER3_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 2/5/1: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$VOTES" != 1 ] && echoErr "ERROR TEST 2/5/1: Expected votes is '1', but got '$VOTES'" && exit 1
  # <--- Repeat decline transaction

  # TEST 2/6
  # Last approve transaction --->
  echoInfo "TEST 2/6"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 7 - 200)) # -7 sent -200 reward
  TESTER4_BALANCE_EXPECTED=$((1000000000000 - 150 + 200)) # -150 fee +200 reward
  TESTER5_BALANCE_EXPECTED=$(($TESTER5_BALANCE_EXPECTED + 7)) # +7 received

  approveTransaction tester4 tester1 "$HASH" 150 ukex

  POOL=$(getCustodyPool tester1)
  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER4_BALANCE_REAL=$(showBalance tester4 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 2/6: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER4_BALANCE_EXPECTED" != "$TESTER4_BALANCE_REAL" ] && echoErr "ERROR TEST 2/6: Expected tester4 account balance to be '$TESTER4_BALANCE_EXPECTED', but got '$TESTER4_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 2/6: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$POOL" != "null" ] && [ "$POOL" != "{}" ] && echoErr "ERROR TEST 2/6: Expected empty pool" && exit 1
  # <--- Last approve transaction

# <--- Send tokens with enabled custodians and disabled password protection

# TEST 3
# Send tokens with disabled custodians and enabled password protection --->

  # TEST 3/1
  # Disable custodians and enable password protection --->
  echoInfo "TEST 3/1"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 150)) # -150 fee

  enableCustody tester1 0 1 0 0 150 ukex

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 3/1: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  # <--- Disable custodians and enable password protection

  # TEST 3/2
  # Send tokens --->
  echoInfo "TEST 3/2"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 150)) # -150 fee
  TESTER5_BALANCE_EXPECTED=$TESTER5_BALANCE_EXPECTED # same

  HASH=$(custodySendTokens tester1 tester5 7 ukex 150 ukex 1000 ukex $PASSWORD)

  POOL=$(getCustodyPool tester1)
  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 3/2: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 3/2: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$POOL" == "null" ] || [ "$POOL" == "{}" ] && echoErr "ERROR TEST 3/2: Expected not empty pool" && exit 1
  # <--- Send tokens

  # TEST 3/3
  # Password confirm transaction --->
  echoInfo "TEST 3/3"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 7)) # -7 sent
  TESTER5_BALANCE_EXPECTED=$(($TESTER5_BALANCE_EXPECTED - 150 + 7)) # -150 fee + 7 received

  passwordConfirmTransaction tester5 tester1 "$HASH" "$PASSWORD" 150 ukex

  POOL=$(getCustodyPool tester1)
  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 3/3: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 3/3: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$POOL" != "null" ] && [ "$POOL" != "{}" ] && echoErr "ERROR TEST 3/3: Expected empty pool" && exit 1
  # <--- Password confirm transaction

# <--- Send tokens with disabled custodians and enabled password protection

# TEST 4
# Send tokens with enabled custodians and enabled password protection --->

  # TEST 4/1
  # Drop custodians --->
  echoInfo "TEST 4/1"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 150)) # -150 fee

  dropCustody tester1 150 ukex

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 4/1: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  # <--- Disable custodians

  # TEST 4/2
  # Send tokens --->
  echoInfo "TEST 4/2"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 150 - 7)) # -150 fee -7 sent
  TESTER5_BALANCE_EXPECTED=$(($TESTER5_BALANCE_EXPECTED + 7)) # +7 sent

  HASH=$(custodySendTokens tester1 tester5 7 ukex 150 ukex 600 ukex $PASSWORD)

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
  POOL=$(getCustodyPool tester1)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 4/2: Expected tester5 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 4/2: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$POOL" != "null" ] && [ "$POOL" != "{}" ] && echoErr "ERROR TEST 4/2: Expected empty pool" && exit 1
  # <--- Send tokens

  # TEST 4/3
  # Enable custodians and password protection --->
  echoInfo "TEST 4/3"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 150)) # -150 fee

  enableCustody tester1 30 1 0 0 150 ukex 0

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 4/3: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  # <--- Enable custodians and password protection

  # TEST 4/3
  # Add custodians --->
  echoInfo "TEST 4/4"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 150)) # -150 fee

  addCustodians tester1 "$INPUT_CUSTODIANS" 150 ukex

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 4/4: Expected tester5 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  # <--- Add custodians

  # TEST 4/5
  # Send tokens --->
  echoInfo "TEST 4/5"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 150)) # -150 fee
  TESTER5_BALANCE_EXPECTED=$(($TESTER5_BALANCE_EXPECTED)) # same

  HASH=$(custodySendTokens tester1 tester5 7 ukex 150 ukex 600 ukex $PASSWORD)

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
  POOL=$(getCustodyPool tester1)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 4/5: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 4/5: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$POOL" == "null" ] || [ "$POOL" == "{}" ] && echoErr "ERROR TEST 4/5: Expected not empty pool" && exit 1
  # <--- Send tokens

  # TEST 4/6
  # Approve transaction --->
  echoInfo "TEST 4/6"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 200)) # -200 reward
  TESTER4_BALANCE_EXPECTED=$(($TESTER4_BALANCE_EXPECTED - 150 + 200)) # -150 fee + 200 reward
  TESTER5_BALANCE_EXPECTED=$(($TESTER5_BALANCE_EXPECTED)) # same

  approveTransaction tester4 tester1 "$HASH" 150 ukex

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER4_BALANCE_REAL=$(showBalance tester4 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
  POOL=$(getCustodyPool tester1)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 4/6: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER4_BALANCE_EXPECTED" != "$TESTER4_BALANCE_REAL" ] && echoErr "ERROR TEST 4/6: Expected tester4 account balance to be '$TESTER4_BALANCE_EXPECTED', but got '$TESTER4_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 4/6: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$POOL" == "null" ] || [ "$POOL" == "{}" ] && echoErr "ERROR TEST 4/6: Expected not empty pool" && exit 1
  # <--- Approve transaction

  # TEST 4/7
  # Password confirm transaction --->
  echoInfo "TEST 4/7"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 7)) # -7 sent
  TESTER5_BALANCE_EXPECTED=$(($TESTER5_BALANCE_EXPECTED - 150 + 7)) # -150 fee +7 sent

  passwordConfirmTransaction tester5 tester1 "$HASH" "$PASSWORD" 150 ukex

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
  POOL=$(getCustodyPool tester1)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 4/7: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 4/7: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$POOL" != "null" ] && [ "$POOL" != "{}" ] && echoErr "ERROR TEST 4/7: Expected empty pool" && exit 1
  # <--- Password confirm transaction

  # TEST 4/8
  # Send tokens --->
  echoInfo "TEST 4/8"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 150)) # -150 fee
  TESTER5_BALANCE_EXPECTED=$(($TESTER5_BALANCE_EXPECTED)) # same

  HASH=$(custodySendTokens tester1 tester5 7 ukex 150 ukex 600 ukex $PASSWORD)

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
  POOL=$(getCustodyPool tester1)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 4/8: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 4/8: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$POOL" == "null" ] || [ "$POOL" == "{}" ] && echoErr "ERROR TEST 4/8: Expected not empty pool 6" && exit 1
  # <--- Send tokens

  # TEST 4/9
  # Password confirm transaction --->
  echoInfo "TEST 4/9"
  TESTER5_BALANCE_EXPECTED=$(($TESTER5_BALANCE_EXPECTED - 150)) # same

  passwordConfirmTransaction tester5 tester1 "$HASH" "$PASSWORD" 150 ukex

  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
  POOL=$(getCustodyPool tester1)

  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 4/9: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$POOL" == "null" ] || [ "$POOL" == "{}" ] && echoErr "ERROR TEST 4/9: Expected not empty pool 7" && exit 1
  # <--- Password confirm transactio

  # TEST 4/10
  # Last approve transaction --->
  echoInfo "TEST 4/10"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 7 - 200)) # -7 sent -200 reward
  TESTER4_BALANCE_EXPECTED=$(($TESTER4_BALANCE_EXPECTED - 150 + 200)) # 150 fee +200 reward
  TESTER5_BALANCE_EXPECTED=$(($TESTER5_BALANCE_EXPECTED + 7)) # +7 received

  approveTransaction tester4 tester1 "$HASH" 150 ukex

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER4_BALANCE_REAL=$(showBalance tester4 ukex)
  TESTER5_BALANCE_REAL=$(showBalance tester5 ukex)
  POOL=$(getCustodyPool tester1)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 4/10: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER4_BALANCE_EXPECTED" != "$TESTER4_BALANCE_REAL" ] && echoErr "ERROR TEST 4/10: Expected tester4 account balance to be '$TESTER4_BALANCE_EXPECTED', but got '$TESTER4_BALANCE_REAL'" && exit 1
  [ "$TESTER5_BALANCE_EXPECTED" != "$TESTER5_BALANCE_REAL" ] && echoErr "ERROR TEST 4/10: Expected tester5 account balance to be '$TESTER5_BALANCE_EXPECTED', but got '$TESTER5_BALANCE_REAL'" && exit 1
  [ "$POOL" != "null" ] && [ "$POOL" != "{}" ] && echoErr "ERROR TEST 4/10: Expected empty pool 8" && exit 1
  # <--- Last approve transaction

echoInfo "INFO: $TEST_NAME - Integration Test - END, elapsed: $(prettyTime "$(timerSpan $TEST_NAME)")"

# TEST 5
# CHeck next control and target options --->

  # TEST 5/1
  # Enable custodians and set next controll address tester2 --->
  echoInfo "TEST 5/1"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 150)) # -150 fee

  enableCustody tester1 60 0 0 0 150 ukex tester2

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)

  # todo: get next controller and test it

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 5/1: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  # <--- Enable custodians

  # TEST 5/2
  # Add custodians test double SHA protection and next controller with target
  echoInfo "TEST 5/2"
  TESTER2_BALANCE_EXPECTED=$(($TESTER2_BALANCE_EXPECTED - 150)) # -150 fee

  addCustodians tester2 "$INPUT_CUSTODIANS" 150 ukex tester3 tester1

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER2_BALANCE_REAL=$(showBalance tester2 ukex)


  # todo: get next controller and test it

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 5/2: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER2_BALANCE_EXPECTED" != "$TESTER2_BALANCE_REAL" ] && echoErr "ERROR TEST 5/2: Expected tester2 account balance to be '$TESTER2_BALANCE_EXPECTED', but got '$TESTER2_BALANCE_REAL'" && exit 1

  # <--- Add custodians test double SHA protection

echoInfo "INFO: $TEST_NAME - Integration Test - END, elapsed: $(prettyTime $(timerSpan $TEST_NAME))"