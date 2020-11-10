package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	interx "github.com/KiraCore/sekai/INTERX/config"
	functions "github.com/KiraCore/sekai/INTERX/functions"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmJsonRPCTypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

const (
	queryFunctions = "/api/functions"
	queryWithdraws = "/api/withdraws"
	queryDeposits  = "/api/deposits"
)

// Transaction is a struct to be used for transaction
type Transaction struct {
	Address string `json:"address"`
	Type    string `json:"type"`
	Denom   string `json:"denom,omitempty"`
	Amount  int64  `json:"amount"`
}

// TransactionResult is a struct to be used for query transaction response
type TransactionResult struct {
	Time int64         `json:"time"`
	Txs  []Transaction `json:"txs"`
}

// RegisterTxQueryRoutes registers tx query routers.
func RegisterTxQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc("/api/functions/{address}", QueryFunctions(rpcAddr)).Methods(GET)
	r.HandleFunc(queryWithdraws, QueryWithdraws(rpcAddr)).Methods(GET)
	r.HandleFunc(queryDeposits, QueryDeposits(rpcAddr)).Methods(GET)

	AddRPCMethod(GET, queryFunctions, "This is an API to list functions and metadata.", true)
	AddRPCMethod(GET, queryWithdraws, "This is an API to query withdraw transactions.", true)
	AddRPCMethod(GET, queryDeposits, "This is an API to query deposit transactions.", true)
}

func queryFunctionsHandle(rpcAddr string, address string) (interface{}, interface{}, int) {
	permittedTxTypes, err := GetPermittedTxTypes(rpcAddr, address)

	if err != nil {
		return ServeError(0, "", err.Error(), http.StatusBadRequest)
	}

	type FunctionsResponse struct {
		FunctionID  string                                 `json:"function_id"`
		Description string                                 `json:"description"`
		Parameters  map[string]functions.FunctionMetaField `json:"parameters"`
	}

	response := map[string]FunctionsResponse{}

	for k, v := range permittedTxTypes {
		if function, ok := functions.AllFunctions[k]; ok {
			response[k] = FunctionsResponse{
				FunctionID:  v,
				Description: function.Description,
				Parameters:  function.Parameters,
			}
		}
	}

	return response, nil, http.StatusOK
}

// QueryFunctions is a function to list functions and metadata.
func QueryFunctions(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		bech32addr := queries["address"]
		request := InterxRequest{
			Method:   r.Method,
			Endpoint: queryBalances,
			Params:   []byte(bech32addr),
		}
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response, response.Error, statusCode = queryFunctionsHandle(rpcAddr, bech32addr)

		WrapResponse(w, request, *response, statusCode, false)
	}
}

func toSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func parseTxType(txType string) string {
	return toSnakeCase(txType)
}

func searchTxHashHandle(rpcAddr string, sender string, recipient string, txType string, limit int, txMinHeight int64, txMaxHeight int64) (*tmTypes.ResultTxSearch, error) {
	var events = make([]string, 0, 4)

	if sender != "" {
		events = append(events, fmt.Sprintf("transfer.sender='%s'", sender))
	}

	if recipient != "" {
		events = append(events, fmt.Sprintf("transfer.recipient='%s'", recipient))
	}

	if txType != "all" {
		events = append(events, fmt.Sprintf("message.action='%s'", txType))
	}

	if txMinHeight >= 0 {
		events = append(events, fmt.Sprintf("tx.height>=%d", txMinHeight))
	}

	if txMaxHeight >= 0 {
		events = append(events, fmt.Sprintf("tx.height<=%d", txMaxHeight))
	}

	// search transactions
	resp, err := http.Get(fmt.Sprintf("%s/tx_search?query=\"%s\"&per_page=%d&order_by=\"desc\"", rpcAddr, strings.Join(events, "%20AND%20"), limit))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	response := new(tmJsonRPCTypes.RPCResponse)

	if err := json.Unmarshal(respBody, response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, errors.New(response.Error.Message)
	}

	result := new(tmTypes.ResultTxSearch)
	if err := tmjson.Unmarshal(response.Result, result); err != nil {
		return nil, fmt.Errorf("error unmarshalling result: %w", err)
	}

	return result, nil
}

