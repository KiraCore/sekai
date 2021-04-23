package interx

import (
	"encoding/json"
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	functions "github.com/KiraCore/sekai/INTERX/functions"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterInterxQueryRoutes registers query routers.
func RegisterInterxQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryRPCMethods, QueryRPCMethods(rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryInterxFunctions, QueryInterxFunctions(rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryStatus, QueryStatusRequest(rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryInterxFunctions, "This is an API to query interx functions.", true)
	common.AddRPCMethod("GET", config.QueryStatus, "This is an API to query interx status.", true)
}

// QueryRPCMethods is a function to query RPC methods.
func QueryRPCMethods(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response = common.RPCMethods

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}

func queryInterxFunctionsHandle(rpcAddr string) (interface{}, interface{}, int) {
	functions := functions.GetInterxFunctions()

	return functions, nil, http.StatusOK
}

// QueryInterxFunctions is a function to list functions and metadata.
func QueryInterxFunctions(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response, response.Error, statusCode = queryInterxFunctionsHandle(rpcAddr)

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}

func queryStatusHandle(rpcAddr string) (interface{}, interface{}, int) {
	type StatusTempResponse struct {
		InterxInfo struct {
			PubKey interface{} `json:"pub_key,omitempty"`
		} `json:"interx_info,omitempty"`
		Moniker           string `json:"moniker"`
		KiraAddr          string `json:"kira_addr"`
		GenesisChecksum   string `json:"genesis_checksum"`
		ChainID           string `json:"chain_id"`
		Version           string `json:"version"`
		LatestBlockHeight string `json:"latest_block_height"`
		CatchingUp        bool   `json:"catching_up"`
		SentryNodeID      string `json:"sentry_node_id,omitemtpy"`
		PrivSentryNodeID  string `json:"priv_sentry_node_id,omitemtpy"`
		ValidatorNodeID   string `json:"validator_node_id,omitemtpy"`
		SeedNodeID        string `json:"seed_node_id,omitemtpy"`
		InterxVersion     string `json:"interx_version"`
	}

	result := StatusTempResponse{}

	// Handle Interx Pubkey
	pubkeyBytes, err := config.EncodingCg.Amino.MarshalJSON(config.Config.PubKey)
	if err != nil {
		common.GetLogger().Error("[query-status] Failed to marshal interx pubkey", err)
		return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
	}

	err = json.Unmarshal(pubkeyBytes, &result.InterxInfo.PubKey)
	if err != nil {
		common.GetLogger().Error("[query-status] Failed to add interx pubkey to status response", err)
		return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
	}

	// Handle Genesis
	genesis, checksum, err := GetGenesisResults(rpcAddr)
	if err != nil {
		common.GetLogger().Error("[query-status] Failed to query genesis ", err)
		return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
	}

	// Get Kira Status
	sentryStatus := common.GetKiraStatus((rpcAddr))

	if sentryStatus != nil {
		result.Moniker = sentryStatus.NodeInfo.Moniker
		result.Version = sentryStatus.NodeInfo.Version
	}

	result.SentryNodeID = config.Config.SentryNodeID
	result.PrivSentryNodeID = config.Config.PrivSentryNodeID
	result.ValidatorNodeID = config.Config.ValidatorNodeID
	result.SeedNodeID = config.Config.SeedNodeID

	result.KiraAddr = config.Config.Address
	result.GenesisChecksum = checksum
	result.ChainID = genesis.ChainID

	result.LatestBlockHeight = sentryStatus.SyncInfo.LatestBlockHeight
	result.CatchingUp = sentryStatus.SyncInfo.CatchingUp

	result.InterxVersion = config.InterxVersion

	return result, nil, http.StatusOK
}

// QueryStatusRequest is a function to query status.
func QueryStatusRequest(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-status] Entering status query")

		if !common.RPCMethods["GET"][config.QueryStatus].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryStatus].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-status] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryStatusHandle(rpcAddr)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryStatus].CachingEnabled)
	}
}
