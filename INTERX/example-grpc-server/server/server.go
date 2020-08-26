package server

import (
	"context"
	"sync"

	pbExample "github.com/KiraCore/sekai/INTERX/proto-gen"
)

// Backend implements the protobuf interface
type Backend struct {
	mu    *sync.RWMutex
}

// New initializes a new Backend struct.
func New() *Backend {
	return &Backend{
		mu: &sync.RWMutex{},
	}
}

// Handle Ping.
func (b *Backend) Ping(ctx context.Context, _ *pbExample.PingRequest) (*pbExample.PingResponse, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	response := &pbExample.PingResponse{
		Message: "success",
	}

	return response, nil
}
