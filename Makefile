GO ?= go
COVERAGE_MIN ?= 80
GOLANGCI_LINT_VERSION ?= v2.12.2
GOVULNCHECK_VERSION ?= v1.6.0
CORPUS_COMMIT ?= 5301155192d91d74d337899cecc59dbda59c4c17
VERSION ?= dev
TOOLS_BIN := $(CURDIR)/.tools/bin
GOLANGCI_LINT := $(TOOLS_BIN)/golangci-lint
GOVULNCHECK := $(TOOLS_BIN)/govulncheck
CORPUS_CACHE ?= $(CURDIR)/.cache/basic-computer-games/$(CORPUS_COMMIT)
RELEASE_LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build check clean corpus-fetch corpus-playable corpus-smoke coverage-check fmt fmt-check fuzz lint release-build release-check test tools vet vuln

build:
	mkdir -p bin
	$(GO) build -o bin/go-basic ./cmd/go-basic

release-build:
	rm -rf dist
	mkdir -p dist/staging
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -trimpath -ldflags "$(RELEASE_LDFLAGS)" -o dist/staging/go-basic ./cmd/go-basic
	tar -C dist/staging -czf dist/go-basic_$(VERSION)_linux_amd64.tar.gz go-basic
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -trimpath -ldflags "$(RELEASE_LDFLAGS)" -o dist/staging/go-basic ./cmd/go-basic
	tar -C dist/staging -czf dist/go-basic_$(VERSION)_linux_arm64.tar.gz go-basic
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -trimpath -ldflags "$(RELEASE_LDFLAGS)" -o dist/staging/go-basic ./cmd/go-basic
	tar -C dist/staging -czf dist/go-basic_$(VERSION)_darwin_amd64.tar.gz go-basic
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -trimpath -ldflags "$(RELEASE_LDFLAGS)" -o dist/staging/go-basic ./cmd/go-basic
	tar -C dist/staging -czf dist/go-basic_$(VERSION)_darwin_arm64.tar.gz go-basic
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags "$(RELEASE_LDFLAGS)" -o dist/staging/go-basic.exe ./cmd/go-basic
	zip -qj dist/go-basic_$(VERSION)_windows_amd64.zip dist/staging/go-basic.exe
	rm -rf dist/staging
	cd dist && sha256sum *.tar.gz *.zip > SHA256SUMS

release-check: release-build
	cd dist && sha256sum --check SHA256SUMS

fmt:
	$(GO) fmt ./...

fmt-check:
	@output="$$(find cmd internal pkg test -type f -name '*.go' -exec gofmt -l {} +)"; \
	test -z "$$output" || { echo "Unformatted Go files:"; echo "$$output"; exit 1; }

corpus-fetch:
	$(GO) run ./cmd/corpus-fetch -commit $(CORPUS_COMMIT) -target $(CORPUS_CACHE)

corpus-smoke: corpus-fetch
	$(GO) run ./cmd/corpus-smoke -corpus $(CORPUS_CACHE) -commit $(CORPUS_COMMIT)

corpus-playable:
	$(GO) test -count=1 ./test -run '^TestCLI$$'

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
	rm -rf bin dist .tools
