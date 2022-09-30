#!/bin/bash

sekaid tx customslashing pause --from=validator --keyring-backend=test --home=$HOME/.sekaid --chain-id=testing --fees=100ukex --yes --broadcast-mode=block 
