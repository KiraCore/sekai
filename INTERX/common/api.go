package common

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"

	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/database"
	"github.com/KiraCore/sekai/INTERX/global"
	"github.com/KiraCore/sekai/INTERX/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmJsonRPCTypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

// MakeTendermintRPCRequest is a function to make GET request
func MakeTendermintRPCRequest(rpcAddr string, url string, query string) (interface{}, interface{}, int) {
	endpoint := fmt.Sprintf("%s%s?%s", rpcAddr, url, query)
	// GetLogger().Info("[rpc-call] Entering tendermint rpc call: ", endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		GetLogger().Error("[rpc-call] Unable to connect to ", endpoint)
		return ServeError(0, "", err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	response := new(types.RPCResponse)
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		GetLogger().Error("[rpc-call] Unable to decode response: : ", err)
		return nil, err.Error(), resp.StatusCode
	}

	return response.Result, response.Error, resp.StatusCode
}

// MakeGetRequest is a function to make GET request
func MakeGetRequest(rpcAddr string, url string, query string) (Result interface{}, Error interface{}, StatusCode int) {
	endpoint := fmt.Sprintf("%s%s?%s", rpcAddr, url, query)
	// GetLogger().Info("[rpc-call] Entering rpc call: ", endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		GetLogger().Error("[rpc-call] Unable to connect to ", endpoint)
		return ServeError(0, "", err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	StatusCode = resp.StatusCode

	err = json.NewDecoder(resp.Body).Decode(&Result)
	if err != nil {
		GetLogger().Error("[rpc-call] Unable to decode response: : ", err)
		Error = err.Error()
	}

	return Result, Error, StatusCode
}

// DownloadResponseToFile is a function to save GET response as a file
func DownloadResponseToFile(rpcAddr string, url string, query string, filepath string) error {
	endpoint := fmt.Sprintf("%s%s?%s", rpcAddr, url, query)
	// GetLogger().Info("[rpc-call] Entering rpc call: ", endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		GetLogger().Error("[rpc-call] Unable to connect to ", endpoint)
		return err
	}
	defer resp.Body.Close()

	fileout, _ := os.Create(filepath)
	defer fileout.Close()

	global.Mutex.Lock()
	io.Copy(fileout, resp.Body)
	global.Mutex.Unlock()

	return err
}

func makePostRequest(r *http.Request) (*types.RPCResponse, error) {
	endpoint := fmt.Sprintf("%s%s", r.Host, r.URL)
	// GetLogger().Info("[rpc-call] Entering rpc call: ", endpoint)

	resp, err := http.PostForm(endpoint, r.Form)
	if err != nil {
		GetLogger().Error("[rpc-call] Unable to connect to ", endpoint)
		return nil, err
	}
	defer resp.Body.Close()

	result := new(types.RPCResponse)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		GetLogger().Error("[rpc-call] Unable to decode response: : ", err)
		return nil, err
	}

	return result, nil
}

// GetAccountBalances is a function to get balances of an address
func GetAccountBalances(gwCosmosmux *runtime.ServeMux, r *http.Request, bech32addr string) []types.Coin {
	_, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		GetLogger().Error("[grpc-call] Invalid bech32addr: ", bech32addr)
		return nil
	}

	r.URL.Path = fmt.Sprintf("/api/cosmos/bank/balances/%s", base64.URLEncoding.EncodeToString([]byte(bech32addr)))
	r.URL.RawQuery = ""
	r.Method = "GET"

	// GetLogger().Info("[grpc-call] Entering grpc call: ", r.URL.Path)

	recorder := httptest.NewRecorder()
	gwCosmosmux.ServeHTTP(recorder, r)
	resp := recorder.Result()

	type BalancesResponse struct {
		Balances []types.Coin `json:"balances"`
	}

	result := BalancesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		GetLogger().Error("[grpc-call] Unable to decode response: ", err)
	}

	return result.Balances
}

