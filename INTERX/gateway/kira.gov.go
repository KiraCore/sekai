package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	common "github.com/KiraCore/sekai/INTERX/common"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

const (
	queryDataReferenceKeys = "/api/kira/gov/data_keys"
	queryDataReference     = "/api/kira/gov/data"
)

// RegisterKiraGovRoutes registers kira gov query routers.
func RegisterKiraGovRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(queryDataReferenceKeys, QueryDataReferenceKeysRequest(gwCosmosmux, rpcAddr)).Methods(GET)
	r.HandleFunc("/api/kira/gov/data/{key}", QueryDataReferenceRequest(gwCosmosmux, rpcAddr)).Methods(GET)

	AddRPCMethod(GET, queryDataReferenceKeys, "This is an API to query all data reference keys.", true)
	AddRPCMethod(GET, queryDataReference, "This is an API to query data reference by key.", true)
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

	return ServeGRPC(r, gwCosmosmux)
}

// QueryDataReferenceKeysRequest is a function to query data reference keys.
func QueryDataReferenceKeysRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := GetInterxRequest(r)
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !rpcMethods[GET][queryDataReferenceKeys].Enabled {
			response.Response, response.Error, statusCode = ServeError(0, "", "", http.StatusForbidden)
		} else {
			if rpcMethods[GET][queryDataReferenceKeys].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = queryDataReferenceKeysHandle(r, gwCosmosmux)
		}

		WrapResponse(w, request, *response, statusCode, rpcMethods[GET][queryTotalSupply].CachingEnabled)
	}
}

func queryDataReferenceHandle(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	success, failure, status := ServeGRPC(r, gwCosmosmux)

	if success != nil {
		type DataReferenceTempResponse struct {
			Data common.DataReferenceEntry `json:"data"`
		}
		result := DataReferenceTempResponse{}

		byteData, _ := json.Marshal(success)
		json.Unmarshal(byteData, &result)

		success = result.Data

		common.UpdateKey(mux.Vars(r)["key"], result.Data)
	}

	return success, failure, status
}

// QueryDataReferenceRequest is a function to query data reference by key.
func QueryDataReferenceRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		key := queries["key"]
		request := InterxRequest{
			Method:   r.Method,
			Endpoint: queryDataReference,
			Params:   []byte(key),
		}
		response := GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !rpcMethods[GET][queryDataReference].Enabled {
			response.Response, response.Error, statusCode = ServeError(0, "", "", http.StatusForbidden)
		} else {
			if rpcMethods[GET][queryDataReference].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					WrapResponse(w, request, *response, statusCode, false)
					return
				}
			}

			response.Response, response.Error, statusCode = queryDataReferenceHandle(r, gwCosmosmux)
		}

		WrapResponse(w, request, *response, statusCode, rpcMethods[GET][queryDataReference].CachingEnabled)
	}
}
