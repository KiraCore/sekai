package cosmos

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/gateway/interx"
	"github.com/KiraCore/sekai/INTERX/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	abciTypes "github.com/tendermint/tendermint/abci/types"
)

// RegisterCosmosBlockRoutes registers query routers.
func RegisterCosmosBlockRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryBlocks, QueryBlocksRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryBlockByHeightOrHash, QueryBlockByHeightOrHashRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryBlockTransactions, QueryBlockTransactionsRequest(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryBlocks, "This is an API to query blocks by parameters.", true)
	common.AddRPCMethod("GET", config.QueryBlockByHeightOrHash, "This is an API to query block by height", true)
	common.AddRPCMethod("GET", config.QueryBlockTransactions, "This is an API to query block transactions by height", true)
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

	return common.MakeGetRequest(rpcAddr, "/blockchain", strings.Join(events, "&"))
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
	success, err, statusCode := common.MakeGetRequest(rpcAddr, "/block", fmt.Sprintf("height=%s", height))

	if err != nil {
		success, err, statusCode = common.MakeGetRequest(rpcAddr, "/block_by_hash", fmt.Sprintf("hash=%s", height))
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
				common.GetLogger().Info("[query-blocks-by-height] Seach from the cache: ", height)
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

func getTransactions(attributes []sdk.Attribute) []types.Transaction {
	txs := []types.Transaction{}
	var evMap = make(map[string]string)
	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		value := string(attribute.GetValue())

		if _, ok := evMap[key]; ok {
			coin, err := sdk.ParseCoinNormalized(evMap["amount"])
			if err == nil {
				tx := types.Transaction{}
				tx.From = evMap["sender"]
				tx.To = evMap["recipient"]
				tx.Amount = coin.Amount.Int64()
				tx.Denom = coin.GetDenom()
				txs = append(txs, tx)
			}

			evMap = make(map[string]string)
		}

		evMap[key] = value
	}

	if _, ok := evMap["amount"]; ok {
		coin, err := sdk.ParseCoinNormalized(evMap["amount"])
		if err == nil {
			tx := types.Transaction{}
			tx.From = evMap["sender"]
			tx.To = evMap["recipient"]
			tx.Amount = coin.Amount.Int64()
			tx.Denom = coin.GetDenom()
			txs = append(txs, tx)
		}

		evMap = make(map[string]string)
	}

	return txs
}

func getTransactionsFromTm(attributes []abciTypes.EventAttribute) []types.Transaction {
	txs := []types.Transaction{}
	var evMap = make(map[string]string)
	for _, attribute := range attributes {
		key := string(attribute.GetKey())
		value := string(attribute.GetValue())

		if _, ok := evMap[key]; ok {
			coin, err := sdk.ParseCoinNormalized(evMap["amount"])
			if err == nil {
				tx := types.Transaction{}
				tx.From = evMap["sender"]
				tx.To = evMap["recipient"]
				tx.Amount = coin.Amount.Int64()
				tx.Denom = coin.GetDenom()
				txs = append(txs, tx)
			}

			evMap = make(map[string]string)
		}

		evMap[key] = value
	}

	if _, ok := evMap["amount"]; ok {
		coin, err := sdk.ParseCoinNormalized(evMap["amount"])
		if err == nil {
			tx := types.Transaction{}
			tx.From = evMap["sender"]
			tx.To = evMap["recipient"]
			tx.Amount = coin.Amount.Int64()
			tx.Denom = coin.GetDenom()
			txs = append(txs, tx)
		}

		evMap = make(map[string]string)
	}

	return txs
}

// QueryBlockTransactionsHandle is a function to query transactions of a block.
func QueryBlockTransactionsHandle(rpcAddr string, height string) (interface{}, interface{}, int) {
	blockHeight, _ := strconv.Atoi(height)
	response, err := interx.SearchTxHashHandle(rpcAddr, "", "", "", 0, int64(blockHeight), int64(blockHeight))
	if err != nil {
		return common.ServeError(0, "transaction query failed", "", http.StatusBadRequest)
	}

	searchResult := types.TransactionSearchResult{}

	searchResult.TotalCount = response.TotalCount
	searchResult.Txs = []types.TransactionResult{}

	for _, transaction := range response.Txs {
		txResult := types.TransactionResult{}

		txResult.Hash = hex.EncodeToString(transaction.Hash)
		txResult.Status = "Success"
		if transaction.TxResult.Code != 0 {
			txResult.Status = "Failure"
		}

		txResult.BlockHeight = transaction.Height
		txResult.BlockTimestamp, err = common.GetBlockTime(rpcAddr, transaction.Height)
		if err != nil {
			common.GetLogger().Error("[query-transactions] Block not found: ", transaction.Height)
			return common.ServeError(0, "", fmt.Sprintf("block not found: %d", transaction.Height), http.StatusInternalServerError)
		}
		txResult.Confirmation = common.NodeStatus.Block - transaction.Height + 1
		txResult.GasWanted = transaction.TxResult.GetGasWanted()
		txResult.GasUsed = transaction.TxResult.GetGasUsed()

		txResult.Transactions = []types.Transaction{}

		logs, err := sdk.ParseABCILogs(transaction.TxResult.GetLog())
		if err != nil {
			common.GetLogger().Error("[query-transactions] Failed to parse ABCI logs: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		transferTxs := []types.Transaction{}
		for _, log := range logs {
			for _, event := range log.GetEvents() {
				if event.GetType() == "transfer" {
					transferTxs = append(transferTxs, getTransactions(event.GetAttributes())...)
				}
			}
		}
		common.GetLogger().Info(transferTxs)
		for _, event := range transaction.TxResult.Events {
			if event.GetType() == "transfer" {
				txs := getTransactionsFromTm(event.GetAttributes())
				common.GetLogger().Info(txs)
				for _, tx := range txs {
					index := 0
					for index < len(transferTxs) {
						if transferTxs[index].From == tx.From &&
							transferTxs[index].To == tx.To &&
							transferTxs[index].Amount == tx.Amount &&
							transferTxs[index].Denom == tx.Denom {
							break
						}
						index++
					}

					if index < len(transferTxs) {
						common.GetLogger().Info("transfer found")
						txResult.Transactions = append(txResult.Transactions, tx)
						transferTxs = append(transferTxs[:index], transferTxs[index+1:]...)
					} else {
						common.GetLogger().Info("fee found")
						txResult.Fees = append(txResult.Transactions, tx)
					}
				}
			}
		}

		common.GetLogger().Info("txResult", txResult)
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

		common.GetLogger().Info("[query-block-transactions-by-height] Entering Block query by height: %s", height)

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
