package dataapi

import "github.com/KiraCore/sekai/INTERX/types/rosetta"

type MempoolRequest NetworkRequest

type MempoolResponse struct {
	TransactionIdentifiers []rosetta.TransactionIdentifier `json:"transaction_identifiers"`
}

type MempoolTransactionRequest struct {
	NetworkIdentifier     rosetta.NetworkIdentifier     `json:"network_identifier"`
	TransactionIdentifier rosetta.TransactionIdentifier `json:"transaction_identifier"`
}

type MempoolTransactionResponse struct {
	Transaction rosetta.Transaction `json:"transaction"`
	Metadata    interface{}         `json:"metadata,omitempty"`
}
