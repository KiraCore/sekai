#!/bin/bash

# query validator account by address
sekaid query customstaking validator --addr $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)

# query all validators with specific filter
sekaid query customstaking validators --moniker="" --addr="" --val-addr="" --moniker="" --status="" --pubkey="" --proposer=""
