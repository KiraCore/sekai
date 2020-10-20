package gateway

import "time"

// RPCResponse is a struct of RPC response
type RPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}

// ProxyResponseError is a struct to be used for proxy response error
type ProxyResponseError struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

// ProxyResponse is a struct to be used for proxy response
type ProxyResponse struct {
	Chainid     string      `json:"chain_id"`
	Block       int64       `json:"block"`
	Blocktime   string      `json:"block_time"`
	Timestamp   int64       `json:"timestamp"`
	Response    interface{} `json:"response,omitempty"`
	Error       interface{} `json:"error,omitempty"`
	Signature   string      `json:"signature,omitempty"`
	Hash        string      `json:"hash,omitempty"`
	RequestHash string      `json:"request_hash,omitempty"`
}

// ResponseSign is a struct to be used for response sign
type ResponseSign struct {
	Chainid   string `json:"chain_id"`
	Block     int64  `json:"block"`
	Blocktime string `json:"block_time"`
	Timestamp int64  `json:"timestamp"`
	Response  string `json:"response"`
}

// Coin is a struct for coin
type Coin struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

// FaucetAccountInfo is a struct to be used for Faucet Account Info
type FaucetAccountInfo struct {
	Address  string `json:"address"`
	Balances []Coin `json:"balances"`
}

// RPCMethod is a struct to be used for rpc_methods API
type RPCMethod struct {
	Description    string  `json:"description"`
	Enabled        bool    `json:"enabled"`
	RateLimit      float64 `json:"rate_limit,omitempty"`
	AuthRateLimit  float64 `json:"auth_rate_limit,omitempty"`
	CachingEnabled bool    `json:"caching_enabled"`
}

// InterxRequest is a struct to be used for request hash
type InterxRequest struct {
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
	Params   []byte `json:"params"`
}

// InterxResponse is a struct to be used for response caching
type InterxResponse struct {
	Response ProxyResponse `json:"response"`
	Status   int           `json:"status"`
	ExpireAt time.Time     `json:"expire_at"`
}

const (
	// GET is a constant to refer GET HTTP Method
	GET string = "GET"
	// POST is a constant to refer POST HTTP Method
	POST string = "POST"
)

var rpcMethods = make(map[string]map[string]RPCMethod)
var interxChainID string = ""
