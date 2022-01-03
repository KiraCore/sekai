package kira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/types"
	govTypes "github.com/KiraCore/sekai/INTERX/types/kira/gov"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterKiraGovProposalRoutes registers kira gov proposal query routers.
func RegisterKiraGovProposalRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(config.QueryProposals, QueryProposalsRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryProposal, QueryProposalRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryVoters, QueryVotersRequest(gwCosmosmux, rpcAddr)).Methods("GET")
	r.HandleFunc(config.QueryVotes, QueryVotesRequest(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", config.QueryProposals, "This is an API to query all proposals.", true)
	common.AddRPCMethod("GET", config.QueryProposal, "This is an API to query a proposal by a given id.", true)
	common.AddRPCMethod("GET", config.QueryVoters, "This is an API to query voters by a given proposal_id.", true)
	common.AddRPCMethod("GET", config.QueryVotes, "This is an API to query votes by a given proposal_id.", true)
}

func queryProposalsHandler(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	queries := r.URL.Query()
	voter := queries["voter"]
	key := queries["key"]
	offset := queries["offset"]
	limit := queries["limit"]
	countTotal := queries["count_total"]
	reverse := queries["reverse"]

	var events = make([]string, 0, 6)
	if len(voter) == 1 {
		events = append(events, fmt.Sprintf("voter=%s", voter[0]))
	}
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
	if len(reverse) == 1 {
		events = append(events, fmt.Sprintf("reverse=%s", reverse[0]))
	}

	r.URL.RawQuery = strings.Join(events, "&")

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

					common.GetLogger().Info("[query-proposal] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryProposalHandler(r, gwCosmosmux, proposalID)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryProposal].CachingEnabled)
	}
}

func queryVotersHandler(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	success, failure, statusCode := common.ServeGRPC(r, gwCosmosmux)

	if success != nil {
		result := struct {
			Voters []struct {
				Address     []byte                `json:"address,omitempty"`
				Roles       []string              `json:"roles,omitempty"`
				Status      string                `json:"status,omitempty"`
				Votes       []string              `json:"votes,omitempty"`
				Permissions *govTypes.Permissions `json:"permissions,omitempty"`
				Skin        uint64                `json:"skin,string,omitempty"`
			} `json:"voters,omitempty"`
		}{}

		byteData, err := json.Marshal(success)
		if err != nil {
			common.GetLogger().Error("[query-voters] Invalid response format: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		err = json.Unmarshal(byteData, &result)
		if err != nil {
			common.GetLogger().Error("[query-voters] Invalid response format: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		voters := make([]govTypes.Voter, 0)

		for _, voter := range result.Voters {
			newVoter := govTypes.Voter{}

			newVoter.Address = sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), voter.Address)
			newVoter.Roles = voter.Roles
			newVoter.Status = voter.Status
			newVoter.Votes = voter.Votes

			newVoter.Permissions.Blacklist = make([]string, 0)
			for _, black := range voter.Permissions.Blacklist {
				newVoter.Permissions.Blacklist = append(newVoter.Permissions.Blacklist, govTypes.PermValue_name[int32(black)])
			}
			newVoter.Permissions.Whitelist = make([]string, 0)
			for _, white := range voter.Permissions.Whitelist {
				newVoter.Permissions.Whitelist = append(newVoter.Permissions.Whitelist, govTypes.PermValue_name[int32(white)])
			}

			newVoter.Skin = voter.Skin

			voters = append(voters, newVoter)
		}

		success = voters
	}

	return success, failure, statusCode
}

// QueryVotersRequest is a function to voters by a given proposal_id.
func QueryVotersRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		proposalID := queries["proposal_id"]
		request := types.InterxRequest{
			Method:   r.Method,
			Endpoint: config.QueryVoters,
			Params:   []byte(proposalID),
		}
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-voters] Entering proposal query by proposal_id: ", proposalID)

		if !common.RPCMethods["GET"][config.QueryVoters].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryVoters].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-voters] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryVotersHandler(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryVoters].CachingEnabled)
	}
}

func queryVotesHandler(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	success, failure, statusCode := common.ServeGRPC(r, gwCosmosmux)

	if success != nil {
		result := struct {
			Votes []struct {
				ProposalID uint64 `json:"proposal_id,string,omitempty"`
				Voter      []byte `json:"voter,omitempty"`
				Option     string `json:"option,omitempty"`
			} `json:"votes,omitempty"`
		}{}

		byteData, err := json.Marshal(success)
		if err != nil {
			common.GetLogger().Error("[query-votes] Invalid response format: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		err = json.Unmarshal(byteData, &result)
		if err != nil {
			common.GetLogger().Error("[query-votes] Invalid response format: ", err)
			return common.ServeError(0, "", err.Error(), http.StatusInternalServerError)
		}

		votes := make([]govTypes.Vote, 0)

		for _, vote := range result.Votes {
			newVote := govTypes.Vote{}

			newVote.ProposalID = vote.ProposalID
			newVote.Voter = sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), vote.Voter)
			newVote.Option = vote.Option

			votes = append(votes, newVote)
		}

		success = votes
	}

	return success, failure, statusCode
}

// QueryVotesRequest is a function to votes by a given proposal_id.
func QueryVotesRequest(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		proposalID := queries["proposal_id"]
		request := types.InterxRequest{
			Method:   r.Method,
			Endpoint: config.QueryVotes,
			Params:   []byte(proposalID),
		}
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		common.GetLogger().Info("[query-votes] Entering proposal query by proposal_id: ", proposalID)

		if !common.RPCMethods["GET"][config.QueryVotes].Enabled {
			response.Response, response.Error, statusCode = common.ServeError(0, "", "API disabled", http.StatusForbidden)
		} else {
			if common.RPCMethods["GET"][config.QueryVotes].CachingEnabled {
				found, cacheResponse, cacheError, cacheStatus := common.SearchCache(request, response)
				if found {
					response.Response, response.Error, statusCode = cacheResponse, cacheError, cacheStatus
					common.WrapResponse(w, request, *response, statusCode, false)

					common.GetLogger().Info("[query-votes] Returning from the cache")
					return
				}
			}

			response.Response, response.Error, statusCode = queryVotesHandler(r, gwCosmosmux)
		}

		common.WrapResponse(w, request, *response, statusCode, common.RPCMethods["GET"][config.QueryVotes].CachingEnabled)
	}
}
