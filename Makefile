test:
	@go test -v ./...

tidy:  ## Get the dependencies
	@go mod tidy

compile: tidy  ## compiles server code
	@go build ./...

test-race-cond:
	@go test -v -race ./...

run: compile ## Run the program
	@go run cmd/filestore/*.go

coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic
	@cat cover.out >> coverage.txt

check: ## runs code linting and formatting
	@golangci-lint run --disable=typecheck ./...

.PHONY: fix
fix:
	for file in `golangci-lint --max-same-issues=1000 --max-issues-per-linter=0 run ./...|grep 'goimports'|cut -f 1 -d:`; do	\
		goimports -local "github.com/hungaikev/rdiff" -w $$file;	\
	done


help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: compile test test-race-cond check tidy run coverage help