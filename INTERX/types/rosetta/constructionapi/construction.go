package constructionapi

import "github.com/KiraCore/sekai/INTERX/types/rosetta"

type ConstructionCombineRequest struct {
	NetworkIdentifier   rosetta.NetworkIdentifier `json:"network_identifier"`
	UnsignedTransaction string                    `json:"unsigned_transaction"`
	Signatures          []rosetta.Signature       `json:"signatures"`
}

type ConstructionCombineResponse struct {
	SignedTransaction string `json:"signed_transaction"`
}

type ConstructionDeriveRequest struct {
	NetworkIdentifier rosetta.NetworkIdentifier `json:"network_identifier"`
	PublicKey         rosetta.PublicKey         `json:"public_key"`
	Metadata          interface{}               `json:"metadata,omitempty"`
}

type ConstructionDeriveResponse struct {
	Address           string                    `json:"address,omitempty"`
	AccountIdentifier rosetta.AccountIdentifier `json:"account_identifier,omitempty"`
	Metadata          interface{}               `json:"metadata,omitempty"`
}

type TransactionIdentifierResponse struct {
	TransactionIdentifier rosetta.TransactionIdentifier `json:"transaction_identifier"`
	Metadata              interface{}                   `json:"metadata,omitempty"`
}

type ConstructionHashRequest struct {
	NetworkIdentifier rosetta.NetworkIdentifier `json:"network_identifier"`
	SignedTransaction string                    `json:"signed_transaction"`
}

type ConstructionHashResponse TransactionIdentifierResponse

type ConstructionMetadataRequest struct {
	NetworkIdentifier rosetta.NetworkIdentifier `json:"network_identifier"`
	Options           interface{}               `json:"options,omitempty"`
	PublicKeys        []rosetta.PublicKey       `json:"public_keys,omitempty"`
}

type ConstructionMetadataResponse struct {
	Metadata     interface{}      `json:"metadata"`
	SuggestedFee []rosetta.Amount `json:"suggested_fee,omitempty"`
}

type ConstructionParseRequest struct {
	NetworkIdentifier rosetta.NetworkIdentifier `json:"network_identifier"`
	Signed            bool                      `json:"signed"`
	Tranaction        string                    `json:"transaction"`
}

type ConstructionParseResponse struct {
	Operations               []rosetta.Operation         `json:"operations"`
	Signers                  []string                    `json:"signers,omitempty"`
	AccountIdentifierSigners []rosetta.AccountIdentifier `json:"account_identifier_signers,omitempty"`
	Metadata                 interface{}                 `json:"metadata,omitempty"`
}

type ConstructionPayloadsRequest struct {
	NetworkIdentifier rosetta.NetworkIdentifier `json:"network_identifier"`
	Operations        []rosetta.Operation       `json:"operations"`
	Metadata          interface{}               `json:"metadata,omitempty"`
	PublicKeys        []rosetta.PublicKey       `json:"public_keys,omitempty"`
}

type ConstructionPayloadsResponse struct {
	UnsignedTransaction string                   `json:"unsigned_transaction"`
	Payloads            []rosetta.SigningPayload `json:"payloads"`
}

type ConstructionPreprocessRequest struct {
	NetworkIdentifier      rosetta.NetworkIdentifier `json:"network_identifier"`
	Operations             []rosetta.Operation       `json:"operations"`
	Metadata               interface{}               `json:"metadata,omitempty"`
	MaxFee                 []rosetta.Amount          `json:"max_fee,omitempty"`
	SuggestedFeeMultiplier float64                   `json:"suggested_fee_multiplier,omitempty"`
}

type ConstructionPreprocessResponse struct {
	Options            interface{}                 `json:"options,omitempty"`
	RequiredPublicKeys []rosetta.AccountIdentifier `json:"required_public_keys,omitempty"`
}

type ConstructionSubmitRequest struct {
	NetworkIdentifier rosetta.NetworkIdentifier `json:"network_identifier"`
	SignedTransaction string                    `json:"signed_transaction"`
}

type ConstructionSubmitResponse TransactionIdentifierResponse
