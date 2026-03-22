BINARY    := investcopilot
CMD       := ./cmd/investcopilot
BUILD_DIR := ./bin

GOFLAGS := -trimpath
LDFLAGS := -s -w

.PHONY: all build lint test clean

all: lint test build

## build: compile the binary to ./bin/investcopilot
build:
	@mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) $(CMD)

## lint: run go vet + staticcheck (install: go install honnef.co/go/tools/cmd/staticcheck@latest)
lint:
	go vet ./...
	staticcheck ./...

## test: run all tests with race detector
test:
	go test -race -count=1 ./...

## test/cover: run tests and open coverage report
test/cover:
	go test -race -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## tidy: tidy and verify go modules
tidy:
	go mod tidy
	go mod verify

## clean: remove build artifacts
clean:
	rm -rf $(BUILD_DIR) coverage.out

## help: list available targets
help:
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sort
