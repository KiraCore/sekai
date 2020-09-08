package gateway

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strings"

	grpclog "google.golang.org/grpc/grpclog"

	"github.com/rakyll/statik/fs"
	"github.com/KiraCore/sekai/INTERX/insecure"

	grpcHandler "github.com/KiraCore/sekai/INTERX/handler/grpc-handler"
	rpcHandler "github.com/KiraCore/sekai/INTERX/handler/rpc-handler"
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

// Run runs the gRPC-Gateway, dialling the provided address.
func Run(grpcAddr string, rpcAddr string, log grpclog.LoggerV2) error {
	gwCosmosmux, err := grpcHandler.GetServeMux(grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %w", err)
	}

	oaHander := getOpenAPIHandler()

	// PORT: Empty parameters mean use port 11000.
	port := os.Getenv("PORT")
	if port == "" {
		port = "11000"
	}

	gatewayAddr := "0.0.0.0:" + port
	gwServer := &http.Server{
		Addr: gatewayAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if grpcHandler.ServeGRPC(w,r, gwCosmosmux, rpcAddr) {
				return
			}

			if rpcHandler.ServeRPC(w, r, rpcAddr) {
				return
			}

			oaHander.ServeHTTP(w, r)
		}),
	}
	
	// SERVE_HTTP: Empty parameters mean use the TLS Config specified with the server.
	if strings.ToLower(os.Getenv("SERVE_HTTP")) == "false" {
		gwServer.TLSConfig = &tls.Config {
			Certificates: []tls.Certificate{insecure.Cert},
		}
		
		log.Info("Serving gRPC-Gateway and OpenAPI Documentation on https://", gatewayAddr)
		return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServeTLS("", ""))
	}

	log.Info("Serving gRPC-Gateway and OpenAPI Documentation on http://", gatewayAddr)
	return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServe())
}
