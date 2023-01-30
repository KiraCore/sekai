#!/bin/bash

# command to generate recovery secret
sekaid tx recovery generate-recovery-secret 10a0fbe01030000122300000000000
NONCE=00
PROOF=29eeb09666c4792c314e631063932621f573cbb07af7274657d1314e1892eb93
CHALLENGE=220c7e47a53ef4c2161035308d4fdc52f619e54f8a657b208ba3708139fdc03d

# commmand to create an account
sekaid keys add recovery --keyring-backend=test
RECOVERY=$(sekaid keys show -a recovery --keyring-backend=test)

# commands to query
sekaid query recovery recovery-record $(sekaid keys show -a validator --keyring-backend=test)

# command to register recovery secret
sekaid tx recovery register-recovery-secret $CHALLENGE $NONCE $PROOF --from=validator --keyring-backend=test --chain-id=testing --fees=1000ukex -y --home=$HOME/.sekaid --broadcast-mode=block

# command to rotate
sekaid tx recovery rotate-recovery-address $RECOVERY $PROOF --from=validator --chain-id=testing --keyring-backend=test --fees=1000ukex -y --home=$HOME/.sekaid --yes --broadcast-mode=block

# command to check validators after rotation
sekaid query customstaking validators --moniker="" --addr="" --val-addr="" --moniker="" --status="" --pubkey="" --proposer=""
