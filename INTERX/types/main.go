package types

import (
	"time"
)

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

// ProxyResponseError is a struct to be used for proxy response error
type ProxyResponseError struct {
	Code    int    `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

// InterxResponse is a struct to be used for response caching
type InterxResponse struct {
	Response ProxyResponse `json:"response"`
	Status   int           `json:"status"`
	ExpireAt time.Time     `json:"expire_at"`
}

// DataReferenceEntry is a struct to be used for data reference
type DataReferenceEntry struct {
	Hash      string `json:"hash"`
	Reference string `json:"reference"`
	Encoding  string `json:"encoding"`
	Size      uint64 `json:"size"`
}

// RPCMethod is a struct to be used for rpc_methods API
type RPCMethod struct {
	Description     string  `json:"description"`
	Enabled         bool    `json:"enabled"`
	RateLimit       float64 `json:"rate_limit,omitempty"`
	AuthRateLimit   float64 `json:"auth_rate_limit,omitempty"`
	CachingEnabled  bool    `json:"caching_enabled"`
	CachingDuration int64   `json:"caching_duration,omitempty"`
}

// RPCResponse is a struct of RPC response
type RPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// ResponseSign is a struct to be used for response sign
type ResponseSign struct {
	Chainid   string `json:"chain_id"`
	Block     int64  `json:"block"`
	Blocktime string `json:"block_time"`
	Timestamp int64  `json:"timestamp"`
	Response  string `json:"response"`
}

// DepositWithdrawTransaction is a struct to be used for deposit/withdraw transaction
type DepositWithdrawTransaction struct {
	Address string `json:"address"`
	Type    string `json:"type"`
	Denom   string `json:"denom,omitempty"`
	Amount  int64  `json:"amount"`
}

// DepositWithdrawResult is a struct to be used for query deposit/withdraw transaction response
type DepositWithdrawResult struct {
	Time int64                        `json:"time"`
	Txs  []DepositWithdrawTransaction `json:"txs"`
}

// Transaction is a struct to be used for query transaction response
type Transaction struct {
	From   string `json:"from,omitempty"`
	To     string `json:"to,omitempty"`
	Amount int64  `json:"amount,omitempty"`
	Denom  string `json:"denom,omitempty"`
}

// TransactionResult is a struct to be used for query transaction response
type TransactionResult struct {
	Hash           string        `json:"hash"`
	Status         string        `json:"status"`
	BlockHeight    int64         `json:"block_height"`
	BlockTimestamp int64         `json:"block_timestamp"`
	Confirmation   int64         `json:"confirmation"`
	Transactions   []Transaction `json:"transactions"`
	Fees           []Transaction `json:"fees"`
	GasWanted      int64         `json:"gas_wanted"`
	GasUsed        int64         `json:"gas_used"`
}

// TransactionSearchResult is a struct to be used for query transaction response
type TransactionSearchResult struct {
	Txs        []TransactionResult `json:"txs"`
	TotalCount int                 `json:"total_count"`
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

// InterxRequest is a struct to be used for request hash
type InterxRequest struct {
	Method   string `json:"method"`
	Endpoint string `json:"endpoint"`
	Params   []byte `json:"params"`
}

const (
	// GET is a constant to refer GET HTTP Method
	GET string = "GET"
	// POST is a constant to refer POST HTTP Method
	POST string = "POST"
)
