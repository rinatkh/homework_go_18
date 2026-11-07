.DEFAULT_GOAL := help
GO ?= go
DOCKER ?= docker
COMPOSE=$(DOCKER) compose --env-file .env -f deploy/compose.yaml
PROTOC ?= protoc
TOOLS=$(CURDIR)/.tools
PATH:=$(TOOLS):$(PATH)
export PATH
.PHONY: help env up down build test vet test-race check tools proto worker client logs kafka-consume redis-cli clean
help: ## List starter workflow
	@awk 'BEGIN{FS=":.*## "} /^[a-zA-Z_-]+:.*## /{printf "make %-15s %s\n",$$1,$$2}' $(MAKEFILE_LIST)
env: ## Copy local config once
	@python3 scripts/create-env.py
up: env ## Start ready infra + starter API; unfinished worker/client are opt-in
	$(COMPOSE) up --build -d --wait --wait-timeout 180
down: env ## Keep data volumes
	$(COMPOSE) --profile exercise down
build: ## Compile all starter binaries
	@mkdir -p bin
	@set -e; for name in document-api document-worker document-client healthcheck migrate; do $(GO) build -buildvcs=false -o bin/$$name ./cmd/$$name; done
test: ## Existing tests pass; TODO behavior still needs student tests
	$(GO) test ./...
vet: ## Static checks
	$(GO) vet ./...
test-race: ## Race checks
	$(GO) test -race ./...
check: test vet test-race build ## Starter baseline check
tools: ## Pinned protoc plugins
	@mkdir -p $(TOOLS)
	GOBIN=$(TOOLS) $(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5
	GOBIN=$(TOOLS) $(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
proto: ## Reproduce generated contract
	@mkdir -p pkg/proto
	cd api/proto && $(PROTOC) -I . --go_out=../../pkg/proto --go_opt=paths=source_relative --go-grpc_out=../../pkg/proto --go-grpc_opt=paths=source_relative document/v1/document.proto
worker: env ## TODO-07: complete command to run worker with exercise profile
	@echo 'TODO-07: complete worker target'; exit 1
client: env ## TODO-07: complete command to run one-shot gRPC client
	@echo 'TODO-07: complete client target'; exit 1
logs: env ## Process output
	$(COMPOSE) logs -f --tail=100
kafka-consume: env ## Observe both event types without using application group
	$(COMPOSE) exec kafka /opt/kafka/bin/kafka-console-consumer.sh --bootstrap-server kafka:9092 --topic document.events --from-beginning
redis-cli: env ## Inspect TTL receipts
	$(COMPOSE) exec redis redis-cli
clean: ## Remove local build/tools
	rm -rf bin .tools

fmt: ## Format source files
	gofmt -w $$(find cmd internal pkg -name '*.go')
fmt-check: ## Check formatting
	@test -z "$$(gofmt -l $$(find cmd internal pkg -name '*.go'))"
ci: fmt-check check ## Starter CI without infrastructure or unfinished acceptance behavior
acceptance: ## Intentionally red until TODO implementations are completed
	$(GO) test -tags=acceptance ./internal/worker ./internal/grpcapi ./internal/automation
acceptance-redis: ## Real Redis receipt contract; TEST_REDIS_URL required
	$(GO) test -tags=integration ./internal/repository/redis
migrate-up: env ## Goose Up in Compose
	$(COMPOSE) run --rm migrate /app/migrate up
migrate-down: env ## Goose Down by one migration
	$(COMPOSE) run --rm migrate /app/migrate down
migrate-status: env ## Goose Status
	$(COMPOSE) run --rm migrate /app/migrate status
migrate-local-up: ## Native goose; export DATABASE_URL
	$(GO) run ./cmd/migrate up
migrate-local-down: ## Native goose Down
	$(GO) run ./cmd/migrate down
migrate-local-status: ## Native goose Status
	$(GO) run ./cmd/migrate status
compose-config: env ## Validate Compose only
	$(COMPOSE) config --quiet

.PHONY: PATH help env up down build test vet test-race check tools proto worker client logs kafka-consume redis-cli clean fmt fmt-check ci acceptance acceptance-redis migrate-up migrate-down migrate-status migrate-local-up migrate-local-down migrate-local-status compose-config
