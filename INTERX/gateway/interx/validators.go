package interx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/tasks"
	"github.com/KiraCore/sekai/INTERX/types"
	"github.com/KiraCore/sekai/INTERX/types/kira"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	tmRPCTypes "github.com/tendermint/tendermint/rpc/core/types"
)

// RegisterValidatorsQueryRoutes registers validators query routers.
func RegisterValidatorsQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryConsensus, QueryConsensus(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryValidators, QueryValidators(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryValidatorInfos, QueryValidatorInfos(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryConsensus, "This is an API to query consensus.", true)
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

func queryValidatorsHandle(r *http.Request, gwCosmosmux *runtime.ServeMux, rpcAddr string) (interface{}, interface{}, int) {
	queries := r.URL.Query()
	address := queries["address"]
	valkey := queries["valkey"]
	pubkey := queries["pubkey"]
	moniker := queries["moniker"]
	status := queries["status"]
	offset := queries["offset"]
	limit := queries["limit"]
	proposer := queries["proposer"]
	countTotal := queries["count_total"]
	all := queries["all"]

	if len(all) == 1 && all[0] == "true" {
		// Query All Validators
		return tasks.AllValidators, nil, http.StatusOK
	}

	response := struct {
		Validators []types.QueryValidator `json:"validators,omitempty"`
		Pagination struct {
			Total int `json:"total,string,omitempty"`
		} `json:"pagination,omitempty"`
	}{}

	validators := tasks.AllValidators.Validators
	if len(countTotal) == 1 && countTotal[0] == "true" {
		response.Pagination.Total = len(validators)
	}

	from := 0
	count := len(validators)
	if len(offset) == 1 {
		from, _ = strconv.Atoi(offset[0])
	}
	if len(limit) == 1 {
		count, _ = strconv.Atoi(limit[0])
	}

	for index, validator := range validators {
		if from > index || index >= from+count {
			continue
		}
		if len(address) == 1 && validator.Address != address[0] {
			continue
		}
		if len(valkey) == 1 && validator.Valkey != valkey[0] {
			continue
		}
		if len(pubkey) == 1 && validator.Pubkey != pubkey[0] {
			continue
		}
		if len(proposer) == 1 && validator.Proposer != proposer[0] {
			continue
		}
		if len(moniker) == 1 && validator.Moniker != moniker[0] {
			continue
		}
		if len(status) == 1 && validator.Status != status[0] {
			continue
		}

		response.Validators = append(response.Validators, validator)
	}
	return response, nil, http.StatusOK
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

					common.GetLogger().Info("[query-validators] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryValidatorsHandle(r, gwCosmosmux, rpcAddr)
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

					common.GetLogger().Info("[query-validator-info] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryValidatorInfosHandle(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, true)
	}
}

func GetValidator(consAddrHex string) string {
	bytes, err := hex.DecodeString(consAddrHex)
	if err != nil {
		return ""
	}

	return sdk.AccAddress(bytes).String()
}

func queryConsensusHandle(r *http.Request, gwCosmosmux *runtime.ServeMux, rpcAddr string) (interface{}, interface{}, int) {
	// Query All Validators
	AllValidators := tasks.AllValidators.Validators

	var catching_up bool
	success, failure, statusCode := common.MakeTendermintRPCRequest(rpcAddr, "/status", "")
	if success != nil {
		type TempResponse struct {
			SyncInfo struct {
				CatchingUp bool `json:"catching_up"`
			} `json:"sync_info"`
		}
		result := TempResponse{}

		byteData, err := json.Marshal(success)
		if err != nil {
			common.GetLogger().Error("[query-consensus] Invalid response format", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		err = json.Unmarshal(byteData, &result)
		if err != nil {
			common.GetLogger().Error("[query-consensus] Invalid response format", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		catching_up = result.SyncInfo.CatchingUp
	} else {
		return success, failure, statusCode
	}

	// Query consensus
	success, failure, statusCode = common.MakeTendermintRPCRequest(rpcAddr, "/dump_consensus_state", "")
	if success != nil {
		consensus := tmRPCTypes.ResultDumpConsensusState{}

		byteData, err := json.Marshal(success)
		if err != nil {
			common.GetLogger().Error("[query-consensus] Invalid response format: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		err = json.Unmarshal(byteData, &consensus)
		if err != nil {
			common.GetLogger().Error("[query-consensus] Invalid response format: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		roundState := kira.RoundState{}

		err = json.Unmarshal(consensus.RoundState, &roundState)
		if err != nil {
			common.GetLogger().Error("[query-consensus] Invalid round state: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		response := kira.ConsensusResponse{}
		response.Height = roundState.Height
		response.Round = roundState.Round
		response.Step = roundState.Step.String()
		response.StartTime = roundState.StartTime
		response.CommitTime = roundState.CommitTime
		response.TriggeredTimeoutPrecommit = roundState.TriggeredTimeoutPrecommit
		response.ConsensusStopped = common.IsConsensusStopped(len(roundState.Validators.Validators))

		if catching_up {
			response.ConsensusStopped = false
		}

		response.AverageBlockTime = common.GetAverageBlockTime()

		// response.Proposer
		for _, validator := range AllValidators {
			if validator.Proposer == roundState.Validators.Proposer.Address {
				response.Proposer = validator.Address
				break
			}
		}

		validators := make([]string, 0)
		for i := range roundState.Validators.Validators {
			for _, validator := range AllValidators {
				if validator.Proposer == roundState.Validators.Validators[i].Address {
					validators = append(validators, validator.Address)
					break
				}
			}
		}

		response.Precommits = make([]string, 0)
		response.Prevotes = make([]string, 0)
		response.Noncommits = make([]string, 0)

		flag := make([]bool, len(validators))
		for i, vote := range roundState.Votes {
			for j := range vote.Precommits {
				if vote.Precommits[j] != "nil-Vote" {
					flag[j] = true
					if j == len(roundState.Votes)-1 {
						response.Precommits = append(response.Precommits, validators[j])
					}
				}
			}
			if i == len(roundState.Votes)-1 {
				for j := range vote.Prevotes {
					if vote.Prevotes[j] != "nil-Vote" {
						response.Prevotes = append(response.Prevotes, validators[j])
					}
				}
			}
		}

		for i := range flag {
			if !flag[i] {
				response.Noncommits = append(response.Noncommits, validators[i])
			}
		}

		return response, failure, statusCode
	}

	return success, failure, statusCode
}

// QueryConsensus is a function to query consensus.
func QueryConsensus(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		if !common.RPCMethods["GET"][config.QueryConsensus].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryConsensus].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-consensus] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryConsensusHandle(r, gwCosmosmux, rpcAddr)
		}

		common.WrapResponse(w, request, *response, statusCode, true)
	}
}
