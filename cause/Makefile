# Makefile for Error Cause Package

.PHONY: help test test-verbose test-race test-cover bench lint fmt vet check clean doc serve-doc

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Testing
test: ## Run all tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-race: ## Run tests with race detection
	go test -race ./...

test-cover: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-cover-func: ## Show coverage by function
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

##@ Benchmarking
bench: ## Run benchmarks
	go test -bench=. -benchmem

bench-cpu: ## Run CPU benchmarks with profiling
	go test -bench=. -cpuprofile=cpu.prof
	@echo "CPU profile generated: cpu.prof"

bench-mem: ## Run memory benchmarks with profiling
	go test -bench=. -memprofile=mem.prof
	@echo "Memory profile generated: mem.prof"

bench-txt: ## Run benchmarks and output in txt format (bench.txt)
	go test -bench=. -benchmem -run=^$ > bench.txt
	@echo "Benchmark txt output written to bench.txt"

benchstat: ## Compare two benchmark outputs using benchstat (usage: make benchstat OLD=bench_old.txt NEW=bench_new.txt)
	@if [ -z "$(OLD)" ] || [ -z "$(NEW)" ]; then \
		echo "Usage: make benchstat OLD=old.txt NEW=new.txt"; \
		exit 1; \
	fi
	benchstat $(OLD) $(NEW)

##@ Code Quality
lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	gofmt -s -w .

vet: ## Run go vet
	go vet ./...

check: fmt vet lint test ## Run all checks (format, vet, lint, test)

##@ Development
tidy: ## Tidy go modules
	go mod tidy

download: ## Download dependencies
	go mod download

verify: ## Verify dependencies
	go mod verify

##@ Documentation
doc: ## Generate and serve documentation
	godoc -http=:6060 -play
	@echo "Documentation server running at http://localhost:6060/pkg/github.com/alextanhongpin/errors/cause/"

serve-doc: ## Serve documentation (alias for doc)
	$(MAKE) doc

##@ Examples
run-examples: ## Run example code
	go run examples.go

##@ Cleanup
clean: ## Clean build artifacts
	go clean ./...
	rm -f coverage.out coverage.html
	rm -f cpu.prof mem.prof
	rm -f *.test

##@ Release
tag: ## Create a git tag (usage: make tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then echo "Please specify VERSION: make tag VERSION=v1.0.0"; exit 1; fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)

##@ Installation
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install golang.org/x/tools/cmd/godoc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

##@ CI/CD
ci: check bench ## Run CI pipeline (checks and benchmarks)

##@ Build Info
info: ## Show build information
	@echo "Go version: $(shell go version)"
	@echo "Go env GOOS: $(shell go env GOOS)"
	@echo "Go env GOARCH: $(shell go env GOARCH)"
	@echo "Git commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
	@echo "Git branch: $(shell git branch --show-current 2>/dev/null || echo 'unknown')"

##@ Performance
profile-cpu: ## Profile CPU usage
	go test -cpuprofile=cpu.prof -bench=.
	go tool pprof cpu.prof

profile-mem: ## Profile memory usage
	go test -memprofile=mem.prof -bench=.
	go tool pprof mem.prof

profile-block: ## Profile blocking operations
	go test -blockprofile=block.prof -bench=.
	go tool pprof block.prof

##@ Security
sec-scan: ## Run security scanner (requires gosec)
	@command -v gosec >/dev/null 2>&1 || { echo "gosec not installed. Run: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; exit 1; }
	gosec ./...

##@ Dependencies
deps-list: ## List all dependencies
	go list -m all

deps-graph: ## Show dependency graph (requires graphviz)
	go mod graph | dot -T svg -o deps.svg
	@echo "Dependency graph generated: deps.svg"

deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy

deps-check: ## Check for outdated dependencies
	go list -u -m all

##@ Git Hooks
install-hooks: ## Install git hooks
	@echo "Installing git hooks..."
	@mkdir -p .git/hooks
	@echo '#!/bin/sh\nmake check' > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Pre-commit hook installed"

##@ Misc
todo: ## Find TODO comments in code
	@grep -rn "TODO\|FIXME\|XXX" --include="*.go" . || echo "No TODOs found"

lines: ## Count lines of code
	@find . -name "*.go" -not -path "./vendor/*" | xargs wc -l | tail -1
