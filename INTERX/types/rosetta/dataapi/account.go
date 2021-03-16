package dataapi

import "github.com/KiraCore/sekai/INTERX/types/rosetta"

type AccountBalanceRequest struct {
	NetworkIdentifier rosetta.NetworkIdentifier      `json:"network_identifier"`
	AccountIdentifier rosetta.AccountIdentifier      `json:"account_identifier"`
	BlockIdentifier   rosetta.PartialBlockIdentifier `json:"block_identifier,omitempty"`
	Currencies        []rosetta.Currency             `json:"currencies,omitempty"`
}

type AccountBalanceResponse struct {
	BlockIdentifier rosetta.BlockIdentifier `json:"block_identifier"`
	Balances        []rosetta.Amount        `json:"balances"`
	Metadata        interface{}             `json:"metadata,omitempty"`
}

type AccountCoinsRequest struct {
	NetworkIdentifier rosetta.NetworkIdentifier `json:"network_identifier"`
	AccountIdentifier rosetta.AccountIdentifier `json:"account_identifier"`
	IncludeMempool    bool                      `json:"include_mempool"`
	Currencies        []rosetta.Currency        `json:"currencies,omitempty"`
}

type AccountCoinsResponse struct {
	BlockIdentifier rosetta.BlockIdentifier `json:"block_identifier"`
	Coins           []rosetta.Coin          `json:"coins"`
	Metadata        interface{}             `json:"metadata,omitempty"`
}
