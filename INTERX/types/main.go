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

// TxAmount is a struct to be used for query transaction response
type TxAmount struct {
	Amount int64  `json:"amount,omitempty"`
	Denom  string `json:"denom,omitempty"`
}

// Transaction is a struct to be used for query transaction response
type Transaction struct {
	Type    string     `json:"type,omitempty"`
	From    string     `json:"from,omitemtpy"`
	To      string     `json:"to,omitempty"`
	Amounts []TxAmount `json:"amounts,omitempty"`
}

// TransactionResult is a struct to be used for query transaction response
type TransactionResult struct {
	Hash           string        `json:"hash"`
	Status         string        `json:"status"`
	BlockHeight    int64         `json:"block_height"`
	BlockTimestamp int64         `json:"block_timestamp"`
	Confirmation   int64         `json:"confirmation"`
	Transactions   []Transaction `json:"transactions"`
	Fees           []TxAmount    `json:"fees"`
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

type QueryValidator struct {
	Address    string `protobuf:"bytes,1,opt,name=address,proto3" json:"address"`
	Valkey     string `protobuf:"bytes,2,opt,name=valkey,proto3" json:"valkey"`
	Pubkey     string `protobuf:"bytes,3,opt,name=pubkey,proto3" json:"pubkey"`
	Proposer   string `protobuf:"bytes,4,opt,name=proposer,proto3" json:"proposer"`
	Moniker    string `protobuf:"bytes,5,opt,name=moniker,proto3" json:"moniker"`
	Website    string `protobuf:"bytes,6,opt,name=website,proto3" json:"website"`
	Social     string `protobuf:"bytes,7,opt,name=social,proto3" json:"social"`
	Identity   string `protobuf:"bytes,8,opt,name=identity,proto3" json:"identity"`
	Commission string `protobuf:"bytes,9,opt,name=commission,proto3" json:"commission"`
	Status     string `protobuf:"bytes,10,opt,name=status,proto3" json:"status"`
	Rank       int64  `protobuf:"varint,11,opt,name=rank,proto3" json:"rank"`
	Streak     int64  `protobuf:"varint,12,opt,name=streak,proto3" json:"streak"`
	Mischance  int64  `protobuf:"varint,13,opt,name=mischance,proto3" json:"mischance"`
}

type AllValidators struct {
	Status struct {
		ConsensusStopped   bool `json:"consensus_stopped"`
		ActiveValidators   int  `json:"active_validators"`
		PausedValidators   int  `json:"paused_validators"`
		InactiveValidators int  `json:"inactive_validators"`
		JailedValidators   int  `json:"jailed_validators"`
		TotalValidators    int  `json:"total_validators"`
		WaitingValidators  int  `json:"waiting_validators"`
	} `json:"status"`
	Waiting    []string         `json:"waiting"`
	Validators []QueryValidator `json:"validators"`
}

const (
	// GET is a constant to refer GET HTTP Method
	GET string = "GET"
	// POST is a constant to refer POST HTTP Method
	POST string = "POST"
)
