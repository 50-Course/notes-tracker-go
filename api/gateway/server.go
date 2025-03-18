package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/uptrace/bunrouter"
	"google.golang.org/grpc"

	api "github.com/50-Course/notes-tracker/shared/proto"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Represents an API Gateway that will be used to route requests to the appropriate service
type Gateway struct {
	grpcClient api.TaskServiceClient
}

// creates a new Task API Gateway instance
// grpcAddress: the address of the gRPC server
// returns a new Gateway instance and an error if any
func NewGateway(grpcAddress string) (*Gateway, error) {
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := api.NewTaskServiceClient(conn)
	return &Gateway{grpcClient: client}, nil
}

/// --- API Gateway Handlers ---

// Handles the request to list all tasks
//
// It retrieves tasks from the gRPC service and returns them as a JSON response.
//
// Parameters:
//
//	w: The http.ResponseWriter to write the response to.
//	req: The bunrouter.Request containing the HTTP request.
//
// Returns:
//
//	error: An error if the operation fails, or nil if successful.
//	  - If the gRPC call fails, it returns a 500 Internal Server Error with a JSON payload
//	    containing an "error" and "error_message" field.
//	  - If successful, it returns a 200 OK with a JSON payload containing the list of tasks.
//
// Example:
//
//	Request: GET /tasks
//	Response (Success): 200 OK, JSON: {"tasks": [...]}
//	Response (Error):   500 Internal Server Error, JSON: {"error": "Failed to list tasks", "error_message": "grpc: ..."}
func (g *Gateway) ListTasksHandler(w http.ResponseWriter, req bunrouter.Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := g.grpcClient.ListTasks(ctx, &api.ListTasksRequest{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return bunrouter.JSON(w, bunrouter.H{
			"error":         "Failed to list tasks",
			"error_message": err.Error(),
		})
	}

	return bunrouter.JSON(w, resp)
}

// Handles the request to create a new task
//
// It creates a new task using the gRPC service and returns the created task as a JSON response.
//
// Parameters:
//
//	w: The http.ResponseWriter to write the response to.
//	req: The bunrouter.Request containing the HTTP request.
//
// Returns:
//
//	error: An error if the operation fails, or nil if successful.
//	  - If the gRPC call fails, it returns a 500 Internal Server Error with a JSON payload
//	    containing an "error" and "error_message" field.
//	  - If successful, it returns a 201 Created with a JSON payload containing the created task.
//
// Example:
//
//	Request: POST /tasks
//	Body: {"title": "Task 1", "description": "Description 1"}
//	Response (Success): 201 Created, JSON: {"task": {...}}
//	Response (Error):   500 Internal Server Error, JSON: {"error": "Failed to create task", "error_message": "grpc: ..."}
func (g *Gateway) CreateTaskHandler(w http.ResponseWriter, req bunrouter.Request) error {
	// Serializer for the request body
	var requestSerializer struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(req.Body).Decode(&requestSerializer); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return bunrouter.JSON(w, bunrouter.H{
			"error":         "Failed to create task",
			"error_message": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.grpcClient.CreateTask(ctx, &api.CreateTaskRequest{
		Title:       requestSerializer.Title,
		Description: requestSerializer.Description,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return bunrouter.JSON(w, bunrouter.H{
			"error":         "Failed to create task",
			"error_message": err.Error(),
		})
	}

	w.WriteHeader(http.StatusCreated)
	return bunrouter.JSON(w, resp)
}

// Handles the request to get a task by ID
//
// It retrieves a task by ID from the gRPC service and returns it as a JSON response.
//
// Parameters:
//
//	w: The http.ResponseWriter to write the response to.
//	req: The bunrouter.Request containing the HTTP request.
//
// Returns:
//
//	error: An error if the operation fails, or nil if successful.
//	  - If the gRPC call fails, it returns a 500 Internal Server Error with a JSON payload
//	    containing an "error" and "error_message" field.
//	  - If successful, it returns a 200 OK with a JSON payload containing the task.
//
// Example:
//
//	Request: GET /tasks/1
//	Response (Success): 200 OK, JSON: {"task": {...}}
//	Response (Error):   500 Internal Server Error, JSON: {"error": "Failed to get task", "error_message": "grpc: ..."}
func (g *Gateway) GetTaskHandler(w http.ResponseWriter, req bunrouter.Request) error {
	id := req.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.grpcClient.GetTask(ctx, &api.GetTaskRequest{Id: id})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return bunrouter.JSON(w, bunrouter.H{
			"error":         "Task not found",
			"error_message": err.Error(),
		})
	}

	return bunrouter.JSON(w, resp)
}

// Handles the request to update a task
//
// Tasks are updated with some new information using the gRPC service and returns the updated task as a JSON response.
//
// Parameters:
//
//	w: The http.ResponseWriter to write the response to.
//	req: The bunrouter.Request containing the HTTP request.
//
// Returns:
//
//	error: An error if the operation fails, or nil if successful.
//	  - If the gRPC call fails, it returns a 500 Internal Server Error with a JSON payload
//	    containing an "error" and "error_message" field.
//	  - If successful, it returns a 200 OK with a JSON payload containing the updated task.
func (g *Gateway) UpdateTaskHandler(w http.ResponseWriter, req bunrouter.Request) error {
	id := req.Param("id")

	var updateRequestSerializer struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(req.Body).Decode(&updateRequestSerializer); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return bunrouter.JSON(w, bunrouter.H{
			"error":         "Invalid request payload. Please review the request body and try again",
			"error_message": err.Error(),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := g.grpcClient.UpdateTask(ctx, &api.UpdateTaskRequest{
		Id:          id,
		Title:       updateRequestSerializer.Title,
		Description: updateRequestSerializer.Description,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return bunrouter.JSON(w, bunrouter.H{
			"error":         "Internal Server Error. Failed to update task",
			"error_message": err.Error(),
		})
	}

	// TODO: add a "message" field to the response
	return bunrouter.JSON(w, resp)
}

// Handles the request to delete a task
//
// It deletes a task by ID using the gRPC service.
//
// Parameters:
//
//	w: The http.ResponseWriter to write the response to.
//	req: The bunrouter.Request containing the HTTP request.
//
// Returns:
//
//	error: An error if the operation fails, or nil if successful.
//	  - If the gRPC call fails, it returns a 500 Internal Server Error with a JSON payload
//	    containing an "error" and "error_message" field.
//	  - If successful, it returns a 204 No Content response.
func (g *Gateway) DeleteTaskHandler(w http.ResponseWriter, req bunrouter.Request) error {
	taskID := req.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := g.grpcClient.DeleteTask(ctx, &api.DeleteTaskRequest{Id: taskID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return bunrouter.JSON(w, bunrouter.H{
			"error":         "Internal Server Error. Failed to delete task",
			"error_message": err.Error(),
		})
	}

	w.WriteHeader(http.StatusNoContent)
	return bunrouter.JSON(w, bunrouter.H{
		"message": "Task deleted successfully",
	})
}

// Initalizes a HTTP router with gRPC integration
//
// @title Notes Tracker API
// @version 1
// @description This is the API Gateway for the Notes Tracker, a simple task management application. handling HTTP requests and translating them to gRPC calls.
// @contact.name 50-Course
// @contact.url https://github.com/50-Course
// @license MIT
// @BasePath /api/v1
// @schemes http
func NewServer(gateway *Gateway) *bunrouter.Router {
	router := bunrouter.New()

	router.GET("/health", func(w http.ResponseWriter, req bunrouter.Request) error {
		return bunrouter.JSON(w, map[string]string{"message": "Hello, World!"})
	})

	router.WithGroup("/tasks", func(r *bunrouter.Group) {
		r.GET("", gateway.ListTasksHandler)
		r.POST("", gateway.CreateTaskHandler)
		r.GET("/:id", gateway.GetTaskHandler)
		r.PUT("/:id", gateway.UpdateTaskHandler)
		r.DELETE("/:id", gateway.DeleteTaskHandler)
	})

	// OpenAPI documentation
	// serve redoc by default
	router.GET("/api/v1/docs", func(w http.ResponseWriter, req bunrouter.Request) error {
		http.ServeFile(w, req.Request, "./docs/redoc.html")
		return nil
	})

	router.GET("/swagger/*", func(w http.ResponseWriter, req bunrouter.Request) error {
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		).ServeHTTP(w, req.Request)
		return nil
	})

	return router
}

func main() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	grpcAddress, addrExists := os.LookupEnv("INTERNAL_SERVER_ADDRESS")
	gatewayPort, portExists := os.LookupEnv("API_GATEWAY_PORT")

	if !addrExists {
		log.Fatal("INTERNAL_SERVER_ADDRESS not set in environment")
	}

	if !portExists {
		log.Printf("API_GATEWAY_PORT not set in environment. Defaulting to 8080")
		gatewayPort = "8080"
	}

	gateway, err := NewGateway(grpcAddress)
	if err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}

	r := NewServer(gateway)
	log.Printf("API Gateway is running on port %s", gatewayPort)
	log.Fatal(http.ListenAndServe(":"+gatewayPort, r))
}
