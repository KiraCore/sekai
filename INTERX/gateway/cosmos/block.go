package cosmos

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterCosmosBlockRoutes registers query routers.
func RegisterCosmosBlockRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryBlocks, QueryBlocksRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryBlockByHeight, QueryBlockByHeightRequest(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryBlocks, "This is an API to query blocks by parameters.", true)
	common.AddRPCMethod("GET", config.QueryBlockByHeight, "This is an API to query block by height", true)
}

func queryBlocksHandle(rpcAddr string, r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	_ = r.ParseForm()

	minHeight := r.FormValue("minHeight")
	maxHeight := r.FormValue("maxHeight")

	var events = make([]string, 0, 2)

	if minHeight != "" {
		events = append(events, fmt.Sprintf("minHeight=%s", minHeight))
	}

	if maxHeight != "" {
		events = append(events, fmt.Sprintf("maxHeight=%s", maxHeight))
	}

	// search blocks

	return common.MakeGetRequest(rpcAddr, "/blockchain", strings.Join(events, "&"))
}

// QueryBlocksRequest is a function to query Blocks.
func QueryBlocksRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-blocks] Entering Blocks query")

		if !common.RPCMethods["GET"][config.QueryBlocks].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryBlocks].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-blocks] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryBlocksHandle(rpcAddr, r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryBlocks].CachingEnabled)
	}
}

func queryBlockByHeightHandle(rpcAddr string, height string) (interface{}, interface{}, int) {
	return common.MakeGetRequest(rpcAddr, "/block", fmt.Sprintf("height=%s", height))
}

// QueryBlockByHeightRequest is a function to query Blocks.
func QueryBlockByHeightRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		height := queries["height"]
		request := types.InterxRequest{
			Method:   r.Method,
			Endpoint: config.QueryBalances,
			Params:   []byte(height),
		}
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-blocks-by-height] Entering Block query by height: %s", height)

		if !common.RPCMethods["GET"][config.QueryBlocks].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryBlocks].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-blocks-by-height] Returning from the cache: %s", height)
					return
				}
			}

			response.Response, response.Error, statusCode = queryBlockByHeightHandle(rpcAddr, height)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryBlocks].CachingEnabled)
	}
}
