# Stage 1: Go
FROM golang:1.21-alpine AS go-builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o /main .


# Stage 2: Python
FROM python:3.11-slim AS python-base

WORKDIR /app

ENV PIP_NO_CACHE_DIR=1 \
    HF_HOME=/app/cache

COPY app/requirements.txt .

RUN pip install -r requirements.txt

RUN python -c "from sentence_transformers import SentenceTransformer; SentenceTransformer('sentence-transformers/all-MiniLM-L6-v2')"

# Stage 3
FROM python:3.11-slim

WORKDIR /app

ENV HF_HOME=/app/cache

RUN apt-get update && \
    apt-get install -y --no-install-recommends curl && \
    rm -rf /var/lib/apt/lists/* && \
    useradd --system --create-home appuser

COPY --from=go-builder /main /app/main
COPY --from=python-base /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=python-base /app/cache /app/cache
COPY app/ /app/

RUN chown -R appuser:appuser /app && chmod +x /app/main

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./main"]