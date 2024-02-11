CWD=$$(pwd)
PKG := "/app"
PKG_LIST := $(shell go list ${PKG}/...)

.PHONY: all dep build test coverage coverhtml lint

all: build

build_image:
	docker build -t dns_server .

build_image_dev:
	docker build -t dns_dev -f Dockerfile.dev .

dev:
	docker exec -it dns_server bash

build: ## Build DNS application
	go build -o dns_server

test: ## Run all go-based tests
	go test -race -coverprofile=coverage.txt -covermode=atomic -v ./...

up:
	docker compose up

down:
	docker compose down

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
