package interx

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
	kiratypes "github.com/KiraCore/sekai/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/rpc/core/types"
)

// RegisterBlockRoutes registers block/transaction query routers.
func RegisterBlockRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryBlocks, QueryBlocksRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryBlockByHeightOrHash, QueryBlockByHeightOrHashRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryBlockTransactions, QueryBlockTransactionsRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryTransactionResult, QueryTransactionResultRequest(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryBlocks, "This is an API to query blocks by parameters.", true)
	common.AddRPCMethod("GET", config.QueryBlockByHeightOrHash, "This is an API to query block by height", true)
	common.AddRPCMethod("GET", config.QueryBlockTransactions, "This is an API to query block transactions by height", true)
	common.AddRPCMethod("GET", config.QueryTransactionResult, "This is an API to query transaction result by hash", true)
}

func queryBlocksHandle(rpcAddr string, r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	_ = r.ParseForm()

	minHeight := r.FormValue("minHeight")
	maxHeight := r.FormValue("maxHeight")

	var events = make([]string, 0, 2)

	if minHeight != "" {
		events = append(events, fmt.Sprintf("minHeight=%s", minHeight))
	}

	if maxHeight != "" {
		events = append(events, fmt.Sprintf("maxHeight=%s", maxHeight))
	}

	// search blocks

	return common.MakeTendermintRPCRequest(rpcAddr, "/blockchain", strings.Join(events, "&"))
}

// QueryBlocksRequest is a function to query Blocks.
func QueryBlocksRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-blocks] Entering Blocks query")

		if !common.RPCMethods["GET"][config.QueryBlocks].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryBlocks].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-blocks] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryBlocksHandle(rpcAddr, r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryBlocks].CachingEnabled)
	}
}

func queryBlockByHeightOrHashHandle(rpcAddr string, height string) (interface{}, interface{}, int) {
	success, err, statusCode := common.MakeTendermintRPCRequest(rpcAddr, "/block", fmt.Sprintf("height=%s", height))

	if err != nil {
		success, err, statusCode = common.MakeTendermintRPCRequest(rpcAddr, "/block_by_hash", fmt.Sprintf("hash=%s", height))
	}

	return success, err, statusCode
}

// QueryBlockByHeightOrHashRequest is a function to query Blocks.
func QueryBlockByHeightOrHashRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		height := queries["height"]
		request := types.InterxRequest{
			Method:   r.Method,
			Endpoint: config.QueryBlockByHeightOrHash,
			Params:   []byte(height),
		}
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-blocks-by-height] Entering Block query by height: ", height)

		if !common.RPCMethods["GET"][config.QueryBlockByHeightOrHash].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryBlockByHeightOrHash].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-blocks-by-height] Returning from the cache: ", height)
					return
				}
			}

			response.Response, response.Error, statusCode = queryBlockByHeightOrHashHandle(rpcAddr, height)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryBlockByHeightOrHash].CachingEnabled)
	}
}

func getTransactionsFromLog(attributes []abciTypes.EventAttribute) []sdk.Coin {
	feeTxs := []sdk.Coin{}

	var evMap = make(map[string]string)
	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		value := string(attribute.GetValue())

		if _, ok := evMap[key]; ok {
			coin, err := sdk.ParseCoinNormalized(evMap["amount"])
			if err == nil {
				feeTx := sdk.Coin{
					Amount: sdk.NewInt(coin.Amount.Int64()),
					Denom:  coin.GetDenom(),
				}
				feeTxs = append(feeTxs, feeTx)
			}

			evMap = make(map[string]string)
		}

		evMap[key] = value
	}

	if _, ok := evMap["amount"]; ok {
		coin, err := sdk.ParseCoinNormalized(evMap["amount"])
		if err == nil {
			feeTx := sdk.Coin{
				Amount: sdk.NewInt(coin.Amount.Int64()),
				Denom:  coin.GetDenom(),
			}
			feeTxs = append(feeTxs, feeTx)
		}
	}

	return feeTxs
}

