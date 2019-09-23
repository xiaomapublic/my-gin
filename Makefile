# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOFMT=$(GOCMD) fmt
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=my-gin
BINARY_LINUX=$(BINARY_NAME)-linux-amd64
BINARY_UNIX=$(BINARY_NAME)-unix

all: run

.PHONY:build

build: ## Build for Linux
	$(GOBUILD) -o $(BINARY_LINUX) -ldflags="-s -w" -tags=jsoniter -v ./

fmt: ## Fmt previous build
	$(GOFMT) $(BINARY_NAME)

clean:  ## Remove previous build
	$(GOCLEAN)
	rm -f $(BINARY_LINUX)
	rm -f $(BINARY_UNIX)

run: ## Run for Linux
	$(GOBUILD) -o $(BINARY_LINUX) -ldflags="-s -w" -tags=jsoniter -v ./
	./$(BINARY_LINUX)

restart: ## Restart for Linux
	kill -INT $$(cat pid)
	$(GOBUILD) -o $(BINARY_LINUX) -ldflags="-s -w" -tags=jsoniter -v ./
	./$(BINARY_LINUX)

deps:
	$(GOGET) github.com/kardianos/govendor
	govendor sync

cross:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/$(BINARY_UNIX) -ldflags="-s -w" -tags=jsoniter -v ./

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
