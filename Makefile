GO ?= go
COVERAGE_MIN ?= 80
GOLANGCI_LINT_VERSION ?= v2.12.2
GOVULNCHECK_VERSION ?= v1.6.0
TOOLS_BIN := $(CURDIR)/.tools/bin
GOLANGCI_LINT := $(TOOLS_BIN)/golangci-lint
GOVULNCHECK := $(TOOLS_BIN)/govulncheck

.PHONY: build check clean coverage-check fmt fmt-check fuzz lint test tools vet vuln

build:
	mkdir -p bin
	$(GO) build -o bin/go-basic ./cmd/go-basic

fmt:
	$(GO) fmt ./...

fmt-check:
	@output="$$(find cmd pkg test -type f -name '*.go' -exec gofmt -l {} +)"; \
	test -z "$$output" || { echo "Unformatted Go files:"; echo "$$output"; exit 1; }

test:
	$(GO) test -count=1 -race -covermode=atomic -coverprofile=coverage.out ./...

coverage-check: test
	@total="$$( $(GO) tool cover -func=coverage.out | awk '/^total:/ { gsub("%", "", $$3); print $$3 }' )"; \
	awk -v total="$$total" -v minimum="$(COVERAGE_MIN)" 'BEGIN { \
		if (total + 0 < minimum + 0) { \
			printf "coverage %.1f%% is below %.1f%%\n", total, minimum; exit 1 \
		} \
		printf "coverage %.1f%% meets %.1f%% minimum\n", total, minimum \
	}'

vet:
	$(GO) vet ./...

$(GOLANGCI_LINT):
	mkdir -p $(TOOLS_BIN)
	GOBIN=$(TOOLS_BIN) $(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

$(GOVULNCHECK):
	mkdir -p $(TOOLS_BIN)
	GOBIN=$(TOOLS_BIN) $(GO) install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)

tools: $(GOLANGCI_LINT) $(GOVULNCHECK)

lint: $(GOLANGCI_LINT)
	mkdir -p $(TOOLS_BIN)/cache
	GOLANGCI_LINT_CACHE=$(TOOLS_BIN)/cache $(GOLANGCI_LINT) run ./...

vuln: $(GOVULNCHECK)
	$(GOVULNCHECK) ./...

fuzz:
	$(GO) test -run=^$$ -fuzz=FuzzLexerTerminates -fuzztime=10s ./pkg/interpreter
	$(GO) test -run=^$$ -fuzz=FuzzParserDoesNotPanic -fuzztime=10s ./pkg/interpreter

check: fmt-check vet coverage-check lint

clean:
	$(GO) clean
	rm -f coverage.out
	rm -rf bin .tools
