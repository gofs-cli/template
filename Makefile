$(shell cp -n .env.example .env)
include .env
export

run: dbup
	@go tool wgo -file=.go -file=.templ -file=.css -file=.ts -xfile=_templ.go ./scripts/build.sh :: ./tmp/main
.PHONY: run

lint:
	@golangci-lint run
.PHONY: lint

build:
	@go build cmd/server/main.go
.PHONY: build

test:
	@go test -v ./...
.PHONY: test

dbup:
	@if [ -n "$$DSN" ]; then \
		docker compose -f docker/docker-compose.yml up -d postgres; \
    fi
.PHONY: dbup

dbdown:
	@docker compose -f docker/docker-compose.yml down
.PHONY: dbdown

gendata: dbdown dbup
	@go test -v -tags=gendata -count=1 ./...
.PHONY: gendata

codegen:
	@go generate ./...
.PHONY: gencode

alpine:
	@echo "Installing alpine.js"
	@cd internal/server/assets/js && { \
		curl -s -O https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js ; mv cdn.min.js alpine.js ; \
		curl -s -O https://cdn.jsdelivr.net/npm/@alpinejs/collapse@3.x.x/dist/cdn.min.js ; mv cdn.min.js collapse.js ; \
		curl -s -O https://cdn.jsdelivr.net/npm/@alpinejs/persist@3.x.x/dist/cdn.min.js ; mv cdn.min.js persist.js ; \
		cd - ; \
	}
.PHONY: alpine

htmx:
	@echo "Installing htmx"
	@cd internal/server/assets/js && { \
			curl -s -O https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js ; \
			curl -s -O https://unpkg.com/htmx-ext-response-targets@2.0.0/response-targets.js ; cd - ; \
		}
.PHONY: htmx

deps: htmx alpine
	@go mod tidy
.PHONY: deps
