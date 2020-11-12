package main

import (
	"flag"
	"io/ioutil"
	"os"

	"google.golang.org/grpc/grpclog"

	"github.com/KiraCore/sekai/INTERX/gateway"
	config "github.com/KiraCore/sekai/INTERX/config"

	// Static files
	_ "github.com/KiraCore/sekai/INTERX/statik"
)

var serverGRPCAddress = flag.String(
	"server-gRPC-address",
	"dns:///0.0.0.0:9090",
	"The address to the gRPC server, in the gRPC standard naming format. ",
)

var serverRPCAddress = flag.String(
	"server-RPC-address",
	"http://0.0.0.0:26657",
	"The address to the RPC server, in the RPC standard naming format. ",
)

func main() {
	flag.Parse()

	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	config.SetConfig();
	
	grpcAddr := os.Getenv("GRPC")
	if len(grpcAddr) == 0 {
		grpcAddr = *serverGRPCAddress
	}
	
	rpcAddr := os.Getenv("RPC")
	if len(rpcAddr) == 0 {
		rpcAddr = *serverRPCAddress
	}

	err := gateway.Run(grpcAddr, rpcAddr, log)
	log.Fatalln(err)
}
