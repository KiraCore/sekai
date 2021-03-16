package rosetta

type SubNetworkIdentifier struct {
	Network  string      `json:"network"`
	Metadata interface{} `json:"metadata,omitempty"`
}

type NetworkIdentifier struct {
	Blockchain           string                `json:"blockchain"`
	Network              string                `json:"network"`
	SubNetworkIdentifier *SubNetworkIdentifier `json:"sub_network_identifier,omitempty"`
}

type BlockIdentifier struct {
	Index int64  `json:"index"`
	Hash  string `json:"hash"`
}

type PartialBlockIdentifier struct {
	Index int64  `json:"index,omitempty"`
	Hash  string `json:"hash,omitempty"`
}

type TransactionIdentifier struct {
	Hash string `json:"hash"`
}

type OperationIdentifier struct {
	Index        int64 `json:"index"`
	NetworkIndex int64 `json:"network_index,omitempty"`
}

type SubAccountIdentifier struct {
	Address  string          `json:"address"`
	Metadata AccountMetadata `json:"metadata,omitempty"`
}

type AccountIdentifier struct {
	Address    string               `json:"address"`
	SubAccount SubAccountIdentifier `json:"sub_account,omitempty"`
	Metadata   AccountMetadata      `json:"metadata,omitempty"`
}

type CoinIdentifier struct {
	Identifier string `json:"identifier"`
}
