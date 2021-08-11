package gateway

import (
	"context"
	"crypto/tls"
	"fmt"
	"mime"
	"net/http"

	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/database"
	"github.com/KiraCore/sekai/INTERX/functions"
	"github.com/KiraCore/sekai/INTERX/insecure"
	cosmosAuth "github.com/KiraCore/sekai/INTERX/proto-gen/cosmos/auth"
	cosmosBank "github.com/KiraCore/sekai/INTERX/proto-gen/cosmos/bank"
	kiraGov "github.com/KiraCore/sekai/INTERX/proto-gen/kira/gov"
	kiraSlashing "github.com/KiraCore/sekai/INTERX/proto-gen/kira/slashing"
	kiraStaking "github.com/KiraCore/sekai/INTERX/proto-gen/kira/staking"
	kiraTokens "github.com/KiraCore/sekai/INTERX/proto-gen/kira/tokens"
	kiraUpgrades "github.com/KiraCore/sekai/INTERX/proto-gen/kira/upgrade"
	"github.com/KiraCore/sekai/INTERX/tasks"
	functionmeta "github.com/KiraCore/sekai/function_meta"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	grpclog "google.golang.org/grpc/grpclog"
)

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

	// With transport credentials
	// if strings.ToLower(os.Getenv("WITH_TRANSPORT_CREDENTIALS")) == "true" {
	// 	security = grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(insecure.CertPool, ""))
	// }

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

	err = cosmosAuth.RegisterQueryHandler(context.Background(), gwCosmosmux, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	err = kiraGov.RegisterQueryHandler(context.Background(), gwCosmosmux, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	err = kiraStaking.RegisterQueryHandler(context.Background(), gwCosmosmux, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	err = kiraSlashing.RegisterQueryHandler(context.Background(), gwCosmosmux, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	err = kiraTokens.RegisterQueryHandler(context.Background(), gwCosmosmux, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	err = kiraUpgrades.RegisterQueryHandler(context.Background(), gwCosmosmux, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	return gwCosmosmux, nil
}

// Run runs the gRPC-Gateway, dialling the provided address.
func Run(configFilePath string, log grpclog.LoggerV2) error {
	config.LoadConfig(configFilePath)
	functions.RegisterInterxFunctions()
	functionmeta.RegisterStdMsgs()

	database.LoadBlockDbDriver()
	database.LoadBlockNanoDbDriver()
	database.LoadFaucetDbDriver()
	database.LoadReferenceDbDriver()

	serveHTTPS := config.Config.ServeHTTPS
	grpcAddr := config.Config.GRPC
	rpcAddr := config.Config.RPC
	port := config.Config.PORT

	gwCosmosmux, err := GetGrpcServeMux(grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	oaHander := getOpenAPIHandler()

	router := mux.NewRouter()
	RegisterRequest(router, gwCosmosmux, rpcAddr)

	router.PathPrefix("/").Handler(oaHander)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodHead, http.MethodOptions, http.MethodPatch, http.MethodConnect, http.MethodTrace},
		AllowCredentials: true,
		ExposedHeaders:   []string{"*"},
	})

	gatewayAddr := "0.0.0.0:" + port
	gwServer := &http.Server{
		Addr:    gatewayAddr,
		Handler: c.Handler(router),
	}

	tasks.RunTasks(gwCosmosmux, rpcAddr, gatewayAddr)

	if serveHTTPS {
		gwServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{insecure.Cert},
		}

		log.Info("Serving gRPC-Gateway and OpenAPI Documentation on https://", gatewayAddr)
		return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServeTLS("", ""))
	}

	log.Info("Serving gRPC-Gateway and OpenAPI Documentation on http://", gatewayAddr)
	return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServe())
}
