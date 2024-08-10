include .env

MAIN_PACKAGE := ./cmd/api/main.go
BINARY_NAME := neoshare
DB_MIGRATION_DIR := ./internal/database/migrations
DB_CONN_STRING := "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_DATABASE}?sslmode=disable"

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## build: Build the production code
.PHONY: build
build:
	@echo "Building binary..."
	@templ generate
	@pnpx tailwindcss -i cmd/web/assets/css/input.css -o cmd/web/assets/css/style.css
	@go build -o bin/${BINARY_NAME} ${MAIN_PACKAGE}

## dev: Run the code development environment
.PHONY: dev
dev:
	@echo "Running development environment..."
	@go run ${MAIN_PACKAGE}

## watch: Run the application with reloading on file changes
.PHONY: watch
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

## gen-repo: Generate the repository with sql queries provided
.PHONY: gen-repo
gen-repo:
	@sqlc generate -f ./internal/database/sqlc.yaml
	@echo "Repository generated..."

## update: Updates the packages and tidy the modfile
.PHONY: update
update:
	@go get -u ./...
	@go mod tidy -v

# ==================================================================================== #
# Database
# ==================================================================================== #

## db-up: Create DB container
.PHONY: db-up
db-up:
	@if docker compose up -d 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up -d; \
	fi
	@echo "DB container is up and running..."

## db-down: Shutdown DB container
.PHONY: db-down
db-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi
	@echo "DB container is down..."

## migrate-up: Migrate to up the table schema
.PHONY: migrate-up
migrate-up:
	@GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=${DB_MIGRATION_DIR} GOOSE_DBSTRING=${DB_CONN_STRING} goose up

## migrate-down: Migrate to down the table schema
.PHONY: migrate-down
migrate-down:
	@GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=${DB_MIGRATION_DIR} GOOSE_DBSTRING=${DB_CONN_STRING} goose down

## migrate-status: Displays about the migration status for the current DB
.PHONY: migrate-status
migrate-status:
	@GOOSE_DRIVER=postgres GOOSE_MIGRATION_DIR=${DB_MIGRATION_DIR} GOOSE_DBSTRING=${DB_CONN_STRING} goose status

## sqlc-gen: Generate the sqlc queries
.PHONY: sqlc-gen
sqlc-gen:
	@sqlc generate -f ./internal/database/sqlc.yaml

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: Format code and tidy modfile
.PHONY: tidy
tidy:
	@echo "Tidying up..."
	@go fmt ./...
	@go mod tidy -v

## lint: Run linter
.PHONY: lint
lint:
	@echo "Linting..."
	@golangci-lint run

## test: Test the application
.PHONY: test
test:
	@echo "Testing..."
	@go test ./tests -v

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## clean: Clean the binary
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -f bin/

## help: Print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