func parseTransaction(rpcAddr string, transaction tmTypes.ResultTx) (types.TransactionResult, error) {
	txResult := types.TransactionResult{}

	tx, err := config.EncodingCg.TxConfig.TxDecoder()(transaction.Tx)
	if err != nil {
		common.GetLogger().Error("[query-transactions] Failed to decode transaction: ", err)
		return txResult, err
	}

	txResult.Hash = hex.EncodeToString(transaction.Hash)
	txResult.Status = "Success"
	if transaction.TxResult.Code != 0 {
		txResult.Status = "Failure"
	}

	txResult.BlockHeight = transaction.Height
	txResult.BlockTimestamp, err = common.GetBlockTime(rpcAddr, transaction.Height)
	if err != nil {
		common.GetLogger().Error("[query-transactions] Block not found: ", transaction.Height)
		return txResult, fmt.Errorf("block not found: %d", transaction.Height)
	}
	txResult.Confirmation = common.NodeStatus.Block - transaction.Height + 1
	txResult.GasWanted = transaction.TxResult.GetGasWanted()
	txResult.GasUsed = transaction.TxResult.GetGasUsed()

	txSigning, ok := tx.(signing.Tx)
	if ok {
		txResult.Memo = txSigning.GetMemo()
	}

	txResult.Msgs = make([]types.TxMsg, 0)
	for _, msg := range tx.GetMsgs() {
		txResult.Msgs = append(txResult.Msgs, types.TxMsg{
			Type: kiratypes.MsgType(msg),
			Data: msg,
		})
	}

	txResult.Transactions = []types.Transaction{}
	txResult.Fees = []sdk.Coin{}

	logs, err := sdk.ParseABCILogs(transaction.TxResult.GetLog())
	if err != nil {
		return txResult, nil
	}

	for _, event := range transaction.TxResult.Events {
		if event.GetType() == "transfer" {
			txResult.Fees = append(txResult.Fees, getTransactionsFromLog(event.GetAttributes())...)
		}
	}

	for index, msg := range tx.GetMsgs() {
		log := logs[index]
		txType := kiratypes.MsgType(msg)
		transfers := []types.Transaction{}

		var evMap = make(map[string]([]sdk.Attribute))
		for _, event := range log.GetEvents() {
			evMap[event.GetType()] = event.GetAttributes()
		}

		if txType == "send" {
			msgSend := msg.(*bank.MsgSend)

			amounts := []sdk.Coin{}
			for _, coin := range msgSend.Amount {
				amounts = append(amounts, sdk.Coin{
					Denom:  coin.GetDenom(),
					Amount: sdk.NewInt(coin.Amount.Int64()),
				})
			}

			transfers = append(transfers, types.Transaction{
				Type:    txType,
				From:    msgSend.FromAddress,
				To:      msgSend.ToAddress,
				Amounts: amounts,
			})
		} else if txType == "multisend" {
			msgMultiSend := msg.(*bank.MsgMultiSend)
			inputs := msgMultiSend.GetInputs()
			outputs := msgMultiSend.GetOutputs()
			if len(inputs) == 1 && len(outputs) == 1 {
				input := inputs[0]
				output := outputs[0]
				amounts := []sdk.Coin{}

				for _, coin := range input.Coins {
					amounts = append(amounts, sdk.Coin{
						Denom:  coin.GetDenom(),
						Amount: sdk.NewInt(coin.Amount.Int64()),
					})
				}

				transfers = append(transfers, types.Transaction{
					Type:    txType,
					From:    input.Address,
					To:      output.Address,
					Amounts: amounts,
				})
			}
		} else if txType == "create_validator" {
			createValidatorMsg := msg.(*staking.MsgCreateValidator)

			transfers = append(transfers, types.Transaction{
				Type: txType,
				From: createValidatorMsg.DelegatorAddress,
				To:   createValidatorMsg.ValidatorAddress,
				Amounts: []sdk.Coin{
					{
						Denom:  createValidatorMsg.Value.Denom,
						Amount: sdk.NewInt(createValidatorMsg.Value.Amount.Int64()),
					},
				},
			})
		} else if txType == "delegate" {
			delegateMsg := msg.(*staking.MsgDelegate)

			transfers = append(transfers, types.Transaction{
				Type: txType,
				From: delegateMsg.DelegatorAddress,
				To:   delegateMsg.ValidatorAddress,
				Amounts: []sdk.Coin{
					{
						Denom:  delegateMsg.Amount.Denom,
						Amount: sdk.NewInt(delegateMsg.Amount.Amount.Int64()),
					},
				},
			})
		} else if txType == "begin_redelegate" {
			reDelegateMsg := msg.(*staking.MsgBeginRedelegate)

			transfers = append(transfers, types.Transaction{
				Type: txType,
				From: reDelegateMsg.ValidatorSrcAddress,
				To:   reDelegateMsg.ValidatorDstAddress,
				Amounts: []sdk.Coin{
					{
						Denom:  reDelegateMsg.Amount.Denom,
						Amount: sdk.NewInt(reDelegateMsg.Amount.Amount.Int64()),
					},
				},
			})
		} else if txType == "begin_unbonding" {
			unDelegateMsg := msg.(*staking.MsgUndelegate)

			transfers = append(transfers, types.Transaction{
				Type: txType,
				From: unDelegateMsg.ValidatorAddress,
				To:   unDelegateMsg.DelegatorAddress,
				Amounts: []sdk.Coin{
					{
						Denom:  unDelegateMsg.Amount.Denom,
						Amount: sdk.NewInt(unDelegateMsg.Amount.Amount.Int64()),
					},
				},
			})
		} else if txType == "withdraw_delegator_reward" {
			var coin sdk.Coin
			if v, found := evMap["withdraw_rewards"]; found && len(v) >= 2 {
				if v[0].GetKey() == "amount" {
					coin, _ = sdk.ParseCoinNormalized(v[0].Value)
				} else if v[1].GetKey() == "amount" {
					coin, _ = sdk.ParseCoinNormalized(v[1].Value)
				}
			}

			withdrawDelegatorRewardMsg := msg.(*distribution.MsgWithdrawDelegatorReward)

			transfers = append(transfers, types.Transaction{
				Type: txType,
				From: withdrawDelegatorRewardMsg.ValidatorAddress,
				To:   withdrawDelegatorRewardMsg.DelegatorAddress,
				Amounts: []sdk.Coin{
					{
						Denom:  coin.Denom,
						Amount: sdk.NewInt(coin.Amount.Int64()),
					},
				},
			})
		} else {
			attributes := []sdk.Attribute{}
			for _, event := range log.GetEvents() {
				if event.GetType() == "transfer" {
					attributes = event.GetAttributes()
				}
			}

			txs := []types.Transaction{}

			var evMap = make(map[string]string)
			for _, attribute := range attributes {
				key := string(attribute.GetKey())
				value := string(attribute.GetValue())

				if _, ok := evMap[key]; ok {
					coin, err := sdk.ParseCoinNormalized(evMap["amount"])
					if err == nil {
						txs = append(txs, types.Transaction{
							Type: txType,
							From: evMap["sender"],
							To:   evMap["recipient"],
							Amounts: []sdk.Coin{
								{
									Denom:  coin.Denom,
									Amount: sdk.NewInt(coin.Amount.Int64()),
								},
							},
						})
					}

					evMap = make(map[string]string)
				}

				evMap[key] = value
			}

			if _, ok := evMap["amount"]; ok {
				coin, err := sdk.ParseCoinNormalized(evMap["amount"])
				if err == nil {
					txs = append(txs, types.Transaction{
						Type: txType,
						From: evMap["sender"],
						To:   evMap["recipient"],
						Amounts: []sdk.Coin{
							{
								Denom:  coin.Denom,
								Amount: sdk.NewInt(coin.Amount.Int64()),
							},
						},
					})
				}
			}

			transfers = append(transfers, txs...)
		}

		for _, transfer := range transfers {
			for _, amount := range transfer.Amounts {
				i := 0
				for i < len(txResult.Fees) {
					if txResult.Fees[i].Amount.Equal(amount.Amount) && txResult.Fees[i].Denom == amount.Denom {
						break
					}
					i++
				}

				if i < len(txResult.Fees) {
					txResult.Fees = append(txResult.Fees[:i], txResult.Fees[i+1:]...)
					break
				}
			}
		}

		txResult.Transactions = append(txResult.Transactions, transfers...)
	}

	return txResult, nil
}

