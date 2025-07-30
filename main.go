package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SentenceInput struct {
	Sentence1 string `json:"sentence1" binding:"required" validate:"min=1"`
	Sentence2 string `json:"sentence2" binding:"required" validate:"min=1"`
}

type SimilarityResponse struct {
	Sentence1  string  `json:"sentence1"`
	Sentence2  string  `json:"sentence2"`
	Similarity float64 `json:"similarity"`
	ProcessedAt string `json:"processed_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Message string `json:"message"`
}

type PythonRequest struct {
	Sentence1 string `json:"sentence1"`
	Sentence2 string `json:"sentence2"`
}

type PythonResponse struct {
	Similarity float64 `json:"similarity"`
	Error string `json:"error, omitempty"`
}

var validate *validator.validate

func init() {
	validate = validator.New()
}

func main() {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	r.use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H {
			"status": "healthy",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"service": "text-similarity-api",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H {
			"message": "Welcome to the Text Similarity API (Go + Python)",
			"version": "2.0.0",
			"endpoints": map[string]string {
				"similarity": "POST /api/v1/similarity",
				"health" : "GET /health",
				"docs" : "GET /docs",
			},
		})
	})

	r.GET("/docs", func(c *gin.Context) {
		docs := map[string]interface{} {
			"title" : "Text Similarity API",
			"description" : "An API to compute semantic similarity between sentences using Go & Python",
			"version" : "1.0.0",
			"endpoints": map[string]interface{} {
				"/api/v1/similarity": map[string]interface{} {
					"method": "POST",
					"description": "Calculate semantic similarity between two sentences",
					"request_body": map[string]interface{} {
						"sentence1": "string (required) - First sentence to compare",
						"sentence2": "string (required) - Second sentence to compare",
					},
					"response": map[string]interface{} {
						"sentence1": "string - Echo of first sentence",
						"sentence2": "string - Echo of second sentence",
						"similarity": "float - Similarity score (0.0 to 1.0)",
						"processed_at": "string - ISO timestamp of processing",
					},
					"example_request": map[string]string {
						"sentence1": "AI is transforming the world.",
						"sentence2": "Artificial intelligence is changing society.",
					},
				},
			},
		}
		c.JSON(http.StatusOK, docs)
	})

	v1 := r.Group("/api/v1") {
		v1.POST("/similarity", handleSimilarity)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Text Similarity API server on port %s", port)
	log.Printf("Endpoints available:")
	log.Printf("  GET  /           - API information")
	log.Printf("  GET  /health     - Health check")
	log.Printf("  GET  /docs       - API documentation")
	log.Printf("  POST /api/v1/similarity - Calculate similarity")

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func handleSimilarity(c *gin.Context) {
	var input SentenceInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse {
			Error: "validation_error",
			Message: "Invalid input format: " + err.Error(),
		})
		return
	}

	if err := validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse {
			Error: "validation_error",
			Message: "Validation failed: " + err.Error(),
		})
		return
	}

	input.Sentence1 = strings.TrimSpace(input.Sentence1)
	input.Sentence2 = strings.TrimSpace(input.Sentence2)

	if len(input.Sentence1) == 0 || len(input.Sentence2) == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse {
			Error: "empty_sentences",
			Message: "Both sentences must be non-empty",
		})
		return 
	}

	similarity, err := callPythonService(input)
	if err != nil {
		log.Printf("Error calling Python service: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse {
			Error: "internal_error",
			Message: "Failed to process similarity calculation",
		})
		return
	}

	response := SimilarityResponse {
		Sentence1: input.Sentence1,
		Sentence2: input.Sentence2,
		Similarity: similarity,
		ProcessedAt: time.Now().UTC().Format(time.RFC3339),
	}
	c.JSON(http.StatusOK, response)
}

func callPythonService(input SentenceInput) (float64, error) {
	pythonReq := PythonRequest {
		Sentence1: input.Sentence1,
		Sentence2: input.Sentence2,
	}

	reqData, err := json.Marshal(pythonReq)
	if err != nil {
		return 0, fmt.Errorf("Failed to Marshal request: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "python3", "app/similarity_service.py")
	cmd.Stdin = bytes.NewReader(reqData)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("python script failed: %w, stderr: %s", err, stderr.String())
	}

	var pythonResp PythonResponse
	if err := json.Unmarshal(stdout.Bytes(), &pythonResp); err != nil {
		return 0, fmt.Errorf("failed to parse python response: %w", err)
	}

	if pythonResp.Error != "" {
		return 0, fmt.Errorf("python service error: %s", pythonResp.Error)
	}
	
	return pythonResp.Similarity, nil
}
