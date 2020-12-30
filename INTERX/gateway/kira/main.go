package kira

import (
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterRequest is a function to register requests.
func RegisterRequest(router *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	RegisterKiraGovRoutes(router, gwCosmosmux, rpcAddr)
	RegisterKiraQueryRoutes(router, gwCosmosmux, rpcAddr)
}
