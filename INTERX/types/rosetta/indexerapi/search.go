package indexerapi

import "github.com/KiraCore/sekai/INTERX/types/rosetta"

type SearchTransactionsRequest struct {
	NetworkIdentifier     rosetta.NetworkIdentifier     `json:"network_identifier"`
	Operator              rosetta.Operator              `json:"operator,omitempty"`
	MaxBlock              int64                         `json:"max_block,omitempty"`
	Offset                int64                         `json:"offset,omitempty"`
	Limit                 int64                         `json:"limit,omitempty"`
	TransactionIdentifier rosetta.TransactionIdentifier `json:"transaction_identifier,omitempty"`
	AccountIdentifier     rosetta.AccountIdentifier     `json:"account_identifier,omitempty"`
	CoinIdentifier        rosetta.CoinIdentifier        `json:"coin_identifier,omitempty"`
	Currency              rosetta.Currency              `json:"currency,omitempty"`
	Status                string                        `json:"status,omitempty"`
	Type                  string                        `json:"type,omitempty"`
	Address               string                        `json:"address,omitempty"`
	Success               bool                          `json:"success,omitempty"`
}

type SearchTransactionsResponse struct {
	Transactions []rosetta.BlockTransactions `json:"transactions"`
	TotalCount   int64                       `json:"total_count"`
	NextOffset   int64                       `json:"next_offset,omitempty"`
}
