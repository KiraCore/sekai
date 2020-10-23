package gateway

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	common "github.com/KiraCore/sekai/INTERX/common"
	interx "github.com/KiraCore/sekai/INTERX/config"
	tasks "github.com/KiraCore/sekai/INTERX/tasks"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/crypto/blake2b"
)

// AddRPCMethod is a function to add a RPC method
func AddRPCMethod(method string, url string, description string, canCache bool) {
	newMethod := RPCMethod{}
	newMethod.Description = description
	newMethod.Enabled = true
	newMethod.CachingEnabled = true

	if conf, ok := interx.Config.RPC.API[method][url]; ok {
		newMethod.Enabled = !conf.Disable
		newMethod.CachingEnabled = !conf.CachingDisable
		newMethod.RateLimit = conf.RateLimit
		newMethod.AuthRateLimit = conf.AuthRateLimit
	}

	if !canCache {
		newMethod.CachingEnabled = false
	}

	if _, ok := rpcMethods[method]; !ok {
		rpcMethods[method] = map[string]RPCMethod{}
	}
	rpcMethods[method][url] = newMethod
}

// MakeGetRequest is a function to make GET request
func MakeGetRequest(rpcAddr string, url string, query string) (interface{}, interface{}, int) {
	resp, err := http.Get(fmt.Sprintf("%s%s?%s", rpcAddr, url, query))
	if err != nil {
		return ServeError(0, "", err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	response := new(RPCResponse)
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return nil, nil, resp.StatusCode
	}

	return response.Result, response.Error, resp.StatusCode
}

func makePostRequest(r *http.Request) (*RPCResponse, error) {
	resp, err := http.PostForm(fmt.Sprintf("%s%s", r.Host, r.URL), r.Form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := new(RPCResponse)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetInterxRequest is a function to get Interx Request
func GetInterxRequest(r *http.Request) InterxRequest {
	request := InterxRequest{}

	request.Method = r.Method
	request.Endpoint = fmt.Sprintf("%s", r.URL)
	request.Params, _ = ioutil.ReadAll(r.Body)

	return request
}

// GetHash is a function to get hash
func GetHash(request interface{}) string {
	// Calculate blake2b hash
	requestJSON, _ := json.Marshal(request)
	hash := blake2b.Sum256([]byte(requestJSON))
	return fmt.Sprintf("%X", hash)
}

// GetResponseFormat is a function to get response format
func GetResponseFormat(request InterxRequest, rpcAddr string) *common.ProxyResponse {
	response := new(common.ProxyResponse)
	response.Timestamp = time.Now().Unix()
	response.RequestHash = GetHash(request)
	response.Chainid = tasks.NodeStatus.Chainid
	response.Block = tasks.NodeStatus.Block
	response.Blocktime = tasks.NodeStatus.Blocktime

	return response
}

// GetResponseSignature is a function to get response signature
func GetResponseSignature(response common.ProxyResponse) (string, string) {
	// Get Response Hash
	responseHash := GetHash(response.Response)

	// Generate json to be signed
	sign := new(ResponseSign)
	sign.Chainid = response.Chainid
	sign.Block = response.Block
	sign.Blocktime = response.Blocktime
	sign.Timestamp = response.Timestamp
	sign.Response = responseHash
	signBytes, err := json.Marshal(sign)
	if err != nil {
		return "", responseHash
	}

	// Get Signature
	signature, err := interx.Config.PrivKey.Sign(signBytes)
	if err != nil {
		return "", responseHash
	}

	return base64.StdEncoding.EncodeToString([]byte(signature)), responseHash
}

// SearchCache is a function to search response in cache
func SearchCache(request InterxRequest, response *common.ProxyResponse) (bool, interface{}, interface{}, int) {
	fmt.Println("searching in the cache")

	chainIDHash := GetHash(response.Chainid)
	endpointHash := GetHash(request.Endpoint)
	requestHash := GetHash(request)

	result, err := GetCache(chainIDHash, endpointHash, requestHash)

	if err != nil {
		return false, nil, nil, -1
	}

	if result.ExpireAt.Before(time.Now()) && result.Response.Block != response.Block {
		return false, nil, nil, -1
	}

	return true, result.Response.Response, result.Response.Error, result.Status
}

// WrapResponse is a function to wrap response
func WrapResponse(w http.ResponseWriter, request InterxRequest, response common.ProxyResponse, statusCode int, saveToCashe bool) {
	if saveToCashe {
		fmt.Println("saving in the cache")
		chainIDHash := GetHash(response.Chainid)
		endpointHash := GetHash(request.Endpoint)
		requestHash := GetHash(request)

		err := PutCache(chainIDHash, endpointHash, requestHash, common.InterxResponse{
			Response: response,
			Status:   statusCode,
			ExpireAt: time.Now().Add(time.Duration(interx.Config.RPC.CachingDuration) * time.Second),
		})
		if err != nil {
			fmt.Printf("failed to save in the cache : %s\n", err.Error())
		}
		fmt.Println("save finished")
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Interx_chain_id", response.Chainid)
	w.Header().Add("Interx_block", strconv.FormatInt(response.Block, 10))
	w.Header().Add("Interx_blocktime", response.Blocktime)
	w.Header().Add("Interx_timestamp", strconv.FormatInt(response.Timestamp, 10))
	w.Header().Add("Interx_request_hash", response.RequestHash)

	if response.Response != nil {
		response.Signature, response.Hash = GetResponseSignature(response)

		w.Header().Add("Interx_signature", response.Signature)
		w.Header().Add("Interx_hash", response.Hash)
		w.WriteHeader(statusCode)

		json.NewEncoder(w).Encode(response.Response)
	} else {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response.Error)
	}
}

// ServeGRPC is a function to server GRPC
func ServeGRPC(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	recorder := httptest.NewRecorder()
	gwCosmosmux.ServeHTTP(recorder, r)
	resp := recorder.Result()

	result := new(interface{})
	if json.NewDecoder(resp.Body).Decode(result) == nil {
		if resp.StatusCode == http.StatusOK {
			return result, nil, resp.StatusCode
		}

		return nil, result, resp.StatusCode
	}

	return nil, nil, resp.StatusCode
}

// ServeError is a function to server GRPC
func ServeError(code int, data string, message string, statusCode int) (interface{}, interface{}, int) {
	return nil, common.ProxyResponseError{
		Code:    code,
		Data:    data,
		Message: message,
	}, statusCode
}

// GetAccountBalances is a function to get balances of an address
func GetAccountBalances(gwCosmosmux *runtime.ServeMux, r *http.Request, bech32addr string) []Coin {
	addr, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		return nil
	}

	r.URL.Path = fmt.Sprintf("/api/cosmos/bank/balances/%s", base64.URLEncoding.EncodeToString([]byte(addr)))
	r.URL.RawQuery = ""
	r.Method = "GET"

	recorder := httptest.NewRecorder()
	gwCosmosmux.ServeHTTP(recorder, r)
	resp := recorder.Result()

	type BalancesResponse struct {
		Balances []Coin `json:"balances"`
	}

	result := BalancesResponse{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.Balances
}

// GetAccountNumberSequence is a function to get AccountNumber and Sequence
func GetAccountNumberSequence(gwCosmosmux *runtime.ServeMux, r *http.Request, bech32addr string) (uint64, uint64) {
	addr, err := sdk.AccAddressFromBech32(bech32addr)
	if err != nil {
		return 0, 0
	}

	r.URL.Path = fmt.Sprintf("/api/cosmos/auth/accounts/%s", base64.URLEncoding.EncodeToString([]byte(addr)))
	r.URL.RawQuery = ""
	r.Method = "GET"

	recorder := httptest.NewRecorder()
	gwCosmosmux.ServeHTTP(recorder, r)
	resp := recorder.Result()

	type QueryAccountResponse struct {
		Account struct {
			Address       string `json:"addresss"`
			PubKey        string `json:"pubKey"`
			AccountNumber string `json:"accountNumber"`
			Sequence      string `json:"sequence"`
		} `json:"account"`
	}
	result := QueryAccountResponse{}
	json.NewDecoder(resp.Body).Decode(&result)

	accountNumber, _ := strconv.ParseInt(result.Account.AccountNumber, 10, 64)
	sequence, _ := strconv.ParseInt(result.Account.Sequence, 10, 64)

	return uint64(accountNumber), uint64(sequence)
}

// BroadcastTransaction is a function to post transaction, returns txHash
func BroadcastTransaction(rpcAddr string, txBytes []byte) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/broadcast_tx_commit?tx=0x%X", rpcAddr, txBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	type RPCTempResponse struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  struct {
			Hash string `json:"hash"`
		} `json:"result,omitempty"`
		Error struct {
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}

	result := new(RPCTempResponse)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(result.Error.Message)
	}

	return result.Result.Hash, nil
}
