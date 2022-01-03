#!/bin/bash

# command
sekaid query tokens alias KEX
# response
# allowed_vote_types:
# - "yes"
# - "no"
# decimals: 6
# denoms:
# - ukex
# enactment: 0
# expiration: 0
# icon: myiconurl
# name: Kira
# status: undefined
# symbol: KEX


# command
sekaid query tokens alias KE
# response
# Error: KE symbol does not exist

# command
sekaid query tokens all-aliases --chain-id=testing --home=$HOME/.sekaid
# response
# data:
# - allowed_vote_types:
#   - "yes"
#   - "no"
#   decimals: 6
#   denoms:
#   - ukex
#   enactment: 0
#   expiration: 0
#   icon: myiconurl
#   name: Kira
#   status: undefined
#   symbol: KEX

# command
sekaid query tokens aliases-by-denom ukex --chain-id=testing --home=$HOME/.sekaid
# response
# data:
#   ukex:
#     allowed_vote_types:
#     - "yes"
#     - "no"
#     decimals: 6
#     denoms:
#     - ukex
#     enactment: 0
#     expiration: 0
#     icon: myiconurl
#     name: Kira
#     status: undefined
#     symbol: KEX