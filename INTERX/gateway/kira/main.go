package kira

import (
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterRequest is a function to register requests.
func RegisterRequest(router *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	RegisterKiraGovRoutes(router, gwCosmosmux, rpcAddr)
	RegisterKiraGovProposalRoutes(router, gwCosmosmux, rpcAddr)
	RegisterKiraQueryRoutes(router, gwCosmosmux, rpcAddr)
	RegisterKiraTokensRoutes(router, gwCosmosmux, rpcAddr)
	RegisterKiraUpgradeRoutes(router, gwCosmosmux, rpcAddr)
	RegisterIdentityRegistrarRoutes(router, gwCosmosmux, rpcAddr)
	RegisterKiraGovRoleRoutes(router, gwCosmosmux, rpcAddr)
	RegisterKiraGovPermissionRoutes(router, gwCosmosmux, rpcAddr)
}
