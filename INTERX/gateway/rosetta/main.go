package rosetta

import (
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/KiraCore/sekai/INTERX/gateway/rosetta/data"
)

// RegisterRequest is a function to register requests.
func RegisterRequest(router *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	data.RegisterRequest(router, gwCosmosmux, rpcAddr)
}
