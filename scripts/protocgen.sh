#!/usr/bin/env bash

protoc -I "./x/staking/types/" -I "third_party/proto" --gocosmos_out=plugins=interfacetype+grpc,\
Mgoogle/protobuf/any.proto=github.com/KiraCore/cosmos-sdk/codec/types:. \
./x/staking/types/msg.proto

# move proto files to the right places
cp -r github.com/KiraCore/cosmos-sdk/* ./
rm -rf github.com