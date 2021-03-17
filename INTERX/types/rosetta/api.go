package rosetta

type Version struct {
	RosettaVersion    string      `json:"rosetta_version"`
	NodeVersion       string      `json:"node_version"`
	MiddlewareVersion string      `json:"middleware_version,omitempty"`
	Metadata          interface{} `json:"metadata,omitempty"`
}

type SyncStatus struct {
	CurrentIndex int64  `json:"current_index,omitempty"`
	TargetIndex  int64  `json:"target_index,omitempty"`
	Stage        string `json:"stage,omitempty"`
	Synced       bool   `json:"synced,omitempty"`
}

type Peer struct {
	PeerID   string      `json:"peer_id"`
	Metadata interface{} `json:"metadata,omitempty"`
}
