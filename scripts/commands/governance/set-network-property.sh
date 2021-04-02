#!/bin/bash

sekaid tx customgov proposal set-network-property MIN_TX_FEE 101 --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes