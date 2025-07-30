# Text Similarity API - Go + Python Hybrid
A text similarity API that combines Go's excellent concurrency and HTTP handling with Python's powerful machine learning capabilities.

## Architecture
- **Go Server**: Fast HTTP server with JSON API, request validation, logging, and error handling
- **Python Service**: ML-powered sentence similarity computation using SentenceTransformers
- **Communication**: Go calls Python via subprocess with JSON stdin/stdout communication
- **Model**: Uses `sentence-transformers/all-MiniLM-L6-v2` for semantic similarity

## Features
- ⚡ **High Performance**: Go handles HTTP requests and concurrency efficiently
- 🧠 **ML-Powered**: Python provides state-of-the-art sentence similarity
- 🐳 **Containerized**: Multi-stage Docker build optimizes image size
- 📝 **Well-Documented**: Comprehensive API documentation and examples
- 🔍 **Observable**: Health checks, request logging, and error handling
- 🛡️ **Production-Ready**: Input validation, timeouts, and security features

## Quick Start

### Using Docker (Recommended)

```bash
# Build and run with Docker Compose
docker-compose up --build

# Test the API
curl -X POST http://localhost:8080/api/v1/similarity \
  -H "Content-Type: application/json" \
  -d '{"sentence1": "AI is transforming the world", "sentence2": "Artificial intelligence is changing society"}'
```

### Local Development

```bash
# Set up development environment
make dev-setup

# Run in development mode
make dev

# Test the service
make test
```

## API Endpoints

### POST /api/v1/similarity

Calculate semantic similarity between two sentences.

**Request:**
```json
{
  "sentence1": "AI is transforming the world.",
  "sentence2": "Artificial intelligence is changing society."
}
```

**Response:**
```json
{
  "sentence1": "AI is transforming the world.",
  "sentence2": "Artificial intelligence is changing society.",
  "similarity": 0.7892,
  "processed_at": "2025-07-30T10:30:45Z"
}
```

### GET /health

Health check endpoint for monitoring.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-07-30T10:30:45Z",
  "service": "text-similarity-api"
}
```

### GET /docs

API documentation endpoint.

### GET /

Service information and available endpoints.

## Project Structure

```
.
├── main.go                          # Go HTTP server
├── go.mod                           # Go dependencies
├── app/
│   ├── similarity_service.py        # Python ML service
│   └── requirements.txt             # Python dependencies
├── Dockerfile                       # Multi-stage container build
├── docker-compose.yml              # Container orchestration
├── .dockerignore         
├── Makefile                        # Development commands
├── nginx.conf                      # Nginx server configuration
└── README.md                       
```

## Configuration

Environment variables:

- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin framework mode (`debug`, `release`)

## Development

### Prerequisites

- Go 1.21+
- Python 3.10+
- Docker & Docker Compose (optional)

### Available Commands

```bash
make dev-setup    # Set up development environment
make dev          # Run in development mode
make build        # Build Go application
make test         # Run tests
make docker-build # Build Docker image
make clean        # Clean up artifacts
make help         # Show all commands
```

### Testing

```bash
# Test Python service directly
make test-python

# Test API endpoint
make test-api

# Run all tests
make test
```

## Performance Characteristics
- **Concurrency**: Go handles multiple requests efficiently
- **Memory**: ~1-2GB RAM usage (mainly for ML model)
- **Latency**: Typically 50-200ms per similarity calculation
- **Throughput**: Scales with available CPU cores

## Security Features
- Input validation and sanitization
- Request timeouts (30s default)
- CORS headers configured
- Non-root container user
- No external network dependencies at runtime

## Deployment

### Docker Production

```bash
# Production build
docker-compose --profile production up -d

# With custom configuration
PORT=3000 GIN_MODE=release docker-compose up
```

## Monitoring

- Health endpoint: `GET /health`
- Request logging with timing
- Error tracking and recovery
- Container health checks

## Model Information

- **Model**: `sentence-transformers/all-MiniLM-L6-v2`
- **Size**: ~90MB
- **Performance**: Good balance of speed and accuracy
- **Language**: Optimized for English
- **Output**: Cosine similarity score (0.0 to 1.0)

## Contributing
1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Run `make format lint test`
5. Submit a pull request

### Development Tips

**Hot reloading:**
```bash
# Use air for Go hot reloading
go install github.com/cosmtrek/air@latest
air
```

**Debug Python service:**
```bash
# Test Python service directly
echo '{"sentence1": "test", "sentence2": "example"}' | python3 python_service/similarity_service.py
```

**Profile performance:**
```bash
# Go profiling
go tool pprof http://localhost:8080/debug/pprof/profile
```

## API Examples

### cURL Examples

```bash
# Basic similarity check
curl -X POST http://localhost:8080/api/v1/similarity \
  -H "Content-Type: application/json" \
  -d '{
    "sentence1": "The cat sat on the mat",
    "sentence2": "A feline rested on the rug"
  }'

# Health check
curl http://localhost:8080/health

# API documentation
curl http://localhost:8080/docs
```

### Python Client Example

```python
import requests
import json

def calculate_similarity(sentence1, sentence2, api_url="http://localhost:8080"):
    response = requests.post(
        f"{api_url}/api/v1/similarity",
        json={
            "sentence1": sentence1,
            "sentence2": sentence2
        }
    )
    return response.json()

# Usage
result = calculate_similarity(
    "Machine learning is advancing rapidly",
    "AI technology is progressing quickly"
)
print(f"Similarity: {result['similarity']:.4f}")
```


## Framework Used
- [SentenceTransformers](https://github.com/UKPLab/sentence-transformers)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Hugging Face Transformers](https://github.com/huggingface/transformers)