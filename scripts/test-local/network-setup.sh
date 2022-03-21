#!/usr/bin/env bash
set -e
set -x
. /etc/profile

timerStart NETWORK_SETUP
echoInfo "INFO: NETWORK-SETUP - Integration Test - START"

echoInfo "INFO: Ensuring essential dependencies are installed & up to date"
SYSCTRL_DESTINATION=/usr/local/bin/systemctl2
if [ ! -f $SYSCTRL_DESTINATION ] ; then
    safeWget /usr/local/bin/systemctl2 \
     https://raw.githubusercontent.com/gdraheim/docker-systemctl-replacement/9cbe1a00eb4bdac6ff05b96ca34ec9ed3d8fc06c/files/docker/systemctl.py \
     "e02e90c6de6cd68062dadcc6a20078c34b19582be0baf93ffa7d41f5ef0a1fdd" && \
    chmod +x $SYSCTRL_DESTINATION && \
    systemctl2 --version
fi

echoInfo "INFO: Environment cleanup...."
systemctl2 stop sekai || echoWarn "WARNING: sekai service was NOT running or could NOT be stopped"

setGlobEnv SEKAID_HOME ~/.sekaid-local
loadGlobEnvs

rm -rfv $SEKAID_HOME 
mkdir -p $SEKAID_HOME

echoInfo "INFO: Starting new network..."
CHAIN_ID="localnet-0"
sekaid init --overwrite --chain-id=$CHAIN_ID "KIRA TEST LOCAL VALIDATOR NODE" --home=$SEKAID_HOME
sekaid keys add validator --keyring-backend=test --home=$SEKAID_HOME
sekaid add-genesis-account $(showAddress validator) 300000000000000ukex,300000000000000test,2000000000000000000000000000samolean,1000000lol --home=$SEKAID_HOME
sekaid gentx-claim validator --keyring-backend=test --moniker="GENESIS VALIDATOR" --home=$SEKAID_HOME

cat > /etc/systemd/system/sekai.service << EOL
[Unit]
Description=Local KIRA Test Network
After=network.target
[Service]
MemorySwapMax=0
Type=simple
User=root
WorkingDirectory=/root
ExecStart=$GOBIN/sekaid start --home=$SEKAID_HOME --trace
Restart=always
RestartSec=5
LimitNOFILE=4096
[Install]
WantedBy=default.target
EOL

systemctl2 enable sekai 
systemctl2 start sekai

echoInfo "INFO: Waiting for network to start..." && sleep 3

echoInfo "INFO: Checking network status..."
NETWORK_STATUS_CHAIN_ID=$(showStatus | jq .NodeInfo.network | xargs)

if [ "$CHAIN_ID" != "$NETWORK_STATUS_CHAIN_ID" ] ; then
    echoErr "ERROR: Incorrect chain ID from the status query, expected '$CHAIN_ID', but got $NETWORK_STATUS_CHAIN_ID"
fi

BLOCK_HEIGHT=$(showBlockHeight)
echoInfo "INFO: Waiting for next block to be produced..." && sleep 10
NEXT_BLOCK_HEIGHT=$(showBlockHeight)

if [ $BLOCK_HEIGHT -ge $NEXT_BLOCK_HEIGHT ] ; then
    echoErr "ERROR: Failed to produce next block height, stuck at $BLOCK_HEIGHT"
fi

sekaid version
sleep 10

echoInfo "INFO: NETWORK-SETUP - Integration Test - END, elapsed: $(prettyTime $(timerSpan NETWORK_SETUP))"