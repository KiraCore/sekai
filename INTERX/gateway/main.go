package gateway

import (
	"context"
	"crypto/tls"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/KiraCore/sekai/INTERX/insecure"
	cosmosBank "github.com/KiraCore/sekai/INTERX/proto-gen/cosmos/bank"
	sekaiapp "github.com/KiraCore/sekai/app"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpclog "google.golang.org/grpc/grpclog"
)

var encodingConfig = sekaiapp.MakeEncodingConfig()

// getOpenAPIHandler serves an OpenAPI UI.
func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")

	statikFS, err := fs.New()
	if err != nil {
		// Panic since this is a permanent error.
		panic("creating OpenAPI filesystem: " + err.Error())
	}

	return http.FileServer(statikFS)
}

// GetGrpcServeMux is a function to get ServerMux for GRPC server.
func GetGrpcServeMux(grpcAddr string) (*runtime.ServeMux, error) {
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

// Run runs the gRPC-Gateway, dialling the provided address.
func Run(grpcAddr string, rpcAddr string, log grpclog.LoggerV2) error {
	gwCosmosmux, err := GetGrpcServeMux(grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	oaHander := getOpenAPIHandler()

	// PORT: Empty parameters mean use port 11000.
	port := os.Getenv("PORT")
	if port == "" {
		port = "11000"
	}

	router := mux.NewRouter()
	RegisterRequest(router, gwCosmosmux, rpcAddr)

	router.PathPrefix("/").Handler(oaHander)

	gatewayAddr := "0.0.0.0:" + port
	gwServer := &http.Server{
		Addr:    gatewayAddr,
		Handler: router,
	}

	// SERVE_HTTP: Empty parameters mean use the TLS Config specified with the server.
	if strings.ToLower(os.Getenv("SERVE_HTTP")) == "false" {
		gwServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{insecure.Cert},
		}

		log.Info("Serving gRPC-Gateway and OpenAPI Documentation on https://", gatewayAddr)
		return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServeTLS("", ""))
	}

	log.Info("Serving gRPC-Gateway and OpenAPI Documentation on http://", gatewayAddr)
	return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServe())
}
