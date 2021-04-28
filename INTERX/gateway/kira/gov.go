package kira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	database "github.com/KiraCore/sekai/INTERX/database"
	"github.com/KiraCore/sekai/INTERX/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterKiraGovRoutes registers kira gov query routers.
func RegisterKiraGovRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryDataReferenceKeys, QueryDataReferenceKeysRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryDataReference, QueryDataReferenceRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryNetworkProperties, QueryNetworkPropertiesRequest(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryDataReferenceKeys, "This is an API to query all data reference keys.", true)
	common.AddRPCMethod("GET", config.QueryDataReference, "This is an API to query data reference by key.", true)
	common.AddRPCMethod("GET", config.QueryNetworkProperties, "This is an API to query network properties.", true)
}

func queryDataReferenceKeysHandle(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	queries := r.URL.Query()
	key := queries["key"]
	offset := queries["offset"]
	limit := queries["limit"]
	countTotal := queries["count_total"]

	var events = make([]string, 0, 4)
	if len(key) == 1 {
		events = append(events, fmt.Sprintf("pagination.key=%s", key[0]))
	}
	if len(offset) == 1 {
		events = append(events, fmt.Sprintf("pagination.offset=%s", offset[0]))
	}
	if len(limit) == 1 {
		events = append(events, fmt.Sprintf("pagination.limit=%s", limit[0]))
	}
	if len(countTotal) == 1 {
		events = append(events, fmt.Sprintf("pagination.count_total=%s", countTotal[0]))
	}

	r.URL.RawQuery = strings.Join(events, "&")

	return common.ServeGRPC(r, gwCosmosmux)
}

// QueryDataReferenceKeysRequest is a function to query data reference keys.
func QueryDataReferenceKeysRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-reference-keys] Entering data reference keys query")

		if !common.RPCMethods["GET"][config.QueryDataReferenceKeys].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryDataReferenceKeys].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-reference-keys] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryDataReferenceKeysHandle(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryDataReferenceKeys].CachingEnabled)
	}
}

func queryDataReferenceHandle(r *http.Request, gwCosmosmux *runtime.ServeMux, key string) (interface{}, interface{}, int) {
	success, failure, status := common.ServeGRPC(r, gwCosmosmux)

	if success != nil {
		type DataReferenceTempResponse struct {
			Data types.DataReferenceEntry `json:"data"`
		}
		result := DataReferenceTempResponse{}

		byteData, err := json.Marshal(success)
		if err != nil {
			common.GetLogger().Error("[query-reference] Invalid response format", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		err = json.Unmarshal(byteData, &result)
		if err != nil {
			common.GetLogger().Error("[query-reference] Invalid response format", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		success = result.Data

		filePath := key + filepath.Ext(result.Data.Reference)

		database.AddReference(key, result.Data.Reference, 0, time.Now(), config.DataReferenceRegistry+"/"+common.GetMD5Hash(filePath))
	}

	return success, failure, status
}

// QueryDataReferenceRequest is a function to query data reference by key.
func QueryDataReferenceRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		key := queries["key"]
		request := types.InterxRequest{
			Method:   r.Method,
			Endpoint: config.QueryDataReference,
			Params:   []byte(key),
		}
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-reference] Entering data reference query by key: ", key)

		if !common.RPCMethods["GET"][config.QueryDataReference].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryDataReference].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-reference] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryDataReferenceHandle(r, gwCosmosmux, key)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryDataReference].CachingEnabled)
	}
}

func QueryNetworkPropertiesHandle(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	return common.ServeGRPC(r, gwCosmosmux)
}

// QueryDataReferenceKeysRequest is a function to query data reference keys.
func QueryNetworkPropertiesRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-network-properties] Entering network properties query")

		if !common.RPCMethods["GET"][config.QueryNetworkProperties].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryNetworkProperties].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-network-properties] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = QueryNetworkPropertiesHandle(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryNetworkProperties].CachingEnabled)
	}
}
