package rosetta

type CurrencyMetadata struct {
	Issuer string `json:"Issuer,omitempty"`
}

type AmountMetadata struct {
}

type OperationMetadata struct {
	Asm string `json:"asm,omitempty"`
	Hex string `json:"hex,omitempty"`
}

type NetworkMetadata struct {
	Producer string `json:"producer,omitempty"`
}

type AccountMetadata struct {
}

type TransactionMetadata struct {
	Size     int64 `json:"size,omitempty"`
	LockTime int64 `json:"lockTime,omitempty"`
}

type BlockMetadata struct {
	TransactionsRoot string `json:"transactions_root,omitempty"`
	Difficulty       string `json:"difficulty,omitempty"`
}
