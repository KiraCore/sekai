package kira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/KiraCore/sekai/INTERX/common"
	database "github.com/KiraCore/sekai/INTERX/database"
	"github.com/KiraCore/sekai/INTERX/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterKiraGovRoutes registers kira gov query routers.
func RegisterKiraGovRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(common.QueryDataReferenceKeys, QueryDataReferenceKeysRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc("/api/kira/gov/data/{key}", QueryDataReferenceRequest(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", common.QueryDataReferenceKeys, "This is an API to query all data reference keys.", true)
	common.AddRPCMethod("GET", common.QueryDataReference, "This is an API to query data reference by key.", true)
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

		if !common.RPCMethods["GET"][common.QueryDataReferenceKeys].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][common.QueryDataReferenceKeys].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = queryDataReferenceKeysHandle(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][common.QueryDataReferenceKeys].CachingEnabled)
	}
}

func queryDataReferenceHandle(r *http.Request, gwCosmosmux *runtime.ServeMux, key string) (interface{}, interface{}, int) {
	success, failure, status := common.ServeGRPC(r, gwCosmosmux)

	if success != nil {
		type DataReferenceTempResponse struct {
			Data types.DataReferenceEntry `json:"data"`
		}
		result := DataReferenceTempResponse{}

		byteData, _ := json.Marshal(success)
		json.Unmarshal(byteData, &result)

		success = result.Data

		filePath := key + filepath.Ext(result.Data.Reference)

		database.AddReference(key, result.Data.Reference, 0, time.Now(), common.DataReferenceRegistry+"/"+common.GetMD5Hash(filePath))
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
			Endpoint: common.QueryDataReference,
			Params:   []byte(key),
		}
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !common.RPCMethods["GET"][common.QueryDataReference].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][common.QueryDataReference].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = queryDataReferenceHandle(r, gwCosmosmux, key)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][common.QueryDataReference].CachingEnabled)
	}
}
