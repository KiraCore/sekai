#!/bin/bash

# To get sekaid version, should install binary with `make install`
make install

sekaid version
# 1.0.0

sekaid version --long
# name: sekai
# server_name: sekaid
# version: 1.0.0
# commit: cad573dbd60f477799e241589410a155f988d120
# build_tags: ""
# go: go version go1.16 darwin/amd64
# build_deps:
# - github.com/99designs/keyring@v1.1.6
