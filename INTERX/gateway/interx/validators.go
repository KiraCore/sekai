package interx

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterValidatorsQueryRoutes registers validators query routers.
func RegisterValidatorsQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryValidators, QueryValidators(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryValidatorInfos, QueryValidatorInfos(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryValidators, "This is an API to query validators.", true)
	common.AddRPCMethod("GET", config.QueryValidatorInfos, "This is an API to query validator infos.", true)
}

const (
	// Undefined status
	Undefined string = "UNDEFINED"
	// Active status
	Active string = "ACTIVE"
	// Inactive status
	Inactive string = "INACTIVE"
	// Paused status
	Paused string = "PAUSED"
	// Jailed status
	Jailed string = "JAILED"
)

func queryValidatorsHandle(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	queries := r.URL.Query()
	address := queries["address"]
	valkey := queries["valkey"]
	pubkey := queries["pubkey"]
	moniker := queries["moniker"]
	status := queries["status"]
	key := queries["key"]
	offset := queries["offset"]
	limit := queries["limit"]
	proposer := queries["proposer"]
	countTotal := queries["count_total"]
	all := queries["all"]

	isQueryAll := false

	var events = make([]string, 0, 9)
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
	if len(address) == 1 {
		events = append(events, fmt.Sprintf("address=%s", address[0]))
	}
	if len(valkey) == 1 {
		events = append(events, fmt.Sprintf("valkey=%s", valkey[0]))
	}
	if len(pubkey) == 1 {
		events = append(events, fmt.Sprintf("pubkey=%s", pubkey[0]))
	}
	if len(proposer) == 1 {
		events = append(events, fmt.Sprintf("proposer=%s", proposer[0]))
	}
	if len(moniker) == 1 {
		events = append(events, fmt.Sprintf("moniker=%s", moniker[0]))
	}
	if len(status) == 1 {
		events = append(events, fmt.Sprintf("status=%s", status[0]))
	}
	if len(all) == 1 {
		events = append(events, fmt.Sprintf("all=%s", all[0]))
		isQueryAll = all[0] == "true"
	}

	r.URL.RawQuery = strings.Join(events, "&")

	success, failure, statusCode := common.ServeGRPC(r, gwCosmosmux)
	if success != nil {
		result := struct {
			Validators []types.QueryValidator `json:"validators,omitempty"`
			Actors     []string               `json:"actors,omitempty"`
			Pagination interface{}            `json:"pagination,omitempty"`
		}{}

		byteData, err := json.Marshal(success)
		if err != nil {
			common.GetLogger().Error("[query-reference] Invalid response format: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		err = json.Unmarshal(byteData, &result)
		if err != nil {
			common.GetLogger().Error("[query-reference] Invalid response format: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		if isQueryAll {

			allValidators := types.AllValidators{}

			allValidators.Validators = result.Validators
			allValidators.Waiting = make([]string, 0)
			for _, actor := range result.Actors {
				isWaiting := true
				for _, validator := range result.Validators {
					if validator.Address == actor {
						isWaiting = false
						break
					}
				}

				if isWaiting {
					allValidators.Waiting = append(allValidators.Waiting, actor)
				}
			}

			allValidators.Status.TotalValidators = len(result.Validators)
			allValidators.Status.WaitingValidators = len(allValidators.Waiting)

			allValidators.Status.ActiveValidators = 0
			allValidators.Status.PausedValidators = 0
			allValidators.Status.InactiveValidators = 0
			allValidators.Status.JailedValidators = 0
			for _, validator := range result.Validators {
				if validator.Status == Active {
					allValidators.Status.ActiveValidators++
				}
				if validator.Status == Inactive {
					allValidators.Status.InactiveValidators++
				}
				if validator.Status == Paused {
					allValidators.Status.PausedValidators++
				}
				if validator.Status == Jailed {
					allValidators.Status.JailedValidators++
				}
			}

			allValidators.Status.ConsensusStopped = float64(allValidators.Status.ActiveValidators) < math.Floor(float64(allValidators.Status.TotalValidators)*2/3)+1

			return allValidators, nil, statusCode
		}

		return result, nil, statusCode
	}

	return success, failure, statusCode
}

// QueryValidators is a function to list validators.
func QueryValidators(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !common.RPCMethods["GET"][config.QueryValidators].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryValidators].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[encode-transaction] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryValidatorsHandle(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, true)
	}
}

func queryValidatorInfosHandle(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	queries := r.URL.Query()
	key := queries["key"]
	offset := queries["offset"]
	limit := queries["limit"]
	countTotal := queries["count_total"]
	all := queries["all"]

	var events = make([]string, 0, 9)
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
	if len(all) == 1 {
		events = append(events, fmt.Sprintf("all=%s", all[0]))
	}

	r.URL.RawQuery = strings.Join(events, "&")

	return common.ServeGRPC(r, gwCosmosmux)
}

// QueryValidatorInfos is a function to list validators information.
func QueryValidatorInfos(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !common.RPCMethods["GET"][config.QueryValidatorInfos].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryValidatorInfos].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[encode-transaction] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryValidatorInfosHandle(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, true)
	}
}
