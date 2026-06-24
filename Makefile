COVERAGE=/tmp/struck/coverage.out

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage: '
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# TESTING
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test ./...

## cover: run tests and generate a coverage profile file
.PHONY: cover
cover:
	@mkdir -p $(dir $(COVERAGE))
	go test ./... -coverprofile=$(COVERAGE)

## cover/func: display coverage breakdown per function
.PHONY: cover/func
cover/func: cover
	go tool cover -func=$(COVERAGE)

## cover/html: generate and open HTML coverage report
.PHONY: cover/html
cover/html: cover
	go tool cover -html=$(COVERAGE)

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: tidy module dependencies, and format and modernize all .go files
.PHONY: tidy
tidy:
	go mod tidy
	go mod verify
	# go mod vendor
	go fix ./...
	go fmt ./...

## audit: run quality control checks
.PHONY: audit
audit:
	go mod tidy -diff
	go mod verify
	go vet ./...
	go tool staticcheck ./...
	go test -race -vet=off ./...
	go run ./cmd/examplecheck
