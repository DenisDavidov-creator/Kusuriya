# --- run -------------------------------------------------------------------

init: tools .env docs/swagger.json ## first-time setup: install tools, create .env, generate the spec
	@echo "ready — run 'make up'"

up: ## start Postgres, migrate, seed, API (:8080) and Swagger UI (:8090)
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f web

.env:
	cp .env.example .env

# --- test ------------------------------------------------------------------

test:
	go test ./...

test-integration:
	go test -race -tags integration -count=1 ./...

lint:
	golangci-lint run --build-tags=integration

# --- generate --------------------------------------------------------------

generate: ## regenerate sqlc queries and mockery mocks
	sqlc generate
	go generate ./...

swagger: ## regenerate the OpenAPI spec into docs/ and reload Swagger UI
	swag init -g cmd/kusuriya/main.go --parseInternal -o docs
	-docker compose restart swagger-ui 2>/dev/null

docs/swagger.json: $(shell find cmd internal -name '*.go' ! -name '*_test.go')
	swag init -g cmd/kusuriya/main.go --parseInternal -o docs

tools: ## install sqlc, mockery, swag, migrate, golangci-lint
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/vektra/mockery/v2@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
