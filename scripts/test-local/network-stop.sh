#!/usr/bin/env bash
set -e
set +x
. /etc/profile

TEST_NAME="NETWORK-STOP" && timerStart $TEST_NAME
echoInfo "INFO: $TEST_NAME - Integration Test - START"
set -x

systemctl2 stop sekai || echoWarn "WARNING: sekai service was NOT running or could NOT be stopped"

pkill -15 sekaid || echoWarn "WARNING: Failed to kill sekaid process"

kill -9 $(sudo lsof -t -i:9090) || echoWarn "WARNING: Nothing running on port 9090, or failed to kill processes"
kill -9 $(sudo lsof -t -i:6060) || echoWarn "WARNING: Nothing running on port 6060, or failed to kill processes"
kill -9 $(sudo lsof -t -i:26656) || echoWarn "WARNING: Nothing running on port 26656, or failed to kill processes"
kill -9 $(sudo lsof -t -i:26657) || echoWarn "WARNING: Nothing running on port 26657, or failed to kill processes"
kill -9 $(sudo lsof -t -i:26658) || echoWarn "WARNING: Nothing running on port 26658, or failed to kill processes"

STOP_SUCCESS="false"
sekaid status || STOP_SUCCESS="true"

[ "${STOP_SUCCESS,,}" != "true" ] && echoErr "ERROR: Failed to stop sekaid process!" && exit 1

set +x
echoInfo "INFO: NETWORK-STOP - Integration Test - END, elapsed: $(prettyTime $(timerSpan $TEST_NAME))"