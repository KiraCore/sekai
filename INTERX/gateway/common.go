package gateway

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
func MakeGetRequest(w http.ResponseWriter, rpcAddr string, url string, query string) (*RPCResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s?%s", rpcAddr, url, query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	result := new(RPCResponse)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func makePostRequest(w http.ResponseWriter, r *http.Request) (*RPCResponse, error) {
	resp, err := http.PostForm(fmt.Sprintf("%s%s", r.Host, r.URL), r.Form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	result := new(RPCResponse)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetResponseFormat is a function to get response format
func GetResponseFormat(rpcAddr string) *ProxyResponse {
	response := new(ProxyResponse)
	response.Timestamp = time.Now().Unix()

	r, err := http.Get(fmt.Sprintf("%s/block", rpcAddr))
	if err != nil {
		return response
	}
	defer r.Body.Close()

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
	if json.NewDecoder(r.Body).Decode(result) != nil {
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
	signature, err := interx.PrivKey.Sign(signBytes)
	if err != nil {
		return "", responseHash
	}

	return base64.StdEncoding.EncodeToString([]byte(signature)), responseHash
}

// WrapResponse is a function to wrap response
func WrapResponse(w http.ResponseWriter, response ProxyResponse) {
	if response.Response != nil {
		response.Signature, response.Hash = GetResponseSignature(response)
	}

	json.NewEncoder(w).Encode(response)
}

// ServeGRPC is a function to server GRPC
func ServeGRPC(w http.ResponseWriter, r *http.Request, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	recorder := httptest.NewRecorder()
	gwCosmosmux.ServeHTTP(recorder, r)
	resp := recorder.Result()

	response := GetResponseFormat(rpcAddr)

	result := new(interface{})
	if json.NewDecoder(resp.Body).Decode(result) == nil {
		if resp.StatusCode == 200 {
			response.Response = result
		} else {
			response.Error = result
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	WrapResponse(w, *response)
}

// ServeRPC is a function to server RPC
func ServeRPC(w http.ResponseWriter, result *RPCResponse, rpcAddr string) {
	response := GetResponseFormat(rpcAddr)
	response.Response = result.Result
	response.Error = result.Error

	WrapResponse(w, *response)
}

// ServeError is a function to server GRPC
func ServeError(w http.ResponseWriter, rpcAddr string, code int, data string, message string, statusCode int) {
	response := GetResponseFormat(rpcAddr)

	response.Response = nil
	response.Error = ProxyResponseError{
		Code:    code,
		Data:    data,
		Message: message,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	WrapResponse(w, *response)
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
