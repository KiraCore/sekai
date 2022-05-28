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

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify

containerProtoVer=v0.2
containerProtoImage=tendermintdev/sdk-proto-gen:$(containerProtoVer)
containerProtoGen=cosmos-sdk-proto-gen-$(containerProtoVer)

proto-gen:
	docker run --rm --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(containerProtoImage) sh ./scripts/protocgen.sh

proto-gen-local:
	./scripts/protogen-local.sh

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