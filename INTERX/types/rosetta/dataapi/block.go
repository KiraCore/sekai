package dataapi

import "github.com/KiraCore/sekai/INTERX/types/rosetta"

type BlockRequest struct {
	NetworkIdentifier rosetta.NetworkIdentifier      `json:"network_identifier"`
	BlockIdentifier   rosetta.PartialBlockIdentifier `json:"block_identifier"`
}

type BlockResponse struct {
	Block             rosetta.Block                   `json:"block,omitempty"`
	OtherTransactions []rosetta.TransactionIdentifier `json:"other_transactions,omitempty"`
}

type BlockTransactionRequest struct {
	NetworkIdentifier     rosetta.NetworkIdentifier     `json:"network_identifier"`
	BlockIdentifier       rosetta.BlockIdentifier       `json:"block_identifier"`
	TransactionIdentifier rosetta.TransactionIdentifier `json:"transaction_identifier"`
}

type BlockTransactionResponse struct {
	Transaction rosetta.Transaction `json:"transaction"`
}
