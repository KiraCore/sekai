PACKAGES=$(shell go list ./... | grep -v '/simulation')

all: install
install:
	./scripts/build.sh

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

test:
	./scripts/test.sh

# look into .golangci.yml for enabling / disabling linters
lint:
	@echo "--> Running linter"
	@golangci-lint run
	@go mod verify

proto-gen:
	docker run --rm -v $(CURDIR):/workspace --workdir /workspace tendermintdev/sdk-proto-gen sh ./scripts/protocgen.sh

proto-gen-local:
	./scripts/proto-gen.sh

build:
	./scripts/build.sh

start:
	go run ./cmd/sekaid/main.go

# ./scripts/proto-gen.sh
publish:
	./scripts/publish.sh

