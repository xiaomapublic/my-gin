# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOFMT=$(GOCMD) fmt
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=my-gin
BINARY_UNIX=$(BINARY_NAME)_unix

all: run

build: ## Build for Linux
	$(GOBUILD) -o ./build/$(BINARY_NAME) -ldflags="-s -w" -tags=jsoniter -v ./

fmt: ## Fmt previous build
	$(GOFMT) $(BINARY_NAME)

clean:  ## Remove previous build
	$(GOCLEAN)
	rm -f ./build/$(BINARY_NAME)
	rm -f ./build/$(BINARY_UNIX)

run: ## Run for Linux
	$(GOBUILD) -o ./build/$(BINARY_NAME) -ldflags="-s -w" -tags=jsoniter -v ./
	./build/$(BINARY_NAME)

restart: ## Restart for Linux
	kill -INT $$(cat pid)
	$(GOBUILD) -o ./build/$(BINARY_NAME) -ldflags="-s -w" -tags=jsoniter -v ./
	./build/$(BINARY_NAME)

deps:
	$(GOGET) github.com/kardianos/govendor
	govendor sync

cross:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/$(BINARY_UNIX) -ldflags="-s -w" -tags=jsoniter -v ./

help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
