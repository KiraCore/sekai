package dataapi

import (
	"github.com/KiraCore/sekai/INTERX/types/rosetta"
)

type MetadataRequest struct {
	Metadata interface{} `json:"metadata,omitempty"`
}

type NetworkListRequest MetadataRequest

type NetworkListResponse struct {
	NetworkIdentifiers []rosetta.NetworkIdentifier `json:"network_identifiers"`
}

type NetworkRequest struct {
	NetworkIdentifier rosetta.NetworkIdentifier `json:"network_identifier,omitempty"` // make it omitable
	Metadata          interface{}               `json:"metadata,omitempty"`
}

type NetworkOptionsRequest NetworkRequest

type NetworkOptionsResponse struct {
	Version rosetta.Version `json:"version"`
	Allow   rosetta.Allow   `json:"allow"`
}

type NetworkStatusRequest NetworkRequest

type NetworkStatusResponse struct {
	CurrentBlockIdentifier rosetta.BlockIdentifier `json:"current_block_identifier"`
	CurrentBlockTimestamp  int64                   `json:"current_block_timestamp"`
	GenesisBlockIdentifier rosetta.BlockIdentifier `json:"genesis_block_identifier"`
	OldestBlockIdentifier  rosetta.BlockIdentifier `json:"oldest_block_identifier,omitempty"`
	SyncStatus             rosetta.SyncStatus      `json:"sync_status,omitempty"`
	Peers                  []rosetta.Peer          `json:"peers"`
}
