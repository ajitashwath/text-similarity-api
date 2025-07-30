.PHONY: build run test clean docker-build docker-run docker-stop dev-setup format lint

BINARY_NAME=text-similarity-api
DOCKER_IMAGE=text-similarity-api:latest
PYTHON_SERVICE_DIR=python_service

build:
	@echo "Building Go application..."
	go build -o $(BINARY_NAME) main.go

run: build
	@echo "Starting application..."
	./$(BINARY_NAME)

dev:
	@echo "Running in development mode..."
	GIN_MODE=debug go run main.go

dev-setup:
	@echo "Setting up development environment..."
	@echo "Installing Go dependencies..."
	go mod tidy
	@echo "Setting up Python virtual environment..."
	python3 -m venv venv
	@echo "Installing Python dependencies..."
	./venv/bin/pip install -r $(PYTHON_SERVICE_DIR)/requirements.txt
	@echo "Pre-downloading ML model..."
	./venv/bin/python3 -c "from sentence_transformers import SentenceTransformer; SentenceTransformer('sentence-transformers/all-MiniLM-L6-v2')"
	@echo "Development environment ready!"

test-python:
	@echo "Testing Python service..."
	echo '{"sentence1": "Hello world", "sentence2": "Hi there"}' | python3 $(PYTHON_SERVICE_DIR)/similarity_service.py

test-api:
	@echo "Testing API endpoint..."
	curl -X POST http://localhost:8080/api/v1/similarity \
		-H "Content-Type: application/json" \
		-d '{"sentence1": "AI is transforming the world", "sentence2": "Artificial intelligence is changing society"}'

format:
	@echo "Formatting Go code..."
	go fmt ./...
	goimports -w .

lint:
	@echo "Linting Go code..."
	golangci-lint run

docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	@echo "Starting services with Docker Compose..."
	docker-compose up --build

docker-run-bg:
	@echo "Starting services in background..."
	docker-compose up -d --build

docker-stop:
	@echo "Stopping Docker services..."
	docker-compose down

logs:
	docker-compose logs -f

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	docker-compose down --volumes --remove-orphans
	docker system prune -f

install-tools:
	@echo "Installing development tools..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

test: test-python test-api
	@echo "All tests completed!"

help:
	@echo "Available commands:"
	@echo "  build         - Build the Go application"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run in development mode"
	@echo "  dev-setup     - Set up development environment"
	@echo "  test-python   - Test Python service"
	@echo "  test-api      - Test API endpoint"
	@echo "  format        - Format Go code"
	@echo "  lint          - Lint Go code"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  docker-stop   - Stop Docker services"
	@echo "  clean         - Clean up build artifacts"
	@echo "  help          - Show this help message"