package gateway

import (
	"github.com/KiraCore/sekai/INTERX/gateway/cosmos"
	"github.com/KiraCore/sekai/INTERX/gateway/interx"
	"github.com/KiraCore/sekai/INTERX/gateway/kira"
	"github.com/KiraCore/sekai/INTERX/gateway/rosetta"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterRequest is a function to register requests.
func RegisterRequest(router *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	cosmos.RegisterRequest(router, gwCosmosmux, rpcAddr)
	kira.RegisterRequest(router, gwCosmosmux, rpcAddr)
	interx.RegisterRequest(router, gwCosmosmux, rpcAddr)
	rosetta.RegisterRequest(router, gwCosmosmux, rpcAddr)
}
