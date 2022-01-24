package types

type NodeConfig struct {
	NodeType        string `json:"node_type"`
	SentryNodeID    string `json:"sentry_node_id"`
	SnapshotNodeID  string `json:"snapshot_node_id"`
	ValidatorNodeID string `json:"validator_node_id"`
	SeedNodeID      string `json:"seed_node_id"`
}

type InterxStatus struct {
	ID         string `json:"id"`
	InterxInfo struct {
		PubKey struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"pub_key,omitempty"`
		Moniker           string     `json:"moniker"`
		KiraAddr          string     `json:"kira_addr"`
		KiraPubKey        string     `json:"kira_pub_key"`
		FaucetAddr        string     `jsong:"faucet_addr"`
		GenesisChecksum   string     `json:"genesis_checksum"`
		ChainID           string     `json:"chain_id"`
		Version           string     `json:"version"`
		LatestBlockHeight string     `json:"latest_block_height"`
		CatchingUp        bool       `json:"catching_up"`
		Node              NodeConfig `json:"node"`
	} `json:"interx_info,omitempty"`
	NodeInfo      NodeInfo      `json:"node_info,omitempty"`
	SyncInfo      SyncInfo      `json:"sync_info,omitempty"`
	ValidatorInfo ValidatorInfo `json:"validator_info,omitempty"`
}

type SnapShotChecksumResponse struct {
	Size     int64  `json:"size,omitempty"`
	Checksum string `json:"checksum,omitempty"`
}
