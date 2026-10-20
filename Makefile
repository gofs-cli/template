$(shell cp -n .env.example .env)
include .env
export

# Development tools, pinned. They run with `go run pkg@version`, which keeps
# go.mod free of tool dependencies. To use Go's tool directive instead, run
# `go get -tool <pkg>@<version>` and replace `go run $(X)` with `go tool x`.
SQLC        := github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0
WGO         := github.com/bokwoon95/wgo@v0.5.11
STATICCHECK := honnef.co/go/tools/cmd/staticcheck@2025.1.1
GOVULNCHECK := golang.org/x/vuln/cmd/govulncheck@v1.1.4

run: dbup
	@go run $(WGO) -file=.go -file=.tmpl -file=.css -file=.js -file=.sql go run ./cmd/server
.PHONY: run

build:
	@go build -o ./tmp/server ./cmd/server
.PHONY: build

test:
	@go test ./cmd/... ./internal/...
.PHONY: test

lint:
	@go vet ./cmd/... ./internal/...
	@go run $(STATICCHECK) ./cmd/... ./internal/...
	@go run $(GOVULNCHECK) ./cmd/... ./internal/...
.PHONY: lint

fmt:
	@gofmt -l -w cmd internal
.PHONY: fmt

# Runs sqlc for every package with a queries.sql (see sqlc.yaml).
generate:
	@go run $(SQLC) generate
.PHONY: generate

dbup:
	@if [ -n "$$DSN" ]; then \
		docker compose -f docker/docker-compose.yml up -d postgres; \
	fi
.PHONY: dbup

dbdown:
	@docker compose -f docker/docker-compose.yml down
.PHONY: dbdown

gendata: dbdown dbup
	@go test -v -tags=gendata -count=1 ./internal/...
.PHONY: gendata
