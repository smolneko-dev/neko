.DEFAULT_GOAL = build

.PHONY: help
help: ## Display help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: lint
lint: ## Run linter
	golangci-lint run -v ./...

.PHONY: build
build: ## Prepare binary file
	go build -C ./cmd/app/ -o neko

.PHONY: run
run: ## Run app
	go run ./cmd/app/main.go

.PHONY: test
test: ## Run all tests
	go test -v -race ./...

.PHONY: clean
clean: ## Delete binary file
	-rm -f ./cmd/app/neko
