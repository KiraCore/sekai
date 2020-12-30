package main

import (
	"flag"
	"os"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/gateway"
	_ "github.com/KiraCore/sekai/INTERX/statik"
	"google.golang.org/grpc/grpclog"
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
	log := common.GetLogger()
	grpclog.SetLoggerV2(log)

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
