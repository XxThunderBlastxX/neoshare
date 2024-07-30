MAIN_PACKAGE := ./cmd/api/main.go
BINARY_NAME := neoshare

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## build: Build the production code
.PHONY: build
build:
	@echo "Building binary..."
	@templ generate
	@pnpx tailwindcss -i cmd/web/assets/css/input.css -o cmd/web/assets/css/output.css
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

## update: Updates the packages and tidy the modfile
.PHONY: update
update:
	@go get -u ./...
	@go mod tidy -v


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
