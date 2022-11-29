#!/usr/bin/env bash
set -e
set -x
CURRENT_DIR=$(pwd)

UTILS_VER=$(bash-utils bashUtilsVersion 2> /dev/null || echo "")
[[ $(bash-utils versionToNumber "$UTILS_VER" || echo "0") -ge $(bash-utils versionToNumber "v0.2.13" || echo "1") ]] && \
 UTILS_OLD_VER="false" || UTILS_OLD_VER="true" 

# Installing utils is essential to simplify the setup steps
if [ "$UTILS_OLD_VER" == "true" ] ; then
    echo "INFO: KIRA utils were NOT installed on the system, setting up..." && sleep 2
    TOOLS_VERSION="v0.2.13" && mkdir -p /usr/keys && FILE_NAME="bash-utils.sh" && \
     if [ -z "$KIRA_COSIGN_PUB" ] ; then KIRA_COSIGN_PUB=/usr/keys/kira-cosign.pub ; fi && \
     echo -e "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE/IrzBQYeMwvKa44/DF/HB7XDpnE+\nf+mU9F/Qbfq25bBWV2+NlYMJv3KvKHNtu3Jknt6yizZjUV4b8WGfKBzFYw==\n-----END PUBLIC KEY-----" > $KIRA_COSIGN_PUB && \
     wget "https://github.com/KiraCore/tools/releases/download/$TOOLS_VERSION/${FILE_NAME}" -O ./$FILE_NAME && \
     wget "https://github.com/KiraCore/tools/releases/download/$TOOLS_VERSION/${FILE_NAME}.sig" -O ./${FILE_NAME}.sig && \
     cosign verify-blob --key="$KIRA_COSIGN_PUB" --signature=./${FILE_NAME}.sig ./$FILE_NAME && \
     chmod -v 555 ./$FILE_NAME && ./$FILE_NAME bashUtilsSetup "/var/kiraglob" && . /etc/profile
     echoInfo "Installed bash-utils $(bashUtilsVersion)"
else
    . /etc/profile
    echoInfo "INFO: KIRA utils are up to date, latest version $UTILS_VER" && sleep 2
fi

GO_VER=$(go version 2> /dev/null || echo "")

# install golang if needed
if  ($(isNullOrEmpty "$GO_VER")) || ($(isNullOrEmpty "$GOBIN")) ; then
     GO_TAR="go1.18.3.$(getPlatform)-$(getArch).tar.gz" && rm -rfv /usr/local/go && cd /tmp && rm -fv ./$GO_TAR && \
     wget https://dl.google.com/go/${GO_TAR} && \
     tar -C /usr/local -xvf $GO_TAR && rm -fv ./$GO_TAR && \
     setGlobEnv GOROOT "/usr/local/go" && setGlobPath "\$GOROOT" && \
     setGlobEnv GOBIN "/usr/local/go/bin" && setGlobPath "\$GOBIN" && \
     setGlobEnv GOPATH "$HOME/go" && setGlobPath "\$GOPATH" && \
     setGlobEnv GOCACHE "$HOME/go/cache" && \
     loadGlobEnvs && \
     mkdir -p "$GOPATH/src" "$GOPATH/bin" "$GOCACHE" && \
     chmod -R 777 "$GOPATH" && chmod -R 777 "$GOROOT" && chmod -R 777 "$GOCACHE"
    echoInfo "INFO: Sucessfully intalled go"
fi

go version

# navigate to current direcotry and load global environment variables
cd $CURRENT_DIR

go clean -modcache
EXPECTED_PROTO_DEP_VER="v0.0.2"
BUF_VER=$(buf --version 2> /dev/null || echo "")

if ($(isNullOrEmpty "$BUF_VER")) || [ "$SEKAI_PROTO_DEP_VER" != "$EXPECTED_PROTO_DEP_VER" ] ; then
    GO111MODULE=on 
    go install github.com/bufbuild/buf/cmd/buf@v1.0.0-rc10
    echoInfo "INFO: Sucessfully intalled buf $(buf --version)"

    setGlobEnv GOLANG_PROTOBUF_VERSION "1.27.1" && \
     setGlobEnv GOGO_PROTOBUF_VERSION "1.3.2" && \
     setGlobEnv GRPC_GATEWAY_VERSION "1.14.7" && \
     loadGlobEnvs

    go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest && \
     go install google.golang.org/protobuf/cmd/protoc-gen-go@v${GOLANG_PROTOBUF_VERSION} && \
     go install github.com/gogo/protobuf/protoc-gen-gogo@v${GOGO_PROTOBUF_VERSION} && \
     go install github.com/gogo/protobuf/protoc-gen-gogofast@v${GOGO_PROTOBUF_VERSION} && \
     go install github.com/gogo/protobuf/protoc-gen-gogofaster@v${GOGO_PROTOBUF_VERSION} && \
     go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v${GRPC_GATEWAY_VERSION} && \
     go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v${GRPC_GATEWAY_VERSION} && \
     go install github.com/gogo/protobuf/protoc-gen-gogotypes

    # Following command executes with error requiring us to silence it, however the executable is placed in $GOBIN
    # https://github.com/regen-network/cosmos-proto
    # Original setup originates from Docker Image tendermintdev/sdk-proto-gen:v0.2
    # reference: 
    go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos
    go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

    setGlobEnv SEKAI_PROTO_DEP_VER "$EXPECTED_PROTO_DEP_VER"
fi

CONSTANS_FILE=./types/constants.go
COSMOS_BRANCH=$(grep -Fn -m 1 'CosmosVersion ' $CONSTANS_FILE | rev | cut -d "=" -f1 | rev | xargs | tr -dc '[:alnum:]\-\.' || echo '')
($(isNullOrEmpty "$COSMOS_BRANCH")) && ( echoErr "ERROR: CosmosVersion was NOT found in contants '$CONSTANS_FILE' !" && sleep 5 && exit 1 )

go get github.com/cosmos/cosmos-sdk@$COSMOS_BRANCH

echoInfo "Cleaning up proto gen files..."
rm -rfv ./github.com
cosmos_sdk_dir=$(go list -f '{{ .Dir }}' -m github.com/cosmos/cosmos-sdk@$COSMOS_BRANCH)
kira_dir=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)

echoInfo "Generating protobuf files..."
for dir in $kira_dir; do
  # generate protobuf bind
  buf protoc \
  -I "proto" \
  -I "$cosmos_sdk_dir/third_party/proto" \
  -I "$cosmos_sdk_dir/proto" \
  --gocosmos_out=plugins=interfacetype+grpc,\
Mgoogle/protobuf/any.proto=github.com/cosmos/cosmos-sdk/codec/types:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')

  # generate grpc gateway
  buf protoc \
  -I "proto" \
  -I "$cosmos_sdk_dir/third_party/proto" \
  -I "$cosmos_sdk_dir/proto" \
  --grpc-gateway_out=logtostderr=true:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')
done

cp -r github.com/KiraCore/sekai/* ./
rm -rf github.com

echoInfo "INFO: Success, all proto files were compiled!"
