package kira

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func RegisterIdentityRegistrarRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryIdentityRecord, QueryIdentityRecordRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryIdentityRecordsByAddress, QueryIdentityRecordsByAddressRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryAllIdentityRecords, QueryAllIdentityRecordsRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryIdentityRecordVerifyRequest, QueryIdentityRecordVerifyRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryIdentityRecordVerifyRequestsByRequester, QueryIdentityRecordVerifyRequestsByRequester(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryIdentityRecordVerifyRequestsByApprover, QueryIdentityRecordVerifyRequestsByApprover(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryAllIdentityRecordVerifyRequests, QueryAllIdentityRecordVerifyRequests(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryIdentityRecord, "This is an API to query identity record by id.", true)
	common.AddRPCMethod("GET", config.QueryIdentityRecordsByAddress, "This is an API to query identity records by address.", true)
	common.AddRPCMethod("GET", config.QueryAllIdentityRecords, "This is an API to query all identity records.", true)
	common.AddRPCMethod("GET", config.QueryIdentityRecordVerifyRequest, "This is an API to query identity record verify request.", true)
	common.AddRPCMethod("GET", config.QueryIdentityRecordVerifyRequestsByRequester, "This is an API to query identity record verify request by requester.", true)
	common.AddRPCMethod("GET", config.QueryIdentityRecordVerifyRequestsByApprover, "This is an API to query identity record verify request by approver.", true)
	common.AddRPCMethod("GET", config.QueryAllIdentityRecordVerifyRequests, "This is an API to query all identity record verify requests.", true)
}

func QueryIdentityRecordHandler(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	return common.ServeGRPC(r, gwCosmosmux)
}

func QueryIdentityRecordRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-identity-record] entering query")

		if !common.RPCMethods["GET"][config.QueryIdentityRecord].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryIdentityRecord].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-identity-record] returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = QueryIdentityRecordHandler(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryIdentityRecord].CachingEnabled)
	}
}

func QueryIdentityRecordsByAddressHandler(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	queries := mux.Vars(r)
	creator := queries["creator"]

	accAddr, _ := sdk.AccAddressFromBech32(creator)
	r.URL.Path = fmt.Sprintf("/api/kira/gov/identity_records/%s", base64.URLEncoding.EncodeToString(accAddr.Bytes()))
	return common.ServeGRPC(r, gwCosmosmux)
}

func QueryIdentityRecordsByAddressRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-identity-records-by-address] entering query")

		if !common.RPCMethods["GET"][config.QueryIdentityRecordsByAddress].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryIdentityRecordsByAddress].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-identity-records-by-address] returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = QueryIdentityRecordsByAddressHandler(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryIdentityRecordsByAddress].CachingEnabled)
	}
}

func QueryAllIdentityRecordsHandler(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
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

func QueryAllIdentityRecordsRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-all-identity-records] entering query")

		if !common.RPCMethods["GET"][config.QueryAllIdentityRecords].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryAllIdentityRecords].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-all-identity-records] returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = QueryAllIdentityRecordsHandler(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryAllIdentityRecords].CachingEnabled)
	}
}

func QueryIdentityRecordVerifyRequestHandler(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	return common.ServeGRPC(r, gwCosmosmux)
}

func QueryIdentityRecordVerifyRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-identity-record-verify-request] entering query")

		if !common.RPCMethods["GET"][config.QueryIdentityRecordVerifyRequest].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryIdentityRecordVerifyRequest].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-identity-record-verify-request] returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = QueryIdentityRecordVerifyRequestHandler(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryIdentityRecordVerifyRequest].CachingEnabled)
	}
}

func QueryIdentityRecordVerifyRequestsByRequesterHandler(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	params := mux.Vars(r)
	requester := params["requester"]
	accAddr, _ := sdk.AccAddressFromBech32(requester)
	r.URL.Path = fmt.Sprintf("/api/kira/gov/identity_verify_requests_by_requester/%s", base64.URLEncoding.EncodeToString(accAddr.Bytes()))

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

func QueryIdentityRecordVerifyRequestsByRequester(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-identity-record-verify-request-by-requester] entering query")

		if !common.RPCMethods["GET"][config.QueryIdentityRecordVerifyRequestsByRequester].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryIdentityRecordVerifyRequestsByRequester].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-identity-record-verify-request-by-requester] returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = QueryIdentityRecordVerifyRequestsByRequesterHandler(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryIdentityRecordVerifyRequestsByRequester].CachingEnabled)
	}
}

func QueryIdentityRecordVerifyRequestsByApproverHandler(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	params := mux.Vars(r)
	approver := params["approver"]
	accAddr, _ := sdk.AccAddressFromBech32(approver)
	r.URL.Path = fmt.Sprintf("/api/kira/gov/identity_verify_requests_by_approver/%s", base64.URLEncoding.EncodeToString(accAddr.Bytes()))

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

func QueryIdentityRecordVerifyRequestsByApprover(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-identity-record-verify-request-by-approver] entering query")

		if !common.RPCMethods["GET"][config.QueryIdentityRecordVerifyRequestsByApprover].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryIdentityRecordVerifyRequestsByApprover].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-identity-record-verify-request-by-approver] returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = QueryIdentityRecordVerifyRequestsByApproverHandler(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryIdentityRecordVerifyRequestsByApprover].CachingEnabled)
	}
}

func QueryAllIdentityRecordVerifyRequestsHandler(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
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

func QueryAllIdentityRecordVerifyRequests(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-all-identity-record-verify-requests] entering query")

		if !common.RPCMethods["GET"][config.QueryAllIdentityRecords].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryAllIdentityRecords].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-all-identity-record-verify-requests] returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = QueryAllIdentityRecordVerifyRequestsHandler(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryAllIdentityRecords].CachingEnabled)
	}
}
