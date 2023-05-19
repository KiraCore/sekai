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
KEY7="7902699be42c8a8e46fbbb4501726517e86b22c56a189f7625a6da49081b2451"
KEY71="1048f03db5d45f654b955eae20d84b72673680fb13b318e7da22e8dce58df21c"

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

  enableCustody tester1 60 0 0 0 150 ukex 0 $KEY11

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  CUSTODY_KEY=$(getCustodyKey tester1)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 2/1: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$CUSTODY_KEY" != $KEY11 ] && echoErr "ERROR TEST 2/1: Expected key to be $KEY11, but got '$CUSTODY_KEY'" && exit 1
  # <--- Enable custodians

  # TEST 2/2
  # Add custodians test double SHA protection --->
  echoInfo "TEST 2/2"
  addCustodiansForce tester1 "$INPUT_CUSTODIANS" 1 ukex $KEY1 $KEY21
  addCustodiansForce tester1 "$INPUT_CUSTODIANS" 2 ukex $KEY1 $KEY31

  sleep 10

  TESTER1_BALANCE_EXPECTED=$(showBalance tester1 ukex)
  CUSTODY_KEY=$(getCustodyKey tester1)

  [ "$CUSTODY_KEY" != $KEY21 ] && echoErr "ERROR TEST 2/2: Expected key1 to be $KEY21, but got '$CUSTODY_KEY'" && exit 1
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

  enableCustody tester1 0 1 0 0 150 ukex $KEY2 $KEY31

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  CUSTODY_KEY=$(getCustodyKey tester1)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 3/1: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$CUSTODY_KEY" != $KEY31 ] && echoErr "ERROR TEST 3/1: Expected key to be $KEY31, but got '$CUSTODY_KEY'" && exit 1
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
  # Disable custodians --->
  echoInfo "TEST 4/1"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 150)) # -150 fee

  disableCustody tester1 150 ukex $KEY3

  CUSTODY_KEY=$(getCustodyKey tester1)
  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 4/1: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$CUSTODY_KEY" != '' ] && echoErr "ERROR TEST 4/1: Expected key to be "", but got '$CUSTODY_KEY'" && exit 1
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

  enableCustody tester1 30 1 0 0 150 ukex 0 $KEY41

  CUSTODY_KEY=$(getCustodyKey tester1)
  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 4/3: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$CUSTODY_KEY" != $KEY41 ] && echoErr "ERROR TEST 4/3: Expected key to be $KEY41, but got '$CUSTODY_KEY'" && exit 1
  # <--- Enable custodians and password protection

  # TEST 4/3
  # Add custodians --->
  echoInfo "TEST 4/4"
  TESTER1_BALANCE_EXPECTED=$(($TESTER1_BALANCE_EXPECTED - 150)) # -150 fee

  addCustodians tester1 "$INPUT_CUSTODIANS" 150 ukex $KEY4 $KEY51

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  CUSTODY_KEY=$(getCustodyKey tester1)

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 4/4: Expected tester5 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$CUSTODY_KEY" != $KEY51 ] && echoErr "ERROR TEST 4/4: Expected key to be $KEY51, but got '$CUSTODY_KEY'" && exit 1
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

  enableCustody tester1 60 0 0 0 150 ukex $KEY5 $KEY61 tester2

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  CUSTODY_KEY1=$(getCustodyKey tester1)

  # todo: get next controller and test it

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 5/1: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$CUSTODY_KEY1" != $KEY61 ] && echoErr "ERROR TEST 5/1: Expected key to be $KEY61, but got '$CUSTODY_KEY1'" && exit 1
  # <--- Enable custodians

  # TEST 5/2
  # Add custodians test double SHA protection and next controller with target
  echoInfo "TEST 5/2"
  TESTER2_BALANCE_EXPECTED=$(($TESTER2_BALANCE_EXPECTED - 150)) # -150 fee

  addCustodians tester2 "$INPUT_CUSTODIANS" 150 ukex $KEY6 $KEY71 tester3 tester1

  TESTER1_BALANCE_REAL=$(showBalance tester1 ukex)
  TESTER2_BALANCE_REAL=$(showBalance tester2 ukex)
  CUSTODY_KEY=$(getCustodyKey tester1)

  # todo: get next controller and test it

  [ "$TESTER1_BALANCE_EXPECTED" != "$TESTER1_BALANCE_REAL" ] && echoErr "ERROR TEST 5/2: Expected tester1 account balance to be '$TESTER1_BALANCE_EXPECTED', but got '$TESTER1_BALANCE_REAL'" && exit 1
  [ "$TESTER2_BALANCE_EXPECTED" != "$TESTER2_BALANCE_REAL" ] && echoErr "ERROR TEST 5/2: Expected tester2 account balance to be '$TESTER2_BALANCE_EXPECTED', but got '$TESTER2_BALANCE_REAL'" && exit 1
  [ "$CUSTODY_KEY" != $KEY71 ] && echoErr "ERROR TEST 5/2: Expected key1 to be $KEY71, but got '$CUSTODY_KEY'" && exit 1
  # <--- Add custodians test double SHA protection

echoInfo "INFO: $TEST_NAME - Integration Test - END, elapsed: $(prettyTime $(timerSpan $TEST_NAME))"