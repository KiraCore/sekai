#!/bin/bash

# query validator account by address
sekaid query validator --addr $(sekaid keys show -a validator --keyring-backend=test --home=$HOME/.sekaid)