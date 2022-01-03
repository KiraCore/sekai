package indexerapi

import "github.com/KiraCore/sekai/INTERX/types/rosetta"

type EventsBlocksRequest struct {
	NetworkIdentifier rosetta.NetworkIdentifier `json:"network_identifier"`
	Offset            int64                     `json:"offset,omitempty"`
	Limit             int64                     `json:"limit,omitempty"`
}

type EventsBlocksResponse struct {
	MaxSequence int64                `json:"max_sequence"`
	Events      []rosetta.BlockEvent `json:"events"`
}
