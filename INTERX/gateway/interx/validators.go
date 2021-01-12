package interx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/KiraCore/sekai/INTERX/common"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// RegisterValidatorsQueryRoutes registers validators query routers.
func RegisterValidatorsQueryRoutes(r *mux.Router, gwCosmosmux *runtime.ServeMux, rpcAddr string) {
	r.HandleFunc(common.QueryValidators, QueryValidators(gwCosmosmux, rpcAddr)).Methods("GET")

	common.AddRPCMethod("GET", common.QueryValidators, "This is an API to query validators.", true)
}

func queryValidatorsHandle(r *http.Request, gwCosmosmux *runtime.ServeMux) (interface{}, interface{}, int) {
	queries := r.URL.Query()
	status := queries["status"]
	key := queries["key"]
	offset := queries["offset"]
	limit := queries["limit"]
	countTotal := queries["count_total"]

	var events = make([]string, 0, 5)
	if len(status) == 1 {
		events = append(events, fmt.Sprintf("status=%s", status[0]))
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

	r.URL.RawQuery = strings.Join(events, "&")

	// Todos - need to change to use grpc endpoint
	// return common.ServeGRPC(r, gwCosmosmux)
	type ValidatorsResponse struct {
		Validators []struct {
			Address    string  `json:"address"`
			Valkey     string  `json:"valkey"`
			Pubkey     string  `json:"pubkey"`
			Moniker    string  `json:"moniker"`
			Website    string  `json:"website"`
			Social     string  `json:"social"`
			Identity   string  `json:"identity"`
			Commission sdk.Dec `json:"commission"`
			Status     string  `json:"status"`
			Rank       int64   `json:"rank"`
			Streak     int64   `json:"streak"`
			Mischance  int64   `json:"mischance"`
		} `json:"validators"`
		Pagination struct {
			NextKey []byte `protobuf:"bytes,1,opt,name=next_key,json=nextKey,proto3" json:"next_key,omitempty"`
			Total   uint64 `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
		} `json:"pagination"`
	}

	validatorsJSON := `{
		"validators": [
			{
				"address": "kira1lrh9pcmlnywlpphj8vtm0j0alhrrjwjsdxmjts",
				"valkey": "kiravaloper1lrh9pcmlnywlpphj8vtm0j0alhrrjwjs7q83nu",
				"pubkey": "kiravalconspub1zcjduepqm98ffgul4ppzzur6l67v3mj2vsyc7tr9nrzwk3e0ffx0z7l9sgsqnln467",
				"moniker": "hello",
				"website": "",
				"social": "social",
				"identity": "",
				"commission": "1.000000000000000000",
				"status": "active",
				"rank": 1,
				"streak": 1,
				"mischance": 1
			}
		],
		"pagination": {
			"total": 2
		}
	}`

	response := ValidatorsResponse{}
	if err := json.Unmarshal([]byte(validatorsJSON), &response); err != nil {
		panic(err)
	}

	return response, nil, http.StatusOK
}

// QueryValidators is a function to list validators.
func QueryValidators(gwCosmosmux *runtime.ServeMux, rpcAddr string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := common.GetInterxRequest(r)
		response := common.GetResponseFormat(request, rpcAddr)
		statusCode := http.StatusOK

		response.Response, response.Error, statusCode = queryValidatorsHandle(r, gwCosmosmux)

		common.WrapResponse(w, request, *response, statusCode, false)
	}
}
