.PHONY: build run
help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

run: ## Run the program
	@go run ./main.go $@

build: ## Build the program
	@CGO_ENABLED=0 go build -ldflags "-s -w" -v main.go
