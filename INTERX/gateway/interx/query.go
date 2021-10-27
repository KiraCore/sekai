package interx

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	functions "github.com/KiraCore/sekai/INTERX/functions"
	"github.com/KiraCore/sekai/INTERX/tasks"
	"github.com/KiraCore/sekai/INTERX/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterInterxQueryRoutes registers query routers.
func RegisterInterxQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryRPCMethods, QueryRPCMethods(rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryInterxFunctions, QueryInterxFunctions(rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryStatus, QueryStatusRequest(rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryAddrBook, QueryAddrBook(rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryNetInfo, QueryNetInfo(rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryInterxFunctions, "This is an API to query interx functions.", true)
	common.AddRPCMethod("GET", config.QueryStatus, "This is an API to query interx status.", true)
	common.AddRPCMethod("GET", config.QueryAddrBook, "This is an API to query address book.", true)
	common.AddRPCMethod("GET", config.QueryNetInfo, "This is an API to query net info.", true)
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
	metadata := functions.GetInterxMetadata()

	return metadata, nil, http.StatusOK
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
	result := types.InterxStatus{}

	// Handle Interx Pubkey
	pubkeyBytes, err := config.EncodingCg.Amino.MarshalJSON(&config.Config.PubKey)

	if err != nil {
		common.GetLogger().Error("[query-status] Failed to marshal interx pubkey", err)
		return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
	}

	err = json.Unmarshal(pubkeyBytes, &result.InterxInfo.PubKey)
	if err != nil {
		common.GetLogger().Error("[query-status] Failed to add interx pubkey to status response", err)
		return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
	}

	result.ID = hex.EncodeToString(config.Config.PubKey.Address())

	// Handle Genesis
	genesis, checksum, err := GetGenesisResults(rpcAddr)
	if err != nil {
		common.GetLogger().Error("[query-status] Failed to query genesis ", err)
		return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
	}

	// Get Kira Status
	sentryStatus := common.GetKiraStatus((rpcAddr))

	if sentryStatus != nil {
		result.NodeInfo = sentryStatus.NodeInfo
		result.SyncInfo = sentryStatus.SyncInfo
		result.ValidatorInfo = sentryStatus.ValidatorInfo

		result.InterxInfo.Moniker = sentryStatus.NodeInfo.Moniker

		result.InterxInfo.LatestBlockHeight = sentryStatus.SyncInfo.LatestBlockHeight
		result.InterxInfo.CatchingUp = sentryStatus.SyncInfo.CatchingUp
	}

	result.InterxInfo.Node = config.Config.Node

	result.InterxInfo.KiraAddr = config.Config.Address
	result.InterxInfo.KiraPubKey = config.Config.PubKey.String()
	result.InterxInfo.FaucetAddr = config.Config.Faucet.Address
	result.InterxInfo.GenesisChecksum = checksum
	result.InterxInfo.ChainID = genesis.ChainID

	result.InterxInfo.Version = config.Config.Version

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

func queryAddrBookHandler(rpcAddr string) (interface{}, interface{}, int) {
	return config.LoadAddressBooks(), nil, http.StatusOK
}

// QueryAddrBook is a function to query address book.
func QueryAddrBook(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response, response.Error, statusCode = queryAddrBookHandler(rpcAddr)

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}

func queryNetInfoHandler(rpcAddr string) (interface{}, interface{}, int) {
	netInfo, err := tasks.QueryNetInfo(rpcAddr)
	if err != nil {
		return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
	}

	return netInfo, nil, http.StatusOK
}

// QueryNetInfo is a function to query net info.
func QueryNetInfo(rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response, response.Error, statusCode = queryNetInfoHandler(rpcAddr)

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}
