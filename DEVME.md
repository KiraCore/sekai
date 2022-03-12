# Dev on Windows 10

```
# Open Ubuntu 20.04 WSL 2.0 console

sudo -s

# Install Essential Dependencies

apt-get install -y curl && curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && apt-get update -y && \
 apt-get install -y --allow-unauthenticated --allow-downgrades --allow-remove-essential --allow-change-held-packages \
 software-properties-common wget git nginx apt-transport-https file build-essential net-tools hashdeep \
 protobuf-compiler golang-goprotobuf-dev golang-grpc-gateway golang-github-grpc-ecosystem-grpc-gateway-dev lsb-release \
 clang cmake gcc g++ pkg-config libudev-dev libusb-1.0-0-dev iputils-ping nano jq python python3 python3-pip gnupg \
 bash libglu1-mesa lsof bc dnsutils psmisc netcat  make nodejs tar unzip xz-utils yarn zip p7zip-full ca-certificates \
 containerd docker.io dos2unix

# install systemd alternative
wget https://raw.githubusercontent.com/gdraheim/docker-systemctl-replacement/master/files/docker/systemctl.py -O /usr/local/bin/systemctl2 && \
 chmod +x /usr/local/bin/systemctl2 && \
 systemctl2 --version

# install kira bash helper utils
BRANCH="v0.0.2" && cd /tmp && rm -fv ./i.sh && \
wget https://raw.githubusercontent.com/KiraCore/tools/$BRANCH/bash-utils/install.sh -O ./i.sh && \
 chmod 555 -v ./i.sh && ./i.sh "$BRANCH" "/var/kiraglob" && . /etc/profile && rm -fv ./i.sh

# uninstall golang if needed
( go clean -modcache -cache -n || echo "Failed to cleanup go cache" ) && \
( rm -rfv "$GOROOT" || echo "Failed to cleanup go root" ) && \
( rm -rfv "$GOBIN" || echo "Failed to cleanup go bin" ) && \
( rm -rfv "$GOPATH" || echo "Failed to cleanup go path" ) && \
( rm -rfv "$GOCACHE" || echo "Failed to cleanup go cache" )

# mount C drive or other disk where repo is stored
setGlobLine "mount -t drvfs C:" "mount -t drvfs C: /mnt/c || echo 'Failed to mount C drive'"

# set env variable to your local repos (will vary depending on the user)
setGlobEnv SEKAI_REPO "/mnt/c/Users/asmodat/Desktop/KIRA/KIRA-CORE/GITHUB/sekai" && \
 setGlobEnv INTERX_REPO "/mnt/c/Users/asmodat/Desktop/KIRA/KIRA-CORE/GITHUB/interx" && \
 loadGlobEnvs

# set home directory of your repos
setGlobEnv SEKAID_HOME "/root/.sekaid" && \
 setGlobEnv INTERXD_HOME "/root/.interxd" && \
 loadGlobEnvs

# Ensure you have Docker Desktop installed: https://code.visualstudio.com/blogs/2020/03/02/docker-in-wsl2 & reboot your entire host machine
```

# Clean Clone
```
cd $HOME && rm -fvr ./sekai && SEKAI_BRANCH="master" && \
 git clone https://github.com/KiraCore/sekai.git -b $SEKAI_BRANCH && \
 cd ./sekai
```

## Installation
```
cd $INTERX_REPO

chmod -Rv 777 ./scripts && \
 dos2unix ./scripts/protocgen-local.sh

make install
```

# If any changes are made to modules update protobuf files

```
# NOTE: This is SUPER important that `protocgen.sh` does NOT contain `\r\n` line endings! Ensure proper Unix `LF` EOL Conversion in notepad++ 
cd "$SEKAI_REPO" && \
 chmod +x ./scripts/protocgen.sh && \
 make proto-gen && \
 go mod tidy && \
 make install && echo "SEKAI update success" || echo "SEKAI update failed" 
```

# Updating INTERX GO Modules

Pull dependencies & ensures that the go.mod file matches the source code in the module. This must be done every time Cosmos-SDK version is updated in the `go.mod` file of the root repository. Always ensure that `go.mod` in the root & INTERX folder match `github.com/cosmos/cosmos-sdk` version!

```
cd $SEKAI_REPO/INTERX && \
 go get -u && \
 go mod tidy && echo "INTERX update success" || echo "INTERX update failed" 
```

# Build repo & create service 

```
make install
sekaid version
systemctl2 stop sekai || echo "sekai service was NOT running or could NOT be stopped"

rm -rfv $SEKAID_HOME && mkdir -p $SEKAID_HOME && \
 sekaid init --overwrite --chain-id="localnet-1" "LOCAL KIRA VALIDATOR NODE" --home=$SEKAID_HOME && \
 sekaid keys add validator --keyring-backend=test --home=$SEKAID_HOME && \
 sekaid add-genesis-account $(sekaid keys show validator -a --keyring-backend=test --home=$SEKAID_HOME) 300000000000000ukex,300000000000000test,2000000000000000000000000000samolean,1000000lol --home=$SEKAID_HOME && \
 sekaid gentx-claim validator --keyring-backend=test --moniker="GENESIS VALIDATOR" --home=$SEKAID_HOME && \
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

# Load / Update CLI Helper Scripts
SCRIPTS_BRANCH="testnet" && mkdir -p "/common/kiraglob" && mkdir -p "/common/scripts" &&  \
 wget "https://raw.githubusercontent.com/KiraCore/kira/$SCRIPTS_BRANCH/docker/base-image/scripts/utils.sh" -O "/common/scripts/utils.sh" && \
 wget "https://raw.githubusercontent.com/KiraCore/kira/$SCRIPTS_BRANCH/docker/kira/container/sekaid-helper.sh" -O "/common/scripts/sekaid-helper.sh" && \
 . ~/.bash_aliases

# Auto Execute Helper Scripts (run only once)
echo "source '/common/scripts/utils.sh' || echo 'ERROR: Failed to load utils script'" >> ~/.bash_aliases && \
 echo "source '/common/scripts/sekaid-helper.sh'  || echo 'ERROR: Failed to load sekaid-helper script'" >> ~/.bash_aliases && \
 echo "NETWORK_NAME='localnet-1'" >> ~/.bash_aliases

# Start sekai service
systemctl2 daemon-reload && \
 systemctl2 enable sekai && \
 systemctl2 restart sekai && \
 systemctl2 status sekai


# Stop sekai service
systemctl2 stop sekai || echoWarn "WARNING: Failed to stop KIRA Plan!"

```

# Running Tests
```
gov_test
network_test
network
app
ante_test
evidence_test
keeper_test
simulation_test
types_test
slashing_test
cli_test
staking_test
tokens_test
upgrade_test
```