// QueryBlockTransactionsHandle is a function to query transactions of a block.
func QueryBlockTransactionsHandle(rpcAddr string, height string) (interface{}, interface{}, int) {
	blockHeight, _ := strconv.Atoi(height)
	response, err := SearchTxHashHandle(rpcAddr, "", "", "", 0, 0, int64(blockHeight), int64(blockHeight), "")
	if err != nil {
		return common.ServeError(0, "transaction query failed", "", http.StatusBadRequest)
	}

	searchResult := types.TransactionSearchResult{}

	searchResult.TotalCount = response.TotalCount
	searchResult.Txs = []types.TransactionResult{}

	for _, transaction := range response.Txs {
		txResult, err := parseTransaction(rpcAddr, *transaction)
		if err != nil {
			return common.ServeError(0, "", err.Error(), http.StatusBadRequest)
		}

		searchResult.Txs = append(searchResult.Txs, txResult)
	}

	return searchResult, nil, http.StatusOK
}

// QueryBlockTransactionsRequest is a function to query transactions of a block.
func QueryBlockTransactionsRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		height := queries["height"]
		request := types.InterxRequest{
			Method:   r.Method,
			Endpoint: config.QueryBlockTransactions,
			Params:   []byte(height),
		}
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-block-transactions-by-height] Entering Block query by height: ", height)

		if !common.RPCMethods["GET"][config.QueryBlockTransactions].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryBlockTransactions].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-block-transactions-by-height] Returning from the cache: %s", height)
					return
				}
			}

			response.Response, response.Error, statusCode = QueryBlockTransactionsHandle(rpcAddr, height)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryBlockTransactions].CachingEnabled)
	}
}

