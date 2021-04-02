#!/bin/bash

# query all signing infos as an array
sekaid query customslashing signing-infos
# response
# info:
# - address: kiravalcons166p6nw8gm4cq7xescmyzm8qsqf56za0x5q6ep9
#   inactive_until: "1970-01-01T00:00:00Z"
#   index_offset: "2"
#   missed_blocks_counter: "0"
#   start_height: "0"
#   tombstoned: false
# pagination:
#   next_key: null
#   total: "0"

# query signing info by validator
sekaid query customslashing signing-info $(sekaid tendermint show-validator)