func getBlockHeight(rpcAddr string, hash string) (int64, error) {
	resp, err := http.Get(fmt.Sprintf("%s/tx?hash=%s", rpcAddr, hash))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	response := new(tmJsonRPCTypes.RPCResponse)

	if err := json.Unmarshal(respBody, response); err != nil {
		return 0, err
	}
	if response.Error != nil {
		return 0, errors.New(response.Error.Message)
	}

	result := new(tmTypes.ResultTx)
	if err := tmjson.Unmarshal(response.Result, result); err != nil {
		return 0, fmt.Errorf("error unmarshalling result: %w", err)
	}

	return result.Height, nil
}

func queryTransactionsHandler(rpcAddr string, r *http.Request, isWithdraw bool) (interface{}, interface{}, int) {
	err := r.ParseForm()
	if err != nil {
		return ServeError(0, "failed to parse query parameters", err.Error(), http.StatusBadRequest)
	}

	var (
		account   string = ""
		txType    string = ""
		last      string = ""
		sender    string = ""
		recipient string = ""
		limit     int    = 10
	)

	account = r.FormValue("account")
	if account == "" {
		return ServeError(0, "'account' is not set", "", http.StatusBadRequest)
	}

	if isWithdraw {
		sender = account
	} else {
		recipient = account
	}

	if maxStr := r.FormValue("max"); maxStr != "" {
		if limit, err = strconv.Atoi(maxStr); err != nil {
			return ServeError(0, "failed to parse parameter 'max'", err.Error(), http.StatusBadRequest)
		}
		if limit < 1 || limit > 1000 {
			return ServeError(0, "'max' should be 1 ~ 1000", "", http.StatusBadRequest)
		}
	}

	txType = r.FormValue("type")
	if txType == "" {
		txType = "all"
	}
	txType = parseTxType(txType)

	last = r.FormValue("last")

	var transactions []*tmTypes.ResultTx

	if last == "" {
		searchResult, err := searchTxHashHandle(rpcAddr, sender, recipient, txType, limit, -1, -1)
		if err != nil {
			return ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		transactions = searchResult.Txs
	} else {
		type TxResult struct {
			Height string `json:"height"`
		}

		blockHeight, err := getBlockHeight(rpcAddr, last)
		if err != nil {
			return ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		// get current block
		searchResult, err := searchTxHashHandle(rpcAddr, sender, recipient, txType, limit, blockHeight, blockHeight)
		if err != nil {
			return ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		beforeLast := true
		for _, tx := range searchResult.Txs {
			if beforeLast == false {
				transactions = append(transactions, tx)
			}
			if fmt.Sprintf("0x%X", tx.Hash) == last {
				beforeLast = false
			}

			if len(transactions) == limit {
				break
			}
		}

		if len(transactions) < limit && blockHeight > 0 {
			searchResult, err := searchTxHashHandle(rpcAddr, sender, recipient, txType, limit-len(transactions), -1, blockHeight-1)
			if err != nil {
				return ServeError(0, "", err.Error(), http.StatusInternalServerError)
			}

			for _, tx := range searchResult.Txs {
				transactions = append(transactions, tx)
			}
		}
	}

	var response = make(map[string]TransactionResult)

	for _, transaction := range transactions {
		tx, err := interx.EncodingCg.TxConfig.TxDecoder()(transaction.Tx)
		if err != nil {
			return ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		blockTime, err := GetBlockTime(rpcAddr, transaction.Height)
		if err != nil {
			return ServeError(0, "", fmt.Sprintf("block not found: %d", transaction.Height), http.StatusInternalServerError)
		}

		logs, err := sdk.ParseABCILogs(transaction.TxResult.GetLog())
		if err != nil {
			return ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		var txResponses []Transaction

		for index, msg := range tx.GetMsgs() {
			txType := msg.Type()

			var evMap = make(map[string]([]sdk.Attribute))
			for _, event := range logs[index].GetEvents() {
				evMap[event.GetType()] = event.GetAttributes()
			}

			if txType == "send" {
				msgSend := msg.(*bank.MsgSend)

				if isWithdraw && msgSend.GetFromAddress().String() == account {
					for _, coin := range msgSend.GetAmount() {
						txResponses = append(txResponses, Transaction{
							Address: msgSend.GetToAddress().String(),
							Type:    txType,
							Denom:   coin.GetDenom(),
							Amount:  coin.Amount.Int64(),
						})
					}
				}
				if !isWithdraw && msgSend.GetToAddress().String() == account {
					for _, coin := range msgSend.GetAmount() {
						txResponses = append(txResponses, Transaction{
							Address: msgSend.GetFromAddress().String(),
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
						if input.GetAddress().String() == account {
							// found input
							if len(inputs) == 1 {
								for _, output := range outputs {
									for _, coin := range output.GetCoins() {
										txResponses = append(txResponses, Transaction{
											Address: output.GetAddress().String(),
											Type:    txType,
											Denom:   coin.GetDenom(),
											Amount:  coin.Amount.Int64(),
										})
									}
								}
							} else if len(outputs) == 1 {
								for _, coin := range input.GetCoins() {
									txResponses = append(txResponses, Transaction{
										Address: outputs[0].GetAddress().String(),
										Type:    txType,
										Denom:   coin.GetDenom(),
										Amount:  coin.Amount.Int64(),
									})
								}
							} else {
							}
						}
					}
				} else {
					for _, output := range outputs {
						if output.GetAddress().String() == account {
							// found output
							if len(inputs) == 1 {
								for _, coin := range output.GetCoins() {
									txResponses = append(txResponses, Transaction{
										Address: inputs[0].GetAddress().String(),
										Type:    txType,
										Denom:   coin.GetDenom(),
										Amount:  coin.Amount.Int64(),
									})
								}
							} else if len(outputs) == 1 {
								for _, input := range inputs {
									for _, coin := range input.GetCoins() {
										txResponses = append(txResponses, Transaction{
											Address: input.GetAddress().String(),
											Type:    txType,
											Denom:   coin.GetDenom(),
											Amount:  coin.Amount.Int64(),
										})
									}
								}
							} else {
							}
						}
					}
				}
			} else if txType == "create_validator" {
				createValidatorMsg := msg.(*staking.MsgCreateValidator)

				if isWithdraw && createValidatorMsg.GetDelegatorAddress().String() == account {
					txResponses = append(txResponses, Transaction{
						Address: createValidatorMsg.GetValidatorAddress().String(),
						Type:    txType,
						Denom:   createValidatorMsg.GetValue().Denom,
						Amount:  createValidatorMsg.GetValue().Amount.Int64(),
					})
				} else if !isWithdraw && createValidatorMsg.GetValidatorAddress().String() == account {
					txResponses = append(txResponses, Transaction{
						Address: createValidatorMsg.GetDelegatorAddress().String(),
						Type:    txType,
						Denom:   createValidatorMsg.GetValue().Denom,
						Amount:  createValidatorMsg.GetValue().Amount.Int64(),
					})
				}
			} else if txType == "delegate" {
				delegateMsg := msg.(*staking.MsgDelegate)

				if isWithdraw && delegateMsg.GetDelegatorAddress().String() == account {
					txResponses = append(txResponses, Transaction{
						Address: delegateMsg.GetValidatorAddress().String(),
						Type:    txType,
						Denom:   delegateMsg.GetAmount().Denom,
						Amount:  delegateMsg.GetAmount().Amount.Int64(),
					})
				} else if !isWithdraw && delegateMsg.GetValidatorAddress().String() == account {
					txResponses = append(txResponses, Transaction{
						Address: delegateMsg.GetDelegatorAddress().String(),
						Type:    txType,
						Denom:   delegateMsg.GetAmount().Denom,
						Amount:  delegateMsg.GetAmount().Amount.Int64(),
					})
				}
			} else if txType == "begin_redelegate" {
				reDelegateMsg := msg.(*staking.MsgBeginRedelegate)

				if isWithdraw && reDelegateMsg.GetValidatorSrcAddress().String() == account {
					txResponses = append(txResponses, Transaction{
						Address: reDelegateMsg.GetValidatorDstAddress().String(),
						Type:    txType,
						Denom:   reDelegateMsg.GetAmount().Denom,
						Amount:  reDelegateMsg.GetAmount().Amount.Int64(),
					})
				} else if !isWithdraw && reDelegateMsg.GetValidatorDstAddress().String() == account {
					txResponses = append(txResponses, Transaction{
						Address: reDelegateMsg.GetValidatorSrcAddress().String(),
						Type:    txType,
						Denom:   reDelegateMsg.GetAmount().Denom,
						Amount:  reDelegateMsg.GetAmount().Amount.Int64(),
					})
				}
			} else if txType == "begin_unbonding" {
				unDelegateMsg := msg.(*staking.MsgUndelegate)

				if isWithdraw && unDelegateMsg.GetValidatorAddress().String() == account {
					txResponses = append(txResponses, Transaction{
						Address: unDelegateMsg.GetDelegatorAddress().String(),
						Type:    txType,
						Denom:   unDelegateMsg.GetAmount().Denom,
						Amount:  unDelegateMsg.GetAmount().Amount.Int64(),
					})
				} else if !isWithdraw && unDelegateMsg.GetDelegatorAddress().String() == account {
					txResponses = append(txResponses, Transaction{
						Address: unDelegateMsg.GetValidatorAddress().String(),
						Type:    txType,
						Denom:   unDelegateMsg.GetAmount().Denom,
						Amount:  unDelegateMsg.GetAmount().Amount.Int64(),
					})
				}
			} else if txType == "withdraw_delegator_reward" {
				var coin sdk.Coin
				if v, found := evMap["withdraw_rewards"]; found && len(v) >= 2 {
					if v[0].GetKey() == "amount" {
						coin, _ = sdk.ParseCoin(v[0].GetValue())
					} else if v[1].GetKey() == "amount" {
						coin, _ = sdk.ParseCoin(v[1].GetValue())
					}
				}

				withdrawDelegatorRewardMsg := msg.(*distribution.MsgWithdrawDelegatorReward)

				if isWithdraw && withdrawDelegatorRewardMsg.GetValidatorAddress().String() == account {
					txResponses = append(txResponses, Transaction{
						Address: withdrawDelegatorRewardMsg.GetDelegatorAddress().String(),
						Type:    txType,
						Denom:   coin.Denom,
						Amount:  coin.Amount.Int64(),
					})
				} else if !isWithdraw && withdrawDelegatorRewardMsg.GetDelegatorAddress().String() == account {
					txResponses = append(txResponses, Transaction{
						Address: withdrawDelegatorRewardMsg.GetValidatorAddress().String(),
						Type:    txType,
						Denom:   coin.Denom,
						Amount:  coin.Amount.Int64(),
					})
				}
			}
		}

		response[fmt.Sprintf("0x%X", transaction.Hash)] = TransactionResult{
			Time: blockTime,
			Txs:  txResponses,
		}
	}

	fmt.Println(response)

	return response, nil, http.StatusOK
}

// QueryWithdraws is a function to query withdraw transactions.
func QueryWithdraws(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		fmt.Println(request)

		if !rpcMethods[GET][queryWithdraws].Enabled {
			response.Response, response.Error, statusCode = ServeError(0, "", "", http.StatusForbidden)
		} else {
			if rpcMethods[GET][queryWithdraws].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = queryTransactionsHandler(rpcAddr, r, true)
		}

		WrapResponse(w, request, *response, statusCode, rpcMethods[GET][queryStatus].CachingEnabled)
	}
}

// QueryDeposits is a function to query deposit transactions.
func QueryDeposits(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !rpcMethods[GET][queryDeposits].Enabled {
			response.Response, response.Error, statusCode = ServeError(0, "", "", http.StatusForbidden)
		} else {
			if rpcMethods[GET][queryDeposits].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = queryTransactionsHandler(rpcAddr, r, false)
		}

		WrapResponse(w, request, *response, statusCode, rpcMethods[GET][queryStatus].CachingEnabled)
	}
}
