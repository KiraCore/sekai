package types

type InterxStatus struct {
	ID         string `json:"id"`
	InterxInfo struct {
		PubKey interface{} `json:"pub_key,omitempty"`
	} `json:"interx_info,omitempty"`
	Moniker           string `json:"moniker"`
	KiraAddr          string `json:"kira_addr"`
	GenesisChecksum   string `json:"genesis_checksum"`
	ChainID           string `json:"chain_id"`
	Version           string `json:"version"`
	LatestBlockHeight string `json:"latest_block_height"`
	CatchingUp        bool   `json:"catching_up"`
	SentryNodeID      string `json:"sentry_node_id,omitemtpy"`
	PrivSentryNodeID  string `json:"priv_sentry_node_id,omitemtpy"`
	ValidatorNodeID   string `json:"validator_node_id,omitemtpy"`
	SeedNodeID        string `json:"seed_node_id,omitemtpy"`
	InterxVersion     string `json:"interx_version"`
}
