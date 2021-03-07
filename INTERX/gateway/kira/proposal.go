package kira

import (
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterKiraGovProposalRoutes registers kira gov proposal query routers.
func RegisterKiraGovProposalRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryProposals, QueryProposalsRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryProposal, QueryProposalRequest(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryProposals, "This is an API to query all proposals.", true)
	common.AddRPCMethod("GET", config.QueryProposal, "This is an API to query a proposal by a given id.", true)
}

func queryProposalsHandler(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	return common.ServeGRPC(r, gwCosmosmux)
}

// QueryProposalsRequest is a function to query all proposals.
func QueryProposalsRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-proposals] Entering proposals query")

		if !common.RPCMethods["GET"][config.QueryProposals].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryProposals].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-proposals] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryProposalsHandler(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryProposals].CachingEnabled)
	}
}

func queryProposalHandler(r *http.Request, gwCosmosmux *runtime.ServeMux, proposal_id string) (interface{}, interface{}, int) {
	return common.ServeGRPC(r, gwCosmosmux)
}

// QueryProposalRequest is a function to query a proposal by a given proposal_id.
func QueryProposalRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		proposalID := queries["proposal_id"]
		request := types.InterxRequest{
			Method:   r.Method,
			Endpoint: config.QueryProposal,
			Params:   []byte(proposalID),
		}
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-proposal] Entering proposal query by proposal_id: ", proposalID)

		if !common.RPCMethods["GET"][config.QueryProposal].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryProposal].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-reference] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryProposalHandler(r, gwCosmosmux, proposalID)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryProposal].CachingEnabled)
	}
}
