#!/usr/bin/env bash
set -e
set -x
. ./scripts/sekai-env.sh
. ./scripts/sekai-utils.sh

TEST_NAME="NETWORK-START"
echo "INFO: $TEST_NAME - Integration Test - START"
echo "INFO: Ensuring essential dependencies are installed & up to date"

UTILS_VER=$(bash-utils bashUtilsVersion 2> /dev/null || echo "")
[[ $(bash-utils versionToNumber "$UTILS_VER" 2> /dev/null || echo "0") -ge $(bash-utils versionToNumber "v0.2.13" 2> /dev/null || echo "1") ]] && \
 UTILS_OLD_VER="false" || UTILS_OLD_VER="true" 

# Installing utils is essential to simplify the setup steps
if [ "$UTILS_OLD_VER" != "false" ] ; then
    echo "INFO: KIRA utils were NOT installed on the system, setting up..." && sleep 2
    TOOLS_VERSION="v0.2.13" && mkdir -p /usr/keys && FILE_NAME="bash-utils.sh" && \
     if [ -z "$KIRA_COSIGN_PUB" ] ; then KIRA_COSIGN_PUB=/usr/keys/kira-cosign.pub ; fi && \
     echo -e "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE/IrzBQYeMwvKa44/DF/HB7XDpnE+\nf+mU9F/Qbfq25bBWV2+NlYMJv3KvKHNtu3Jknt6yizZjUV4b8WGfKBzFYw==\n-----END PUBLIC KEY-----" > $KIRA_COSIGN_PUB && \
     wget "https://github.com/KiraCore/tools/releases/download/$TOOLS_VERSION/${FILE_NAME}" -O ./$FILE_NAME && \
     wget "https://github.com/KiraCore/tools/releases/download/$TOOLS_VERSION/${FILE_NAME}.sig" -O ./${FILE_NAME}.sig && \
     cosign verify-blob --key="$KIRA_COSIGN_PUB" --signature=./${FILE_NAME}.sig ./$FILE_NAME && \
     chmod -v 555 ./$FILE_NAME && ./$FILE_NAME bashUtilsSetup "/var/kiraglob"
fi

bash-utils loadGlobEnvs
. /etc/profile
echoInfo "INFO: Bash-utils version: $(bashUtilsVersion)"

SYSCTRL_DESTINATION=/usr/local/bin/systemctl2
if [ ! -f $SYSCTRL_DESTINATION ] ; then
    safeWget $SYSCTRL_DESTINATION \
     https://raw.githubusercontent.com/gdraheim/docker-systemctl-replacement/9cbe1a00eb4bdac6ff05b96ca34ec9ed3d8fc06c/files/docker/systemctl.py \
     "e02e90c6de6cd68062dadcc6a20078c34b19582be0baf93ffa7d41f5ef0a1fdd" && \
    chmod +x $SYSCTRL_DESTINATION
fi

echoInfo "INFO: Systemctrl2 version: $(systemctl2 --version)"

timerStart $TEST_NAME
./scripts/sekai-utils.sh sekaiUtilsSetup

echoInfo "INFO: Environment cleanup...."
NETWORK_NAME="localnet-1"
setGlobEnv SEKAID_HOME ~/.sekaid-$NETWORK_NAME
setGlobEnv NETWORK_NAME $NETWORK_NAME
loadGlobEnvs
. /etc/profile

rm -rfv $SEKAID_HOME 
mkdir -p $SEKAID_HOME

echoInfo "INFO: Starting new network..."
sekaid init --overwrite --chain-id=$NETWORK_NAME "KIRA TEST LOCAL VALIDATOR NODE" --home=$SEKAID_HOME
sekaid keys add validator --keyring-backend=test --home=$SEKAID_HOME
sekaid keys add faucet --keyring-backend=test --home=$SEKAID_HOME
sekaid add-genesis-account $(showAddress validator) 150000000000000ukex,300000000000000test,2000000000000000000000000000samolean,1000000lol --keyring-backend=test --home=$SEKAID_HOME
sekaid add-genesis-account $(showAddress faucet) 150000000000000ukex,300000000000000test,2000000000000000000000000000samolean,1000000lol --keyring-backend=test --home=$SEKAID_HOME
sekaid gentx-claim validator --keyring-backend=test --moniker="GENESIS VALIDATOR" --home=$SEKAID_HOME

CFG="$SEKAID_HOME/config/config.toml"
# set block time to 0.5 second
setTomlVar "[consensus]" timeout_commit "500ms" $CFG
# progress ASAP we have all precommits needed
setTomlVar "[consensus]" skip_timeout_commit "true" $CFG

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

systemctl2 status sekai

echoInfo "INFO: Checking network status..."
NETWORK_STATUS_CHAIN_ID=$(showStatus | jq .NodeInfo.network | xargs)

if [ "$NETWORK_NAME" != "$NETWORK_STATUS_CHAIN_ID" ] ; then
    echoErr "ERROR: Incorrect chain ID from the status query, expected '$NETWORK_NAME', but got $NETWORK_STATUS_CHAIN_ID"
fi

echoInfo "INFO: Waiting for next block to be produced..."
timeout 60 sekai-utils awaitBlocks 2
BLOCK_HEIGHT=$(showBlockHeight)
timeout 60 sekai-utils awaitBlocks 2
NEXT_BLOCK_HEIGHT=$(showBlockHeight)

if [ $BLOCK_HEIGHT -ge $NEXT_BLOCK_HEIGHT ] ; then
    echoErr "ERROR: Failed to produce next block height, stuck at $BLOCK_HEIGHT"
fi

echoInfo "INFO: Printing sekai status..."
showStatus | jq

set +x
echoInfo "INFO: SEKAI $(sekaid version) is running"
echoInfo "INFO: NETWORK-START - Integration Test - END, elapsed: $(prettyTime $(timerSpan $TEST_NAME))"