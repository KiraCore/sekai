package types

type ProtocolVersion struct {
	P2P   string `json:"p2p,omitempty"`
	Block string `json:"block,omitempty"`
	App   string `json:"app,omitempty"`
}

type NodeOtherInfo struct {
	TxIndex    string `json:"tx_index,omitempty"`
	RpcAddress string `json:"rpc_address,omitempty"`
}

type SyncInfo struct {
	LatestBlockHash     string `json:"latest_block_hash,omitempty"`
	LatestAppHash       string `json:"latest_app_hash,omitempty"`
	LatestBlockHeight   string `json:"latest_block_height,omitempty"`
	LatestBlockTime     string `json:"latest_block_time,omitempty"`
	EarliestBlockHash   string `json:"earliest_block_hash,omitempty"`
	EarliestAppHash     string `json:"earliest_app_hash,omitempty"`
	EarliestBlockHeight string `json:"earliest_block_height,omitempty"`
	EarliestBlockTime   string `json:"earliest_block_time,omitempty"`
	CatchingUp          bool   `json:"catching_up,omitempty"`
}

type NodeInfo struct {
	ProtocolVersion ProtocolVersion `json:"protocol_version,omitempty"`
	Id              string          `json:"id,omitempty"`
	ListenAddr      string          `json:"listen_addr,omitempty"`
	Network         string          `json:"network,omitempty"`
	Version         string          `json:"version,omitempty"`
	Channels        string          `json:"channels,omitempty"`
	Moniker         string          `json:"moniker,omitempty"`
	Other           NodeOtherInfo   `json:"other,omitempty"`
}

type PubKey struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type ValidatorInfo struct {
	Address     string  `json:"address,omitempty"`
	PubKey      *PubKey `json:"pub_key,omitempty"`
	VotingPower string  `json:"voting_power,omitempty"`
}

type KiraStatus struct {
	NodeInfo      NodeInfo      `json:"node_info,omitempty"`
	SyncInfo      SyncInfo      `json:"sync_info,omitempty"`
	ValidatorInfo ValidatorInfo `json:"validator_info,omitempty"`
}
