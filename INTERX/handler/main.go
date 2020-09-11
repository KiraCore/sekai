package interx

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"golang.org/x/crypto/blake2b"
)

const (
	// QueryStatus is an endpoint for query status
	QueryStatus string = "/api/cosmos/status"
	// QueryTransactionHash is an endpoint for query transaction hash
	QueryTransactionHash string = "/api/cosmos/tx/"
	// PostTransaction is an endpoint for post transaction
	PostTransaction string = "/api/cosmos/tx"
	// QueryBalances is an endpoint for query balances
	QueryBalances string = "/api/cosmos/bank/balances/"
	// QuerySupply is an endpoint for query total supply
	QuerySupply string = "/api/cosmos/bank/supply"
)

// Endpoints is an array of Proxy endpoints
var Endpoints = []string{
	QueryStatus,
	QueryTransactionHash,
	PostTransaction,
	QueryBalances,
	QuerySupply,
}

// RPCResponse is a struct of RPC response
type RPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}

// ProxyResponse is a struct to be used for proxy response
type ProxyResponse struct {
	Chainid   string      `json:"chain_id"`
	Block     int64       `json:"block"`
	Blocktime string      `json:"block_time"`
	Timestamp int64       `json:"timestamp"`
	Response  interface{} `json:"response,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Signature string      `json:"signature,omitempty"`
	Hash      string      `json:"hash,omitempty"`
}

// ResponseSign is a struct to be used for response sign
type ResponseSign struct {
	Chainid   string `json:"chain_id"`
	Block     int64  `json:"block"`
	Blocktime string `json:"block_time"`
	Timestamp int64  `json:"timestamp"`
	Response  string `json:"response"`
}

// CopyHeader is a function to copy http header
func CopyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			if k != "Content-Length" {
				dst.Add(k, v)
			}
		}
	}
}

// MakeGetRequest is a function to make GET request
func MakeGetRequest(w http.ResponseWriter, r *http.Request) (*RPCResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", r.Host, r.URL))
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}
	defer resp.Body.Close()

	CopyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	result := new(RPCResponse)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}

	return result, nil
}

func makePostRequest(w http.ResponseWriter, r *http.Request) (*RPCResponse, error) {
	resp, err := http.PostForm(fmt.Sprintf("%s%s", r.Host, r.URL), r.Form)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	defer resp.Body.Close()

	CopyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	result := new(RPCResponse)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
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

// GetResponseSignature is a function to get response signature
func GetResponseSignature(response ProxyResponse) (string, string) {
	// Calculate blake2b hash
	responseJSON, err := json.Marshal(response.Response)
	if err != nil {
		return "", ""
	}
	hash := blake2b.Sum256([]byte(responseJSON))
	responseHash := "0x" + hex.EncodeToString(hash[:])

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

	// Generate PrivKey
	privKey := GenEd25519PrivKey()

	// Get Signature
	signature, err := privKey.Sign(signBytes)
	if err != nil {
		return "", responseHash
	}

	return base64.StdEncoding.EncodeToString([]byte(signature)), responseHash
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

// WrapResponse is a function to wrap response
func WrapResponse(w http.ResponseWriter, response ProxyResponse) {
	if response.Response != nil {
		response.Signature, response.Hash = GetResponseSignature(response)
	}

	json.NewEncoder(w).Encode(response)
}
