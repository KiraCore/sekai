#!/usr/bin/env bash
# To run test locally: make network-start && ./scripts/test-local/token-transfers.sh
set -e
set -x
. /etc/profile
. ./scripts/sekai-env.sh

TEST_NAME="POLL" && timerStart $TEST_NAME
echoInfo "INFO: $TEST_NAME - Integration Test - START"

addAccount polltester1
addAccount polltester2
addAccount polltester3
addAccount polltester4
addAccount polltester5

ACCOUNT1_ADDRESS=$(showAddress polltester1)
ACCOUNT2_ADDRESS=$(showAddress polltester2)
ACCOUNT3_ADDRESS=$(showAddress polltester3)
ACCOUNT4_ADDRESS=$(showAddress polltester4)
ACCOUNT5_ADDRESS=$(showAddress polltester5)

sendTokens validator "$ACCOUNT1_ADDRESS" 1000000000000 ukex 100 ukex
sendTokens validator "$ACCOUNT2_ADDRESS" 1000000000000 ukex 100 ukex
sendTokens validator "$ACCOUNT3_ADDRESS" 1000000000000 ukex 100 ukex
sendTokens validator "$ACCOUNT4_ADDRESS" 1000000000000 ukex 100 ukex
sendTokens validator "$ACCOUNT5_ADDRESS" 1000000000000 ukex 100 ukex

createRole validator poll_voter description
assignRole validator poll_voter polltester1
assignRole validator poll_voter polltester2
assignRole validator poll_voter polltester3
assignRole validator poll_voter polltester4
assignRole validator poll_voter polltester5

whitelistPermission validator 66 polltester1
whitelistPermission validator 66 polltester2
whitelistPermission validator 66 polltester3
whitelistPermission validator 66 polltester4
whitelistPermission validator 66 polltester5

#createPoll 1
createPoll polltester1 poll_voter yes,no,mbe string 1 3 1m

#checkPollResult 1
POLL_RESULT=$(getPollResult "$ACCOUNT1_ADDRESS" 0)

[ "$POLL_RESULT" != "POLL_PENDING" ] && echoErr "ERROR: Expected result 0 to be POLL_PENDING, but got '$POLL_RESULT'" && exit 1
# ------------
sleep 60s
POLL_RESULT=$(getPollResult "$ACCOUNT1_ADDRESS" 0)

[ "$POLL_RESULT" != "POLL_RESULT_QUORUM_NOT_REACHED" ] && echoErr "ERROR: Expected result 1 should be POLL_RESULT_QUORUM_NOT_REACHED, but got '$POLL_RESULT'" && exit 1
# ------------

#createPoll 2
createPoll polltester1 poll_voter yes,no,mbe string 1 3 1m

addVote polltester2 2 2 yes
addVote polltester3 2 2 yes
addVote polltester4 2 2 yes
addVote polltester5 2 2 yes

sleep 60s
POLL_RESULT=$(getPollResult "$ACCOUNT1_ADDRESS" 1)

[ "$POLL_RESULT" != "POLL_RESULT_PASSED" ] && echoErr "ERROR: Expected result 2 to be POLL_RESULT_PASSED, but got '$POLL_RESULT'" && exit 1
# ------------

createPoll polltester1 poll_voter yes,no,mbe string 1 3 1m
addVote polltester2 3 3
addVote polltester3 3 3
addVote polltester4 3 3
addVote polltester5 3 2 mbe

sleep 60s
POLL_RESULT=$(getPollResult "$ACCOUNT1_ADDRESS" 2)

[ "$POLL_RESULT" != "POLL_RESULT_REJECTED_WITH_VETO" ] && echoErr "ERROR: Expected result 3 to be POLL_RESULT_REJECTED_WITH_VETO, but got '$POLL_RESULT'" && exit 1
# ------------

createPoll polltester1 poll_voter yes,no,mbe string 1 3 1m
addVote polltester2 4 1
addVote polltester3 4 1
addVote polltester4 4 3
addVote polltester5 4 2 mbe

sleep 60s
POLL_RESULT=$(getPollResult "$ACCOUNT1_ADDRESS" 3)

[ "$POLL_RESULT" != "POLL_RESULT_REJECTED" ] && echoErr "ERROR: Expected result 4 to be POLL_RESULT_REJECTED, but got '$POLL_RESULT'" && exit 1
# ------------

createPoll polltester1 poll_voter 1,2,3 int 1 4 1m
addVote polltester2 5 2 4
addVote polltester3 5 2 4
addVote polltester4 5 2 4

sleep 60s
POLL_RESULT=$(getPollResult "$ACCOUNT1_ADDRESS" 4)

[ "$POLL_RESULT" != "POLL_RESULT_PASSED" ] && echoErr "ERROR: Expected result 6 to be POLL_RESULT_PASSED, but got '$POLL_RESULT'" && exit 1
# ------------

echoInfo "INFO: $TEST_NAME - Integration Test - END, elapsed: $(prettyTime "$(timerSpan $TEST_NAME)")"