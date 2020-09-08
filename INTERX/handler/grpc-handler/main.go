package grpc_handler

import (
	"context"
    "encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/KiraCore/sekai/INTERX/insecure"
	cosmosBank "github.com/KiraCore/sekai/INTERX/proto-gen/cosmos/bank"
	rpcHandler "github.com/KiraCore/sekai/INTERX/handler/rpc-handler"
)

const (
	QueryBalances 	string = "/api/cosmos/bank/balances/"
	QuerySupply 	string = "/api/cosmos/bank/supply"
)

var Endpoints = []string {
	QueryBalances,
	QuerySupply,
}

// Run runs the gRPC-Gateway, dialling the provided address.
func GetServeMux(grpcAddr string) (*runtime.ServeMux, error) {
	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	// WITH_TRANSPORT_CREDENTIALS: Empty parameters mean set transport security.
	security := grpc.WithInsecure()
	if strings.ToLower(os.Getenv("WITH_TRANSPORT_CREDENTIALS")) == "true" {
		security = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(insecure.CertPool, ""))
	}

	conn, err := grpc.DialContext(
		context.Background(),
		grpcAddr,
		security,
		grpc.WithBlock(),
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}

	gwCosmosmux := runtime.NewServeMux()
	err = cosmosBank.RegisterQueryHandler(context.Background(), gwCosmosmux, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}
	
	return gwCosmosmux, nil
}

func ServeGRPC(w http.ResponseWriter, r *http.Request, gwCosmosmux *runtime.ServeMux, rpcAddr string) bool {
	serve := false

	if strings.HasPrefix(r.URL.Path, QueryBalances) && r.Method == http.MethodGet {
		serve = true

		addr, err := sdk.AccAddressFromBech32(strings.TrimPrefix(r.URL.Path, QueryBalances))
		if err != nil {
			fmt.Println(err)
		}
		r.URL.Path = QueryBalances + fmt.Sprintf("%x", addr)
	} else if strings.HasPrefix(r.URL.Path, QuerySupply) && r.Method == http.MethodGet {
		serve = true
	}

	if serve {
		recorder := httptest.NewRecorder()
		gwCosmosmux.ServeHTTP(recorder, r)
		resp := recorder.Result()

		response := rpcHandler.GetResponseFormat(rpcAddr)

		result := new(interface{})
		if json.NewDecoder(resp.Body).Decode(result) == nil {
			if resp.StatusCode == 200 {
				response.Response = result
			} else {
				response.Error = result;
			}
		}
		
		rpcHandler.CopyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(response)

	}

	return serve
}
