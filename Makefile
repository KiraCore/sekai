DOCKER := $(shell which docker)

.PHONY: all install go.sum test test-local lint proto-gen proto-gen-local build start publish

all: install
install:
	./scripts/build.sh

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

test:
	./scripts/test.sh

test-local:
	./scripts/test-local.sh

test-custody:
	./scripts/test-custody.sh

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify

containerProtoVer=0.14.0
containerProtoImage=ghcr.io/cosmos/proto-builder:$(containerProtoVer)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(DOCKER) run --user $(id -u):$(id -g) --rm -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protocgen.sh; 

proto-format:
	@echo "Formatting Protobuf files"
	@$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./proto -name "*.proto" -exec clang-format -i {} \;  

proto-swagger-gen:
	@echo "Generating Protobuf Swagger"
	@$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) \
		sh ./scripts/protoc-swagger-gen.sh; 

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) breaking --against $(HTTPS_GIT)#branch=main

build:
	./scripts/build.sh

start:
	go run ./cmd/sekaid/main.go

# ./scripts/proto-gen.sh
publish:
	./scripts/publish.sh

network-start:
	./scripts/test-local/network-stop.sh
	./scripts/test-local/network-start.sh

network-stop:
	./scripts/test-local/network-stop.sh