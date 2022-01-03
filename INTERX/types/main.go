package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	txSinging "github.com/cosmos/cosmos-sdk/types/tx/signing"
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
	Response             ProxyResponse `json:"response"`
	Status               int           `json:"status"`
	CacheTime            time.Time     `json:"cache_time"`
	CachingDuration      int64         `json:"caching_duration"`
	CachingBlockDuration int64         `json:"caching_block_duration"`
}

// DataReferenceEntry is a struct to be used for data reference
type DataReferenceEntry struct {
	Hash      string `json:"hash"`
	Reference string `json:"reference"`
	Encoding  string `json:"encoding"`
	Size      uint64 `json:"size,string"`
}

// RPCMethod is a struct to be used for rpc_methods API
type RPCMethod struct {
	Description          string  `json:"description"`
	Enabled              bool    `json:"enabled"`
	RateLimit            float64 `json:"rate_limit,omitempty"`
	AuthRateLimit        float64 `json:"auth_rate_limit,omitempty"`
	CachingEnabled       bool    `json:"caching_enabled"`
	CachingDuration      int64   `json:"caching_duration"`
	CachingBlockDuration int64   `json:"caching_block_duration"`
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
	Type    string     `json:"type,omitempty"`
	From    string     `json:"from,omitempty"`
	To      string     `json:"to,omitempty"`
	Amounts []sdk.Coin `json:"amounts,omitempty"`
}

type TxMsg struct {
	Type string  `json:"type"`
	Data sdk.Msg `json:"data"`
}

// TransactionResult is a struct to be used for query transaction response
type TransactionResult struct {
	Hash           string        `json:"hash"`
	Status         string        `json:"status"`
	BlockHeight    int64         `json:"block_height"`
	BlockTimestamp int64         `json:"block_timestamp"`
	Confirmation   int64         `json:"confirmation"`
	Msgs           []TxMsg       `json:"msgs"`
	Transactions   []Transaction `json:"transactions"`
	Fees           []sdk.Coin    `json:"fees"`
	GasWanted      int64         `json:"gas_wanted"`
	GasUsed        int64         `json:"gas_used"`
	Memo           string        `json:"memo"`
}

type TransactionUnconfirmedResult struct {
	Msgs      []TxMsg                 `json:"msgs"`
	Fees      []sdk.Coin              `json:"fees"`
	Gas       uint64                  `json:"gas"`
	Signature []txSinging.SignatureV2 `json:"signature"`
	Memo      string                  `json:"memo"`
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

type TokenAlias struct {
	Decimals int64    `json:"decimals"`
	Denoms   []string `json:"denoms"`
	Icon     string   `json:"icon"`
	Name     string   `json:"name"`
	Symbol   string   `json:"symbol"`
}

type TokenSupply struct {
	Amount int64  `json:"amount,string"`
	Denom  string `json:"denom"`
}

// ID is a field for facuet claim struct.
func (c TokenAlias) ID() (jsonField string, value interface{}) {
	value = c.Symbol
	jsonField = "height"
	return
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

type IdentityRecord struct {
	ID        uint64   `json:"id,string"`
	Key       string   `json:"key"`
	Value     string   `json:"value"`
	Date      string   `json:"date"`
	Verifiers []string `json:"verifiers"`
}

type QueryValidator struct {
	Top int `json:"top,string"`

	Address             string           `json:"address"`
	Valkey              string           `json:"valkey"`
	Pubkey              string           `json:"pubkey"`
	Proposer            string           `json:"proposer"`
	Moniker             string           `json:"moniker"`
	Status              string           `json:"status"`
	Rank                int64            `json:"rank,string"`
	Streak              int64            `json:"streak,string"`
	Mischance           int64            `json:"mischance,string"`
	MischanceConfidence int64            `json:"mischance_confidence,string"`
	Identity            []IdentityRecord `json:"identity,omitempty"`

	// Additional
	StartHeight           int64  `json:"start_height,string"`
	InactiveUntil         string `json:"inactive_until"`
	LastPresentBlock      int64  `json:"last_present_block,string"`
	MissedBlocksCounter   int64  `json:"missed_blocks_counter,string"`
	ProducedBlocksCounter int64  `json:"produced_blocks_counter,string"`

	// From Identity Records
	Description       string `json:"description,omitempty"`
	Website           string `json:"website,omitempty"`
	Logo              string `json:"logo,omitempty"`
	Social            string `json:"social,omitempty"`
	Contact           string `json:"contact,omitempty"`
	Validator_node_id string `json:"validator_node_id,omitempty"`
	Sentry_node_id    string `json:"sentry_node_id,omitempty"`
}

type QueryValidators []QueryValidator

func (s QueryValidators) Len() int {
	return len(s)
}
func (s QueryValidators) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s QueryValidators) Less(i, j int) bool {
	if s[i].Status != s[j].Status {
		if s[j].Status == "ACTIVE" {
			return false
		}
		if s[i].Status == "ACTIVE" {
			return true
		}
		return s[i].Status > s[j].Status
	}
	if s[i].Rank != s[j].Rank {
		return s[i].Rank > s[j].Rank
	}
	if s[i].Streak != s[j].Streak {
		return s[i].Streak > s[j].Streak
	}
	if s[i].MissedBlocksCounter != s[j].MissedBlocksCounter {
		return s[i].MissedBlocksCounter < s[j].MissedBlocksCounter
	}

	return false
}

type AllValidators struct {
	Status struct {
		ActiveValidators   int `json:"active_validators"`
		PausedValidators   int `json:"paused_validators"`
		InactiveValidators int `json:"inactive_validators"`
		JailedValidators   int `json:"jailed_validators"`
		TotalValidators    int `json:"total_validators"`
		WaitingValidators  int `json:"waiting_validators"`
	} `json:"status"`
	Waiting    []string         `json:"waiting"`
	Validators []QueryValidator `json:"validators"`
}

type ValidatorSigningInfo struct {
	Address               string `json:"address"`
	StartHeight           int64  `json:"start_height,string"`
	InactiveUntil         string `json:"inactive_until"`
	MischanceConfidence   int64  `json:"mischance_confidence,string"`
	Mischance             int64  `json:"mischance,string"`
	LastPresentBlock      int64  `json:"last_present_block,string"`
	MissedBlocksCounter   int64  `json:"missed_blocks_counter,string"`
	ProducedBlocksCounter int64  `json:"produced_blocks_counter,string"`
}

const (
	// GET is a constant to refer GET HTTP Method
	GET string = "GET"
	// POST is a constant to refer POST HTTP Method
	POST string = "POST"
)
