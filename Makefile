.PHONY: help docker-up docker-down docker-build docker-restart docker-logs \
        test test-verbose test-cover test-html run build lint clean

APP_NAME    := go-google-cloud-storage
APP_BINARY  := server
DOCKER_APP  := gcs-app
DOCKER_DB   := gcs-postgres

# ──────────────────────────────────────────────
#  Help
# ──────────────────────────────────────────────
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ──────────────────────────────────────────────
#  Docker
# ──────────────────────────────────────────────
docker-up: ## Start all services with Docker Compose
	@echo "🚀 Starting services..."
	docker compose up -d
	@echo "✅ App running at http://localhost:4000"

docker-down: ## Stop and remove all Docker services
	@echo "🛑 Stopping services..."
	docker compose down
	@echo "✅ Services stopped"

docker-down-volumes: ## Stop services and remove volumes (WARNING: deletes DB data)
	@echo "⚠️  Stopping services and removing volumes..."
	docker compose down -v
	@echo "✅ Services stopped, volumes removed"

docker-build: ## Rebuild Docker images (no cache)
	@echo "🔨 Rebuilding Docker images..."
	docker compose build --no-cache
	@echo "✅ Build complete"

docker-restart: docker-down docker-up ## Restart all Docker services

docker-logs: ## Tail logs from app container
	docker compose logs -f app

docker-logs-all: ## Tail logs from all containers
	docker compose logs -f

docker-ps: ## Show running containers status
	docker compose ps

docker-shell: ## Open shell inside the app container
	docker compose exec app sh

docker-db-shell: ## Open PostgreSQL shell
	docker compose exec postgres psql -U postgres -d gcs

# ──────────────────────────────────────────────
#  Development
# ──────────────────────────────────────────────
run: ## Run the app locally (requires PostgreSQL & .env)
	@echo "🏃 Running app..."
	go run main.go

build: ## Build the binary locally
	@echo "🔨 Building $(APP_NAME)..."
	go build -o $(APP_BINARY) .
	@echo "✅ Binary built: ./$(APP_BINARY)"

watch: ## Run app with auto-reload (requires 'air' tool)
	@which air > /dev/null || (echo "❌ air not installed. Run: go install github.com/air-verse/air@latest" && exit 1)
	air

# ──────────────────────────────────────────────
#  Testing
# ──────────────────────────────────────────────
test: ## Run all unit tests
	@echo "🧪 Running tests..."
	go test ./test/... -count=1

test-verbose: ## Run all tests with verbose output
	@echo "🧪 Running tests (verbose)..."
	go test ./test/... -v -count=1

test-cover: ## Run all tests with coverage report
	@echo "🧪 Running tests with coverage..."
	go test ./test/... -cover -coverpkg=./app/... -coverprofile=coverage.out -count=1
	@echo ""
	@echo "📊 Coverage by function:"
	@go tool cover -func=coverage.out | tail -n 1

test-html: test-cover ## Generate HTML coverage report and open in browser
	@echo "📊 Generating HTML coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"
	@open coverage.html 2>/dev/null || xdg-open coverage.html 2>/dev/null || echo "Open coverage.html manually"

test-short: ## Run tests in short mode (skip integration tests)
	@echo "🧪 Running tests (short mode)..."
	go test ./test/... -short -count=1

# ──────────────────────────────────────────────
#  Code Quality
# ──────────────────────────────────────────────
lint: ## Run golangci-lint (requires golangci-lint)
	@which golangci-lint > /dev/null || (echo "❌ golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run ./...

fmt: ## Format all Go source files
	@echo "📝 Formatting code..."
	go fmt ./...
	@echo "✅ Done"

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	go vet ./...
	@echo "✅ No issues"

tidy: ## Tidy Go modules
	@echo "🧹 Tidying modules..."
	go mod tidy
	@echo "✅ Done"

# ──────────────────────────────────────────────
#  Cleanup
# ──────────────────────────────────────────────
clean: ## Clean build artifacts, coverage files, and logs
	@echo "🧹 Cleaning..."
	@rm -f $(APP_BINARY)
	@rm -f coverage.out coverage.html
	@rm -rf logs/*.log
	@echo "✅ Clean"

clean-all: clean docker-down-volumes ## Clean everything including Docker volumes

# ──────────────────────────────────────────────
#  Setup
# ──────────────────────────────────────────────
setup: ## Install Go dependencies
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed"

env: ## Create .env from .env.example if missing
	@if [ ! -f .env ]; then \
		echo "📝 Creating .env from .env.example..."; \
		cp .env.example .env; \
		echo "✅ .env created — please edit with your credentials"; \
	else \
		echo "ℹ️  .env already exists"; \
	fi