// GetAccountNumberSequence is a function to get AccountNumber and Sequence
func GetAccountNumberSequence(gwCosmosmux *runtime.ServeMux, r *http.Request, bech32addr string) (uint64, uint64) {
	_, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		GetLogger().Error("[grpc-call] Invalid bech32addr: ", bech32addr)
		return 0, 0
	}

	r.URL.Path = fmt.Sprintf("/api/cosmos/auth/accounts/%s", base64.URLEncoding.EncodeToString([]byte(bech32addr)))
	r.URL.RawQuery = ""
	r.Method = "GET"

	// GetLogger().Info("[grpc-call] Entering grpc call: ", r.URL.Path)

	recorder := httptest.NewRecorder()
	gwCosmosmux.ServeHTTP(recorder, r)
	resp := recorder.Result()

	type QueryAccountResponse struct {
		Account struct {
			Address       string `json:"addresss"`
			PubKey        string `json:"pubKey"`
			AccountNumber string `json:"account_number"`
			Sequence      string `json:"sequence"`
		} `json:"account"`
	}
	result := QueryAccountResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		GetLogger().Error("[grpc-call] Unable to decode response: ", err)
	}

	accountNumber, _ := strconv.ParseInt(result.Account.AccountNumber, 10, 64)
	sequence, _ := strconv.ParseInt(result.Account.Sequence, 10, 64)

	return uint64(accountNumber), uint64(sequence)
}

// BroadcastTransaction is a function to post transaction, returns txHash
func BroadcastTransaction(rpcAddr string, txBytes []byte) (string, error) {
	endpoint := fmt.Sprintf("%s/broadcast_tx_async?tx=0x%X", rpcAddr, txBytes)
	// GetLogger().Info("[rpc-call] Entering rpc call: ", endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		GetLogger().Error("[rpc-call] Unable to connect to ", endpoint)
		return "", err
	}
	defer resp.Body.Close()

	type RPCTempResponse struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  struct {
			Height string `json:"height"`
			Hash   string `json:"hash"`
		} `json:"result,omitempty"`
		Error struct {
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}

	result := new(RPCTempResponse)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		GetLogger().Error("[rpc-call] Unable to decode response: ", err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		GetLogger().Error("[rpc-call] Unable to broadcast transaction: ", result.Error.Message)
		return "", errors.New(result.Error.Message)
	}

	return result.Result.Hash, nil
}

// GetPermittedTxTypes is a function to get all permitted tx types and function ids
func GetPermittedTxTypes(rpcAddr string, account string) (map[string]string, error) {
	permittedTxTypes := map[string]string{}
	permittedTxTypes["ExampleTx"] = "123"
	return permittedTxTypes, nil
}

// GetBlockTime is a function to get block time
func GetBlockTime(rpcAddr string, height int64) (int64, error) {
	blockTime, err := database.GetBlockTime(height)
	if err == nil {
		return blockTime, nil
	}

	endpoint := fmt.Sprintf("%s/block?height=%d", rpcAddr, height)
	// GetLogger().Info("[rpc-call] Entering rpc call: ", endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		GetLogger().Error("[rpc-call] Unable to connect to ", endpoint)
		return 0, fmt.Errorf("block not found: %d", height)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	response := new(tmJsonRPCTypes.RPCResponse)

	if err := json.Unmarshal(respBody, response); err != nil {
		GetLogger().Error("[rpc-call] Unable to decode response: ", err)
		return 0, err
	}

	if response.Error != nil {
		GetLogger().Error("[rpc-call] Block not found: ", height)
		return 0, fmt.Errorf("block not found: %d", height)
	}

	result := new(tmTypes.ResultBlock)
	if err := tmjson.Unmarshal(response.Result, result); err != nil {
		GetLogger().Error("[rpc-call] Unable to decode response: ", err)
		return 0, err
	}

	blockTime = result.Block.Header.Time.Unix()

	// save block time
	database.AddBlockTime(height, blockTime)

	// save block nano time
	database.AddBlockNanoTime(height, result.Block.Header.Time.UnixNano())

	return blockTime, nil
}

