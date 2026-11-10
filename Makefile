.PHONY: build build-arm64 build-amd64 test install clean help gui gui-build gui-start

BINARY_NAME=kick-interval
BUILD_DIR=build
CMD_DIR=cmd/kick-interval
ELECTRON_DIR=electron

VERSION ?= 0.1.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: build-arm64 build-amd64 ## Build for all macOS architectures

build-arm64: ## Build for macOS ARM64 (Apple Silicon)
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-arm64 ./$(CMD_DIR)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)-arm64"

build-amd64: ## Build for macOS AMD64 (Intel)
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-amd64 ./$(CMD_DIR)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)-amd64"

build-local: ## Build for current platform
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

install: build-local ## Install to /usr/local/bin
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"

gui-install: ## Install Electron dependencies
	@cd $(ELECTRON_DIR) && npm install

gui: build gui-install ## Build backend and run GUI
	@cd $(ELECTRON_DIR) && npm start

gui-build: build gui-install ## Build backend and package GUI app
	@cd $(ELECTRON_DIR) && npm run build

gui-build-mac: build gui-install ## Build backend and package GUI for macOS
	@cd $(ELECTRON_DIR) && npm run build:mac

gui-build-win: build gui-install ## Build backend and package GUI for Windows
	@cd $(ELECTRON_DIR) && npm run build:win

gui-build-linux: build gui-install ## Build backend and package GUI for Linux
	@cd $(ELECTRON_DIR) && npm run build:linux

clean: ## Remove build artifacts
	@rm -rf $(BUILD_DIR)
	@rm -rf $(ELECTRON_DIR)/dist
	@rm -f coverage.out coverage.html
	@echo "Cleaned build artifacts"

deps: ## Download dependencies
	go mod download
	go mod tidy

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: fmt vet ## Run linters

.DEFAULT_GOAL := help
