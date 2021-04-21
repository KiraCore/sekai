GOBIN=${HOME}/go/bin

generate:
	# Generate go, gRPC-Gateway, OpenAPI output.
	#
	# -I declares import folders, in order of importance
	# This is how proto resolves the protofile imports.
	# It will check for the protofile relative to each of these
	# folders and use the first one it finds.
	#
	# --go_out generates go Protobuf output with gRPC plugin enabled.
	# 		paths=source_relative means the file should be generated
	# 		relative to the input proto file.
	# --grpc-gateway_out generates gRPC-Gateway output.
	# --openapiv2_out generates an OpenAPI 2.0 specification for our gRPC-Gateway endpoints.
	#
	docker run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen.sh

	# Generate static assets for OpenAPI UI
	# statik -m -f -src third_party/OpenAPI/

install: go.sum
		go build -o $(GOBIN)/interxd

start:
	go run main.go
