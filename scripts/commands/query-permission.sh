#!/bin/bash

# command
sekaid query customgov permissions $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

# response
# blacklist: []
# whitelist:
# - 4
# - 3