package gateway

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"golang.org/x/crypto/blake2b"
)

func getAPIConfig(url string, method string) (*APIConfig, error) {
	api := Endpoint{}
	api.URL = url
	api.Method = method

	i := sort.Search(len(config), func(i int) bool {
		if api.URL != config[i].API.URL {
			return api.URL < config[i].API.URL
		}
		return api.Method <= config[i].API.Method
	})

	if i < len(config) && config[i].API.URL == api.URL && config[i].API.Method == api.Method {
		conf := new(APIConfig)
		conf.API = config[i].API
		conf.Disable = config[i].Disable
		conf.RateLimit = config[i].RateLimit
		conf.AuthRateLimit = config[i].AuthRateLimit

		return conf, nil
	}

	return nil, errors.New("Not Found")
}

// AddRPCMethod is a function to add a RPC method
func AddRPCMethod(name string, url string, method string) {
	newMethod := RPCMethod{}
	newMethod.API = Endpoint{}
	newMethod.API.URL = url
	newMethod.API.Method = method
	newMethod.Enabled = true

	conf, err := getAPIConfig(url, method)
	if err == nil {
		newMethod.Enabled = !conf.Disable
		newMethod.RateLimit = conf.RateLimit
		newMethod.AuthRateLimit = conf.AuthRateLimit
	}

	rpcMethods[name] = newMethod
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

// GenSecp256k1PrivKey is a function to generate a pubKey
func GenSecp256k1PrivKey() secp256k1.PrivKey {
	return secp256k1.GenPrivKey()
}

// GenEd25519PrivKey is a function to generate a pubKey
func GenEd25519PrivKey() ed25519.PrivKey {
	return ed25519.GenPrivKey()
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
	signature, err := privKey.Sign(signBytes)
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
