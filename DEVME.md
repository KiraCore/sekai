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

# install deb package manager
echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | tee /etc/apt/sources.list.d/goreleaser.list && apt-get update -y && \
	apt install nfpm

FILE_NAME="bash-utils.sh" && TOOLS_VERSION="v0.1.5" && \
 wget "https://github.com/KiraCore/tools/releases/download/$TOOLS_VERSION/${FILE_NAME}" -O ./$FILE_NAME && \
 chmod -v 755 ./$FILE_NAME && ./$FILE_NAME bashUtilsSetup "$GLOBAL_COMMON" && source $FILE_NAME

# install go
setGlobEnv GOROOT /usr/local/go && setGlobEnv GOPATH /home/go && setGlobEnv GOBIN /usr/local/go/bin && \
 loadGlobEnvs && setGlobPath $GOROOT && setGlobPath $GOPATH && setGlobPath $GOBIN && loadGlobEnvs && \
 ( go clean -modcache -cache -n || : ) && rm -rfv "$GOROOT" "$GOBIN" "$GOPATH" && GO_VERSION="1.18.2" && \
 GO_TAR="go$GO_VERSION.$(getPlatform)-$(getArch).tar.gz" && cd /tmp && safeWget ./$GO_TAR https://dl.google.com/go/$GO_TAR \
 "fc4ad28d0501eaa9c9d6190de3888c9d44d8b5fb02183ce4ae93713f67b8a35b,e54bec97a1a5d230fc2f9ad0880fcbabb5888f30ed9666eca4a91c5a32e86cbc" && \
 tar -C /usr/local -xf $GO_TAR &>/dev/null && go version

# mount C drive or other disk where repo is stored
setGlobLine "mount -t drvfs C:" "mount -t drvfs C: /mnt/c || echo 'Failed to mount C drive'"

# set env variable to your local repos (will vary depending on the user)
setGlobEnv SEKAI_REPO "/mnt/c/Users/asmodat/Desktop/KIRA/KIRA-CORE/GITHUB/sekai" && \
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
cd $SEKAI_REPO && \
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

# Start local test network

```
make network-start
loadGlobEnvs
```

# Stop local test network

```
make network-stop
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