// GetBlockNanoTime is a function to get block nano time
func GetBlockNanoTime(rpcAddr string, height int64) (int64, error) {
	blockTime, err := database.GetBlockNanoTime(height)
	if err == nil {
		return blockTime, nil
	}

	endpoint := fmt.Sprintf("%s/block?height=%d", rpcAddr, height)
	// GetLogger().Info("[rpc-call] Entering rpc call: ", endpoint)

	resp, err := http.Get(endpoint)
	if err != nil {
		GetLogger().Error("[rpc-call] Unable to connect to ", endpoint)
		return 0, fmt.Errorf("block not found: %d", height)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	response := new(tmJsonRPCTypes.RPCResponse)

	if err := json.Unmarshal(respBody, response); err != nil {
		GetLogger().Error("[rpc-call] Unable to decode response: ", err)
		return 0, err
	}

	if response.Error != nil {
		GetLogger().Error("[rpc-call] Block not found: ", height)
		return 0, fmt.Errorf("block not found: %d", height)
	}

	result := new(tmTypes.ResultBlock)
	if err := tmjson.Unmarshal(response.Result, result); err != nil {
		GetLogger().Error("[rpc-call] Unable to decode response: ", err)
		return 0, err
	}

	blockTime = result.Block.Header.Time.UnixNano()

	// save block time
	database.AddBlockTime(height, result.Block.Header.Time.Unix())

	// save block nano time
	database.AddBlockNanoTime(height, blockTime)

	return blockTime, nil
}

// GetTokenAliases is a function to get token aliases
func GetTokenAliases(gwCosmosmux *runtime.ServeMux, r *http.Request) []types.TokenAlias {
	tokens, err := database.GetTokenAliases()
	if err == nil {
		return tokens
	}

	r.URL.Path = config.QueryKiraTokensAliases
	r.URL.RawQuery = ""
	r.Method = "GET"

	// GetLogger().Info("[grpc-call] Entering grpc call: ", r.URL.Path)

	recorder := httptest.NewRecorder()
	gwCosmosmux.ServeHTTP(recorder, r)
	resp := recorder.Result()

	type TokenAliasesResponse struct {
		Data []types.TokenAlias `json:"data"`
	}

	result := TokenAliasesResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		GetLogger().Error("[grpc-call] Unable to decode response: ", err)
	}

	// save block time
	database.AddTokenAliases(result.Data)

	return result.Data
}

// GetTokenSupply is a function to get token supply
func GetTokenSupply(gwCosmosmux *runtime.ServeMux, r *http.Request) []types.TokenSupply {
	r.URL.Path = config.QueryTotalSupply
	r.URL.RawQuery = ""
	r.Method = "GET"

	// GetLogger().Info("[grpc-call] Entering grpc call: ", r.URL.Path)

	recorder := httptest.NewRecorder()
	gwCosmosmux.ServeHTTP(recorder, r)
	resp := recorder.Result()

	type TokenAliasesResponse struct {
		Supply []types.TokenSupply `json:"supply"`
	}

	result := TokenAliasesResponse{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		GetLogger().Error("[grpc-call] Unable to decode response: ", err)
	}

	return result.Supply
}

func GetKiraStatus(rpcAddr string) *types.KiraStatus {
	success, _, _ := MakeTendermintRPCRequest(rpcAddr, "/status", "")

	if success != nil {
		result := types.KiraStatus{}

		byteData, err := json.Marshal(success)
		if err != nil {
			GetLogger().Error("[kira-status] Invalid response format", err)
		}

		err = json.Unmarshal(byteData, &result)
		if err != nil {
			GetLogger().Error("[kira-status] Invalid response format", err)
		}

		return &result
	}

	return nil
}

func GetInterxStatus(interxAddr string) *types.InterxStatus {
	success, _, _ := MakeGetRequest(interxAddr, "/api/status", "")

	if success != nil {
		result := types.InterxStatus{}

		byteData, err := json.Marshal(success)
		if err != nil {
			GetLogger().Error("[interx-status] Invalid response format", err)
			return nil
		}

		err = json.Unmarshal(byteData, &result)
		if err != nil {
			GetLogger().Error("[interx-status] Invalid response format", err)
			return nil
		}

		return &result
	}

	return nil
}

func GetSnapshotInfo(interxAddr string) *types.SnapShotChecksumResponse {
	success, _, _ := MakeGetRequest(interxAddr, "/api/snapshot_info", "")

	if success != nil {
		result := types.SnapShotChecksumResponse{}

		byteData, err := json.Marshal(success)
		if err != nil {
			GetLogger().Error("[interx-snapshot_info] Invalid response format", err)
			return nil
		}

		err = json.Unmarshal(byteData, &result)
		if err != nil {
			GetLogger().Error("[interx-snapshot_info] Invalid response format", err)
			return nil
		}

		return &result
	}

	return nil
}
