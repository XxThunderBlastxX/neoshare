MAIN_PACKAGE := ./cmd/api/main.go
BINARY_NAME := neoshare

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## build: Build the production code
.PHONY: build
build:
	@echo "Building binary..."
	@TEMPL_EXPERIMENT=rawgo templ generate
	@pnpx tailwindcss -i cmd/web/assets/css/input.css -o cmd/web/assets/css/style.css
	@go build -o bin/$(BINARY_NAME) $(MAIN_PACKAGE)


## dev: Run the code development environment
.PHONY: dev
dev:
	@echo "Running development environment..."
	@go run $(MAIN_PACKAGE)

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
# DOCKER FOR DB
# ==================================================================================== #

## Create DB container
.PHONY: docker-up
docker-run:
	@if docker compose up -d 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up -d; \
	fi
	@echo "DB container is up and running..."

## Shutdown DB container
.PHONY: docker-down
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi
	@echo "DB container is down..."


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
