package main

import (
	"flag"
	"io/ioutil"
	"mime"
	"net/http"
	"os"

	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc/grpclog"

	"github.com/KiraCore/sekai/INTERX/gateway"

	// Static files
	_ "github.com/KiraCore/sekai/INTERX/statik"
)

// getOpenAPIHandler serves an OpenAPI UI.
func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")

	statikFS, err := fs.New()
	if err != nil {
		panic("creating OpenAPI filesystem: " + err.Error())
	}

	return http.FileServer(statikFS)
}

var serverAddress = flag.String(
	"server-address",
	"dns:///0.0.0.0:10000",
	"The address to the gRPC server, in the gRPC standard naming format. "+
		"See https://github.com/grpc/grpc/blob/master/doc/naming.md for more information.",
)

func main() {
	flag.Parse()

	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	err := gateway.Run(*serverAddress)
	log.Fatalln(err)
}