// QueryTransactionResultHandle is a function to query transaction by a given hash.
func QueryTransactionResultHandle(rpcAddr string, txHash string) (interface{}, interface{}, int) {
	response, err := SearchTxHashHandle(rpcAddr, "", "", "", 0, 0, 0, 0, strings.TrimPrefix(txHash, "0x"))
	if err != nil {
		return common.ServeError(0, "transaction query failed", "", http.StatusBadRequest)
	}

	txResult := types.TransactionResult{}

	for _, transaction := range response.Txs {
		txResult, err = parseTransaction(rpcAddr, *transaction)
		if err != nil {
			return common.ServeError(0, "", err.Error(), http.StatusBadRequest)
		}
	}

	return txResult, nil, http.StatusOK
}

// QueryTransactionResultRequest is a function to query transactions by a given hash.
func QueryTransactionResultRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		txHash := queries["txHash"]
		request := types.InterxRequest{
			Method:   r.Method,
			Endpoint: config.QueryTransactionResult,
			Params:   []byte(txHash),
		}
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-transaction-by-hash] Entering transaction query by hash: %s", txHash)

		if !common.RPCMethods["GET"][config.QueryTransactionResult].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryTransactionResult].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-transaction-by-hash] Returning from the cache: %s", txHash)
					return
				}
			}

			response.Response, response.Error, statusCode = QueryTransactionResultHandle(rpcAddr, txHash)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryTransactionResult].CachingEnabled)
	}
}
