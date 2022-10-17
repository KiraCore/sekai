#!/bin/bash

# e.g. sekaid query customgov execution-fee <msg_type>

# command
sekaid query customgov execution-fee "B"
# response
# fee:
#   default_parameters: "0"
#   execution_fee: "10"
#   failure_fee: "1"
#   name: ABC
#   timeout: "10"
#   transaction_type: B

# genesis fee configuration test
sekaid query customgov execution-fee "A"
# response
# fee:
#   default_parameters: "0"
#   execution_fee: "10"
#   failure_fee: "1"
#   name: Claim Validator Seat
#   timeout: "10"
#   transaction_type: A

# query all execution fees
sekaid query customgov all-execution-fees
# response
# fees:
# - default_parameters: "0"
#   execution_fee: "100"
#   failure_fee: "1000"
#   timeout: "10"
#   transaction_type: activate
# - default_parameters: "0"
#   execution_fee: "100"
#   failure_fee: "1"
#   timeout: "10"
#   transaction_type: claim-councilor
# - default_parameters: "0"
#   execution_fee: "100"
#   failure_fee: "1"
#   timeout: "10"
#   transaction_type: claim-proposal-type-x
# - default_parameters: "0"
#   execution_fee: "100"
#   failure_fee: "1"
#   timeout: "10"
#   transaction_type: claim-validator
# - default_parameters: "0"
#   execution_fee: "100"
#   failure_fee: "100"
#   timeout: "10"
#   transaction_type: pause
# - default_parameters: "0"
#   execution_fee: "10"
#   failure_fee: "1"
#   timeout: "10"
#   transaction_type: submit-proposal-type-x
# - default_parameters: "0"
#   execution_fee: "100"
#   failure_fee: "100"
#   timeout: "10"
#   transaction_type: unpause
# - default_parameters: "0"
#   execution_fee: "100"
#   failure_fee: "1"
#   timeout: "10"
#   transaction_type: upsert-token-alias
# - default_parameters: "0"
#   execution_fee: "100"
#   failure_fee: "1"
#   timeout: "10"
#   transaction_type: veto-proposal-type-x
# - default_parameters: "0"
#   execution_fee: "100"
#   failure_fee: "1"
#   timeout: "10"
#   transaction_type: vote-proposal-type-x