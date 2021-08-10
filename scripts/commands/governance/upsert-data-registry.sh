#!/bin/bash

sekaid tx customgov proposal upsert-data-registry 111 0x0000 "https://github.githubassets.com/images/modules/notifications/inbox-zero.svg" "image/svg+xml" 31088 --keyring-backend=test --from=validator --home=$HOME/.sekaid --title="upsert-data-registry" --description="upsert data" --chain-id=testing --yes --fees=100ukex
sekaid tx customgov proposal vote 1 1 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --yes --fees=100ukex
sekaid query customgov proposals
sekaid query customgov all-data-reference-keys
