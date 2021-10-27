package interx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmJsonRPCTypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

// RegisterInterxTxRoutes registers tx query routers.
func RegisterInterxTxRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryWithdraws, QueryWithdraws(rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryDeposits, QueryDeposits(rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryUnconfirmedTxs, QueryUnconfirmedTxs(rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryKiraFunctions, "This is an API to query kira functions and metadata.", true)
	common.AddRPCMethod("GET", config.QueryWithdraws, "This is an API to query withdraw transactions.", true)
	common.AddRPCMethod("GET", config.QueryDeposits, "This is an API to query deposit transactions.", true)
	common.AddRPCMethod("GET", config.QueryUnconfirmedTxs, "This is an API to query unconfirmed transactions.", true)
}

func toSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// SearchTxHashHandle is a function to query transactions
func SearchTxHashHandle(rpcAddr string, sender string, recipient string, txType string, page int, limit int, txMinHeight int64, txMaxHeight int64, txHash string) (*tmTypes.ResultTxSearch, error) {
	var events = make([]string, 0, 5)

	if sender != "" {
		events = append(events, fmt.Sprintf("transfer.sender='%s'", sender))
	}

	if recipient != "" {
		events = append(events, fmt.Sprintf("transfer.recipient='%s'", recipient))
	}

	if txType != "all" && txType != "" {
		events = append(events, fmt.Sprintf("message.action='%s'", txType))
	}

	if txHash != "" {
		events = append(events, fmt.Sprintf("tx.hash='%s'", txHash))
	}

	if txMinHeight >= 0 {
		events = append(events, fmt.Sprintf("tx.height>=%d", txMinHeight))
	}

	if txMaxHeight >= 0 {
		events = append(events, fmt.Sprintf("tx.height<=%d", txMaxHeight))
	}

	// search transactions
	endpoint := fmt.Sprintf("%s/tx_search?query=\"%s\"&page=%d&&per_page=%d&order_by=\"desc\"", rpcAddr, strings.Join(events, "%20AND%20"), page, limit)
	if page == 0 {
		endpoint = fmt.Sprintf("%s/tx_search?query=\"%s\"&per_page=%d&order_by=\"desc\"", rpcAddr, strings.Join(events, "%20AND%20"), limit)
	}
	common.GetLogger().Info("[query-transaction] Entering transaction search: ", endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		common.GetLogger().Error("[query-transaction] Unable to connect to ", endpoint)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	response := new(tmJsonRPCTypes.RPCResponse)

	if err := json.Unmarshal(respBody, response); err != nil {
		common.GetLogger().Error("[query-transaction] Unable to decode response: ", err)
		return nil, err
	}

	if response.Error != nil {
		common.GetLogger().Error("[query-transaction] Error response:", response.Error.Message)
		return nil, errors.New(response.Error.Message)
	}

	result := new(tmTypes.ResultTxSearch)
	if err := tmjson.Unmarshal(response.Result, result); err != nil {
		common.GetLogger().Error("[query-transaction] Failed to unmarshal result:", err)
		return nil, fmt.Errorf("error unmarshalling result: %w", err)
	}

	return result, nil
}

func getBlockHeight(rpcAddr string, hash string) (int64, error) {
	endpoint := fmt.Sprintf("%s/tx?hash=%s", rpcAddr, hash)
	common.GetLogger().Info("[query-block] Entering block query: ", endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		common.GetLogger().Error("[query-block] Unable to connect to ", endpoint)
		return 0, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	response := new(tmJsonRPCTypes.RPCResponse)

	if err := json.Unmarshal(respBody, response); err != nil {
		common.GetLogger().Error("[query-block] Unable to decode response: ", err)
		return 0, err
	}
	if response.Error != nil {
		common.GetLogger().Error("[query-block] Error response:", response.Error.Message)
		return 0, errors.New(response.Error.Message)
	}

	result := new(tmTypes.ResultTx)
	if err := tmjson.Unmarshal(response.Result, result); err != nil {
		common.GetLogger().Error("[query-block] Failed to unmarshal result:", err)
		return 0, fmt.Errorf("error unmarshalling result: %w", err)
	}

	return result.Height, nil
}

func QueryBlockTransactionsHandler(rpcAddr string, r *http.Request, isWithdraw bool) (interface{}, interface{}, int) {
	err := r.ParseForm()
	if err != nil {
		common.GetLogger().Error("[query-transactions] Failed to parse query parameters:", err)
		return common.ServeError(0, "failed to parse query parameters", err.Error(), http.StatusBadRequest)
	}

	var (
		account   string = ""
		txType    string = ""
		last      string = ""
		sender    string = ""
		recipient string = ""
		pageSize  int    = 10
		page      int    = 1
	)

	account = r.FormValue("account")
	if account == "" {
		common.GetLogger().Error("[query-transactions] 'account' is not set")
		return common.ServeError(0, "'account' is not set", "", http.StatusBadRequest)
	}

	if isWithdraw {
		sender = account
	} else {
		recipient = account
	}

	if pageSizeStr := r.FormValue("page_size"); pageSizeStr != "" {
		if pageSize, err = strconv.Atoi(pageSizeStr); err != nil {
			common.GetLogger().Error("[query-transactions] Failed to parse parameter 'page_size': ", err)
			return common.ServeError(0, "failed to parse parameter 'page_size'", err.Error(), http.StatusBadRequest)
		}
		if pageSize < 1 || pageSize > 1000 {
			common.GetLogger().Error("[query-transactions] Invalid 'page_size' range: ", pageSize)
			return common.ServeError(0, "'page_size' should be 1 ~ 1000", "", http.StatusBadRequest)
		}
	}

	if pageStr := r.FormValue("page"); pageStr != "" {
		if page, err = strconv.Atoi(pageStr); err != nil {
			common.GetLogger().Error("[query-transactions] Failed to parse parameter 'page': ", err)
			return common.ServeError(0, "failed to parse parameter 'page'", err.Error(), http.StatusBadRequest)
		}
	}

	txType = r.FormValue("type")
	if txType == "" {
		txType = "all"
	}

	last = r.FormValue("last")

	var transactions []*tmTypes.ResultTx

	if last == "" {
		searchResult, err := SearchTxHashHandle(rpcAddr, sender, recipient, txType, page, pageSize, -1, -1, "")
		if err != nil {
			common.GetLogger().Error("[query-transactions] Failed to search transaction hash: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		transactions = searchResult.Txs
	} else {
		blockHeight, err := getBlockHeight(rpcAddr, last)
		if err != nil {
			common.GetLogger().Error("[query-transactions] Failed to query block height: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		// get current block
		searchResult, err := SearchTxHashHandle(rpcAddr, sender, recipient, txType, 0, pageSize, blockHeight, blockHeight, "")
		if err != nil {
			common.GetLogger().Error("[query-transactions] Failed to search transaction hash: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		beforeLast := true
		for _, tx := range searchResult.Txs {
			if !beforeLast {
				transactions = append(transactions, tx)
			}
			if fmt.Sprintf("0x%X", tx.Hash) == last {
				beforeLast = false
			}

			if len(transactions) == pageSize {
				break
			}
		}

		if len(transactions) < pageSize && blockHeight > 0 {
			searchResult, err := SearchTxHashHandle(rpcAddr, sender, recipient, txType, 0, pageSize-len(transactions), -1, blockHeight-1, "")
			if err != nil {
				common.GetLogger().Error("[query-transactions] Failed to search transaction hash: ", err)
				return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
			}

			transactions = append(transactions, searchResult.Txs...)
		}
	}

	var txResults = make(map[string]types.DepositWithdrawResult)

	for _, transaction := range transactions {
		tx, err := config.EncodingCg.TxConfig.TxDecoder()(transaction.Tx)
		if err != nil {
			common.GetLogger().Error("[query-transactions] Failed to decode transaction: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		blockTime, err := common.GetBlockTime(rpcAddr, transaction.Height)
		if err != nil {
			common.GetLogger().Error("[query-transactions] Block not found: ", transaction.Height)
			return common.ServeError(0, "", fmt.Sprintf("block not found: %d", transaction.Height), http.StatusInternalServerError)
		}

		logs, err := sdk.ParseABCILogs(transaction.TxResult.GetLog())
		if err != nil {
			common.GetLogger().Error("[query-transactions] Failed to parse ABCI logs: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		var txResponses []types.DepositWithdrawTransaction

		for index, msg := range tx.GetMsgs() {
			txType := kiratypes.MsgType(msg)

			var evMap = make(map[string]([]sdk.Attribute))
			for _, event := range logs[index].GetEvents() {
				evMap[event.GetType()] = event.GetAttributes()
			}

			if txType == "send" {
				msgSend := msg.(*bank.MsgSend)

				if isWithdraw && msgSend.FromAddress == account {
					for _, coin := range msgSend.Amount {
						txResponses = append(txResponses, types.DepositWithdrawTransaction{
							Address: msgSend.ToAddress,
							Type:    txType,
							Denom:   coin.GetDenom(),
							Amount:  coin.Amount.Int64(),
						})
					}
				}
				if !isWithdraw && msgSend.ToAddress == account {
					for _, coin := range msgSend.Amount {
						txResponses = append(txResponses, types.DepositWithdrawTransaction{
							Address: msgSend.FromAddress,
							Type:    txType,
							Denom:   coin.GetDenom(),
							Amount:  coin.Amount.Int64(),
						})
					}
				}
			} else if txType == "multisend" {
				msgMultiSend := msg.(*bank.MsgMultiSend)
				inputs := msgMultiSend.GetInputs()
				outputs := msgMultiSend.GetOutputs()
				if isWithdraw {
					for _, input := range inputs {
						if input.Address == account {
							if len(inputs) == 1 {
								for _, output := range outputs {
									for _, coin := range output.Coins {
										txResponses = append(txResponses, types.DepositWithdrawTransaction{
											Address: output.Address,
											Type:    txType,
											Denom:   coin.GetDenom(),
											Amount:  coin.Amount.Int64(),
										})
									}
								}
							} else if len(outputs) == 1 {
								for _, coin := range input.Coins {
									txResponses = append(txResponses, types.DepositWithdrawTransaction{
										Address: outputs[0].Address,
										Type:    txType,
										Denom:   coin.GetDenom(),
										Amount:  coin.Amount.Int64(),
									})
								}
							}
						}
					}
				} else {
					for _, output := range outputs {
						if output.Address == account {
							if len(inputs) == 1 {
								for _, coin := range output.Coins {
									txResponses = append(txResponses, types.DepositWithdrawTransaction{
										Address: inputs[0].Address,
										Type:    txType,
										Denom:   coin.GetDenom(),
										Amount:  coin.Amount.Int64(),
									})
								}
							} else if len(outputs) == 1 {
								for _, input := range inputs {
									for _, coin := range input.Coins {
										txResponses = append(txResponses, types.DepositWithdrawTransaction{
											Address: input.Address,
											Type:    txType,
											Denom:   coin.GetDenom(),
											Amount:  coin.Amount.Int64(),
										})
									}
								}
							}
						}
					}
				}
			} else if txType == "create_validator" {
				createValidatorMsg := msg.(*staking.MsgCreateValidator)

				if isWithdraw && createValidatorMsg.DelegatorAddress == account {
					txResponses = append(txResponses, types.DepositWithdrawTransaction{
						Address: createValidatorMsg.ValidatorAddress,
						Type:    txType,
						Denom:   createValidatorMsg.Value.Denom,
						Amount:  createValidatorMsg.Value.Amount.Int64(),
					})
				} else if !isWithdraw && createValidatorMsg.ValidatorAddress == account {
					txResponses = append(txResponses, types.DepositWithdrawTransaction{
						Address: createValidatorMsg.DelegatorAddress,
						Type:    txType,
						Denom:   createValidatorMsg.Value.Denom,
						Amount:  createValidatorMsg.Value.Amount.Int64(),
					})
				}
			} else if txType == "delegate" {
				delegateMsg := msg.(*staking.MsgDelegate)

				if isWithdraw && delegateMsg.DelegatorAddress == account {
					txResponses = append(txResponses, types.DepositWithdrawTransaction{
						Address: delegateMsg.ValidatorAddress,
						Type:    txType,
						Denom:   delegateMsg.Amount.Denom,
						Amount:  delegateMsg.Amount.Amount.Int64(),
					})
				} else if !isWithdraw && delegateMsg.ValidatorAddress == account {
					txResponses = append(txResponses, types.DepositWithdrawTransaction{
						Address: delegateMsg.DelegatorAddress,
						Type:    txType,
						Denom:   delegateMsg.Amount.Denom,
						Amount:  delegateMsg.Amount.Amount.Int64(),
					})
				}
			} else if txType == "begin_redelegate" {
				reDelegateMsg := msg.(*staking.MsgBeginRedelegate)

				if isWithdraw && reDelegateMsg.ValidatorSrcAddress == account {
					txResponses = append(txResponses, types.DepositWithdrawTransaction{
						Address: reDelegateMsg.ValidatorDstAddress,
						Type:    txType,
						Denom:   reDelegateMsg.Amount.Denom,
						Amount:  reDelegateMsg.Amount.Amount.Int64(),
					})
				} else if !isWithdraw && reDelegateMsg.ValidatorDstAddress == account {
					txResponses = append(txResponses, types.DepositWithdrawTransaction{
						Address: reDelegateMsg.ValidatorSrcAddress,
						Type:    txType,
						Denom:   reDelegateMsg.Amount.Denom,
						Amount:  reDelegateMsg.Amount.Amount.Int64(),
					})
				}
			} else if txType == "begin_unbonding" {
				unDelegateMsg := msg.(*staking.MsgUndelegate)

				if isWithdraw && unDelegateMsg.ValidatorAddress == account {
					txResponses = append(txResponses, types.DepositWithdrawTransaction{
						Address: unDelegateMsg.DelegatorAddress,
						Type:    txType,
						Denom:   unDelegateMsg.Amount.Denom,
						Amount:  unDelegateMsg.Amount.Amount.Int64(),
					})
				} else if !isWithdraw && unDelegateMsg.DelegatorAddress == account {
					txResponses = append(txResponses, types.DepositWithdrawTransaction{
						Address: unDelegateMsg.ValidatorAddress,
						Type:    txType,
						Denom:   unDelegateMsg.Amount.Denom,
						Amount:  unDelegateMsg.Amount.Amount.Int64(),
					})
				}
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

				if isWithdraw && withdrawDelegatorRewardMsg.ValidatorAddress == account {
					txResponses = append(txResponses, types.DepositWithdrawTransaction{
						Address: withdrawDelegatorRewardMsg.DelegatorAddress,
						Type:    txType,
						Denom:   coin.Denom,
						Amount:  coin.Amount.Int64(),
					})
				} else if !isWithdraw && withdrawDelegatorRewardMsg.DelegatorAddress == account {
					txResponses = append(txResponses, types.DepositWithdrawTransaction{
						Address: withdrawDelegatorRewardMsg.ValidatorAddress,
						Type:    txType,
						Denom:   coin.Denom,
						Amount:  coin.Amount.Int64(),
					})
				}
			}
		}

		txResults[fmt.Sprintf("0x%X", transaction.Hash)] = types.DepositWithdrawResult{
			Time: blockTime,
			Txs:  txResponses,
		}
	}

	res := struct {
		Transactions map[string]types.DepositWithdrawResult `json:"transactions"`
		TotalCount   int                                    `json:"total_count"`
	}{}

	searchResult, err := SearchTxHashHandle(rpcAddr, sender, recipient, txType, 0, pageSize, -1, -1, "")
	res.TotalCount = searchResult.TotalCount
	res.Transactions = txResults

	return res, nil, http.StatusOK
}

// QueryWithdraws is a function to query withdraw transactions.
func QueryWithdraws(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-withdraws] Entering withdraws query")

		if !common.RPCMethods["GET"][config.QueryWithdraws].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryWithdraws].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-withdraws] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = QueryBlockTransactionsHandler(rpcAddr, r, true)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryStatus].CachingEnabled)
	}
}

// QueryDeposits is a function to query deposit transactions.
func QueryDeposits(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Error("[query-deposits] Entering withdraws query")

		if !common.RPCMethods["GET"][config.QueryDeposits].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryDeposits].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-deposits] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = QueryBlockTransactionsHandler(rpcAddr, r, false)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryStatus].CachingEnabled)
	}
}

func searchUnconfirmed(rpcAddr string, limit string) (*tmTypes.ResultUnconfirmedTxs, error) {
	endpoint := fmt.Sprintf("%s/unconfirmed_txs?limit=%s", rpcAddr, limit)
	fmt.Println(endpoint)
	common.GetLogger().Info("[query-unconfirmed-txs] Entering transaction search: ", endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		common.GetLogger().Error("[query-unconfirmed-txs] Unable to connect to ", endpoint)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	response := new(tmJsonRPCTypes.RPCResponse)

	if err := json.Unmarshal(respBody, response); err != nil {
		common.GetLogger().Error("[query-unconfirmed-txs] Unable to decode response: ", err)
		return nil, err
	}

	if response.Error != nil {
		common.GetLogger().Error("[query-unconfirmed-txs] Error response:", response.Error.Message)
		return nil, errors.New(response.Error.Message)
	}

	result := new(tmTypes.ResultUnconfirmedTxs)
	if err := tmjson.Unmarshal(response.Result, result); err != nil {
		common.GetLogger().Error("[query-unconfirmed-txs] Failed to unmarshal result:", err)
		return nil, fmt.Errorf("error unmarshalling result: %w", err)
	}

	return result, nil
}

func queryUnconfirmedTransactionsHandler(rpcAddr string, r *http.Request) (interface{}, interface{}, int) {
	limit := r.FormValue("limit")
	result, err := searchUnconfirmed(rpcAddr, limit)
	if err != nil {
		common.GetLogger().Error("[query-unconfirmed-txs] Failed to query unconfirmed txs: %w ", err)
		return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
	}

	response := struct {
		Count      int                                  `json:"n_txs"`
		Total      int                                  `json:"total"`
		TotalBytes int64                                `json:"total_bytes"`
		Txs        []types.TransactionUnconfirmedResult `json:"txs"`
	}{}

	response.Count = result.Count
	response.Total = result.Total
	response.TotalBytes = result.TotalBytes
	response.Txs = make([]types.TransactionUnconfirmedResult, 0)

	for _, tx := range result.Txs {
		decodedTx, err := config.EncodingCg.TxConfig.TxDecoder()(tx)
		if err != nil {
			common.GetLogger().Error("[post-unconfirmed-txs] Failed to decode transaction: ", err)
			return common.ServeError(0, "failed to decode signed TX", err.Error(), http.StatusBadRequest)
		}

		txResult, ok := decodedTx.(signing.Tx)
		if !ok {
			common.GetLogger().Error("[post-unconfirmed-txs] Failed to decode transaction")
			return common.ServeError(0, "failed to decode signed TX", err.Error(), http.StatusBadRequest)
		}

		signature, _ := txResult.GetSignaturesV2()

		var msgs []types.TxMsg = make([]types.TxMsg, 0)

		for _, msg := range txResult.GetMsgs() {
			msgs = append(msgs, types.TxMsg{
				Type: kiratypes.MsgType(msg),
				Data: msg,
			})
		}

		response.Txs = append(response.Txs, types.TransactionUnconfirmedResult{
			Msgs:      msgs,
			Fees:      txResult.GetFee(),
			Gas:       txResult.GetGas(),
			Signature: signature,
			Memo:      txResult.GetMemo(),
		})
	}

	return response, nil, http.StatusOK
}

// QueryUnconfirmedTxs is a function to query unconfirmed transactions.
func QueryUnconfirmedTxs(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Error("[query-unconfirmed-txs] Entering query")

		if !common.RPCMethods["GET"][config.QueryUnconfirmedTxs].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			response.Response, response.Error, statusCode = queryUnconfirmedTransactionsHandler(rpcAddr, r)
		}

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}
