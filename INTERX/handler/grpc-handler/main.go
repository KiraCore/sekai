package interx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	handler "github.com/KiraCore/sekai/INTERX/handler"
	"github.com/KiraCore/sekai/INTERX/insecure"
	cosmosBank "github.com/KiraCore/sekai/INTERX/proto-gen/cosmos/bank"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GetServeMux is a function to get ServerMux for GRPC server.
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

// ServeGRPC is a function to server GRPC
func ServeGRPC(w http.ResponseWriter, r *http.Request, gwCosmosmux *runtime.ServeMux, rpcAddr string) bool {
	serve := false

	if strings.HasPrefix(r.URL.Path, handler.QueryBalances) && r.Method == http.MethodGet {
		serve = true

		addr, err := sdk.AccAddressFromBech32(strings.TrimPrefix(r.URL.Path, handler.QueryBalances))
		if err != nil {
			fmt.Println(err)
		}

		r.URL.Path = handler.QueryBalances + fmt.Sprintf("%x", addr)
	} else if strings.HasPrefix(r.URL.Path, handler.QuerySupply) && r.Method == http.MethodGet {
		serve = true
	}

	if serve {
		recorder := httptest.NewRecorder()
		gwCosmosmux.ServeHTTP(recorder, r)
		resp := recorder.Result()

		response := handler.GetResponseFormat(rpcAddr)

		result := new(interface{})
		if json.NewDecoder(resp.Body).Decode(result) == nil {
			if resp.StatusCode == 200 {
				response.Response = result
			} else {
				response.Error = result
			}
		}

		handler.CopyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)

		handler.WrapResponse(w, *response)
	}

	return serve
}
