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

	interx "github.com/KiraCore/sekai/INTERX/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/crypto/blake2b"
)

// AddRPCMethod is a function to add a RPC method
func AddRPCMethod(method string, url string, description string) {
	newMethod := RPCMethod{}
	newMethod.Description = description
	newMethod.Enabled = true

	if conf, ok := interx.WhitelistCg[method][url]; ok {
		newMethod.Enabled = !conf.Disable
		newMethod.RateLimit = conf.RateLimit
		newMethod.AuthRateLimit = conf.AuthRateLimit
	}

	if _, ok := rpcMethods[method]; !ok {
		rpcMethods[method] = map[string]RPCMethod{}
	}
	rpcMethods[method][url] = newMethod
}

// MakeGetRequest is a function to make GET request
func MakeGetRequest(rpcAddr string, url string, query string) (*RPCResponse, int, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s?%s", rpcAddr, url, query))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	result := new(RPCResponse)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return result, resp.StatusCode, nil
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

// GetResponseFormat is a function to get response format
func GetResponseFormat(request InterxRequest, rpcAddr string) *ProxyResponse {
	response := new(ProxyResponse)
	response.Timestamp = time.Now().Unix()

	requestJSON, _ := json.Marshal(request)
	hash := blake2b.Sum256([]byte(requestJSON))
	response.RequestHash = fmt.Sprintf("%X", hash)

	resp, err := http.Get(fmt.Sprintf("%s/block", rpcAddr))
	if err != nil {
		return response
	}
	defer resp.Body.Close()

	type RPCTempResponse struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  struct {
			Block struct {
				Header struct {
					Chainid string `json:"chain_id"`
					Height  string `json:"height"`
					Time    string `json:"time"`
				} `json:"header"`
			} `json:"block"`
		} `json:"result"`
		Error interface{} `json:"error"`
	}

	result := new(RPCTempResponse)
	if json.NewDecoder(resp.Body).Decode(result) != nil {
		return response
	}

	response.Chainid = result.Result.Block.Header.Chainid
	response.Block, _ = strconv.ParseInt(result.Result.Block.Header.Height, 10, 64)
	response.Blocktime = result.Result.Block.Header.Time

	return response
}

// GetResponseSignature is a function to get response signature
func GetResponseSignature(response ProxyResponse) (string, string) {
	// Calculate blake2b hash
	responseJSON, err := json.Marshal(response.Response)
	if err != nil {
		return "", ""
	}
	hash := blake2b.Sum256([]byte(responseJSON))
	responseHash := fmt.Sprintf("%X", hash)

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
	signature, err := interx.InterxCg.PrivKey.Sign(signBytes)
	if err != nil {
		return "", responseHash
	}

	return base64.StdEncoding.EncodeToString([]byte(signature)), responseHash
}

// WrapResponse is a function to wrap response
func WrapResponse(w http.ResponseWriter, response ProxyResponse, statusCode int) {
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
func ServeGRPC(w http.ResponseWriter, r *http.Request, gwCosmosmux *runtime.ServeMux, request InterxRequest, rpcAddr string) {
	recorder := httptest.NewRecorder()
	gwCosmosmux.ServeHTTP(recorder, r)
	resp := recorder.Result()

	response := GetResponseFormat(request, rpcAddr)

	result := new(interface{})
	if json.NewDecoder(resp.Body).Decode(result) == nil {
		if resp.StatusCode == 200 {
			response.Response = result
		} else {
			response.Error = result
		}
	}

	WrapResponse(w, *response, resp.StatusCode)
}

// ServeRPC is a function to server RPC
func ServeRPC(w http.ResponseWriter, request InterxRequest, result *RPCResponse, rpcAddr string, statusCode int) {
	response := GetResponseFormat(request, rpcAddr)
	response.Response = result.Result
	response.Error = result.Error

	WrapResponse(w, *response, statusCode)
}

// ServeError is a function to server GRPC
func ServeError(w http.ResponseWriter, request InterxRequest, rpcAddr string, code int, data string, message string, statusCode int) {
	response := GetResponseFormat(request, rpcAddr)

	response.Response = nil
	response.Error = ProxyResponseError{
		Code:    code,
		Data:    data,
		Message: message,
	}

	WrapResponse(w, *response, statusCode)
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

	if resp.StatusCode != 200 {
		return "", errors.New(result.Error.Message)
	}

	return result.Result.Hash, nil
}

// GetChainID is a function to get ChainID
func GetChainID(rpcAddr string) string {
	r, err := http.Get(fmt.Sprintf("%s/block", rpcAddr))
	if err != nil {
		return ""
	}
	defer r.Body.Close()

	type RPCTempResponse struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  struct {
			Block struct {
				Header struct {
					Chainid string `json:"chain_id"`
				} `json:"header"`
			} `json:"block"`
		} `json:"result"`
		Error interface{} `json:"error"`
	}

	result := new(RPCTempResponse)
	if json.NewDecoder(r.Body).Decode(result) != nil {
		return ""
	}

	return result.Result.Block.Header.Chainid
}
