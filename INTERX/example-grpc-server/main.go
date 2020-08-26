package main

import (
	"io/ioutil"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"github.com/KiraCore/sekai/INTERX/insecure"
	pbExample "github.com/KiraCore/sekai/INTERX/proto-gen"
	"github.com/KiraCore/sekai/INTERX/example-grpc-server/server"

	// Static files
	_ "github.com/KiraCore/sekai/INTERX/statik"
)

func main() {
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	addr := "0.0.0.0:10000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	if strings.ToLower(os.Getenv("WITH_TRANSPORT_CREDENTIALS")) == "false" {
		s := grpc.NewServer()
		pbExample.RegisterExampleServiceServer(s, server.New())

		// Serve gRPC Server
		log.Info("Serving gRPC on https://", addr)
		s.Serve(lis)
	} else {
		s := grpc.NewServer(
			// TODO: Replace with your own certificate!
			grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
		)
		pbExample.RegisterExampleServiceServer(s, server.New())
	
		// Serve gRPC Server
		log.Info("Serving gRPC on https://", addr)
		s.Serve(lis)
	}
}
