package rosetta

type OperationStatus struct {
	Status     string `json:"status"`
	Successful bool   `json:"successful"`
}

type BalanceExemption struct {
	SubAccountAddress string        `json:"sub_account_address,omitempty"`
	Currency          Currency      `json:"currency,omitempty"`
	ExemptionType     ExemptionType `json:"exemption_type,omitempty"`
}

type Allow struct {
	OperationStatuses       []OperationStatus  `json:"operation_statuses"`
	OperationTypes          []string           `json:"operation_types"`
	Errors                  []Error            `json:"errors"`
	HistoricalBalanceLookup bool               `json:"historical_balance_lookup"`
	TimestampStartIndex     int64              `json:"timestamp_start_index,omitempty"`
	CallMethods             []string           `json:"call_methods"`
	BalanceExemptions       []BalanceExemption `json:"balance_exemptions"`
	MempoolCoins            bool               `json:"mempool_coins"`
}

type Currency struct {
	Symbol   string           `json:"symbol"`
	Decimals int64            `json:"decimals"`
	Metadata CurrencyMetadata `json:"metadata,omitempty"`
}

type Amount struct {
	Value    string         `json:"value"`
	Currency Currency       `json:"currency"`
	Metadata AmountMetadata `json:"metadata,omitempty"`
}

type Coin struct {
	CoinIdentifier CoinIdentifier `json:"coin_identifier"`
	Amount         Amount         `json:"amount"`
}

type CoinChange struct {
	CoinIdentifier CoinIdentifier `json:"coin_identifier"`
	CoinAction     CoinAction     `json:"coin_action"`
}

type Operation struct {
	OperationIdentifier OperationIdentifier   `json:"operation_identifier"`
	RelatedOperations   []OperationIdentifier `json:"related_operations,omitempty"`
	Type                string                `json:"type"`
	Status              string                `json:"status,omitempty"`
	Account             AccountIdentifier     `json:"account,omitempty"`
	Amount              Amount                `json:"amount,omitempty"`
	CoinChange          CoinChange            `json:"coin_change,omitempty"`
	Metadata            OperationMetadata     `json:"metadata,omitempty"`
}

type RelatedTransaction struct {
	NetworkIdentifier     NetworkIdentifier     `json:"network_identifier,omitempty"`
	TransactionIdentifier TransactionIdentifier `json:"transaction_identifier"`
	Direction             Direction             `json:"direction"`
}

type Transaction struct {
	TransactionIdentifier TransactionIdentifier `json:"transaction_identifier"`
	Operations            []Operation           `json:"operations"`
	RelatedTransactions   []RelatedTransaction  `json:"related_transactions,omitempty"`
	Metadata              TransactionMetadata   `json:"metadata,omitempty"`
}

type Block struct {
	BlockIdentifier       BlockIdentifier `json:"block_identifier"`
	ParentBlockIdentifier BlockIdentifier `json:"parent_block_identifier"`
	Timestamp             int64           `json:"timestamp"`
	Transactions          []Transaction   `json:"transactions"`
	Metadata              BlockMetadata   `json:"metadata,omitempty"`
}

type BlockEvent struct {
	Sequence        int64           `json:"sequence"`
	BlockIdentifier BlockIdentifier `json:"block_identifier"`
	Type            BlockEventType  `json:"type"`
}

type BlockTransactions struct {
	BlockIdentifier BlockIdentifier `json:"block_identifier"`
	Transaction     Transaction     `json:"transaction"`
}
