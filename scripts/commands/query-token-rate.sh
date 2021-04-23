#!/bin/bash

# command
sekaid query tokens rate mykex
# response
# denom: mykex
# fee_payments: true
# rate: "1.500000"

# command
sekaid query tokens rate invalid_denom
# response
# Error: invalid_denom denom does not exist

# command
sekaid query tokens all-rates --chain-id=testing --home=$HOME/.sekaid
# response
# data:
# - denom: ubtc
#   fee_payments: true
#   rate: "0.000010"
# - denom: ukex
#   fee_payments: true
#   rate: "1.000000"
# - denom: xeth
#   fee_payments: true
#   rate: "0.000100"

# command
sekaid query tokens rates-by-denom ukex --chain-id=testing --home=$HOME/.sekaid
# response
# data:
#   ukex:
#     denom: ukex
#     fee_payments: true
#     rate: "1.000000"
