package hosting

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	interx "github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/insecure"
	tasks "github.com/KiraCore/sekai/INTERX/tasks"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	grpclog "google.golang.org/grpc/grpclog"
)

// Run runs the gRPC-Gateway, dialling the provided address.
func Run(log grpclog.LoggerV2) {
	tasks.RunHostingTasks()

	// PORT: Empty parameters mean use port 11000.
	port := os.Getenv("HOSTING_PORT")
	if port == "" {
		port = "12000"
	}

	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(DownloadReference()).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
	})

	hostingAddr := "0.0.0.0:" + port
	gwServer := &http.Server{
		Addr:    hostingAddr,
		Handler: c.Handler(router),
	}

	// SERVE_HTTP: Empty parameters mean use the TLS Config specified with the server.
	if strings.ToLower(os.Getenv("SERVE_HTTP")) == "false" {
		gwServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{insecure.Cert},
		}

		log.Info("File hosting server running on https://", hostingAddr)
		fmt.Println("file hosting server: %w", gwServer.ListenAndServeTLS("", ""))
	} else {
		log.Info("File hosting server running on http://", hostingAddr)
		fmt.Println("file hosting server: %w", gwServer.ListenAndServe())
	}
}

// DownloadReference is a function to download reference.
func DownloadReference() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Path
		if len(filename) != 0 {
			http.ServeFile(w, r, interx.GetReferenceCacheDir()+filename)
		}
	}
}
