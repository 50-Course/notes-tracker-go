package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/50-Course/notes-tracker/cmd/repository"
	"github.com/50-Course/notes-tracker/shared/models"
	api "github.com/50-Course/notes-tracker/shared/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type TaskServiceServer struct {
	api.UnimplementedTaskServiceServer
	repo *repository.TaskRepository
}

// Creates new instance of TaskServiceServer
func NewTaskServiceServer(repo *repository.TaskRepository) *TaskServiceServer {
	return &TaskServiceServer{repo: repo}
}

// Handles our CreateTask RPC call for creating tasks
func (s *TaskServiceServer) CreateTask(ctx context.Context, req *api.CreateTaskRequest) (*api.CreateTaskResponse, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("Title is required")
	}

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	err := s.repo.CreateTask(ctx, task)
	if err != nil {
		// just propagate that error up our handler
		return nil, fmt.Errorf("Error creating task: %v", err)
	}

	return &api.CreateTaskResponse{
		Task: &api.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			CreatedAt:   task.CreatedAt.String(),
		},
	}, nil
}

// Handles call to get a specific task
func (s *TaskServiceServer) GetTask(ctx context.Context, req *api.GetTaskRequest) (*api.GetTaskResponse, error) {
	task, err := s.repo.GetTask(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("Task not found: %w", err)
	}

	return &api.GetTaskResponse{
		Task: &api.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			CreatedAt:   task.CreatedAt.String(),
		},
	}, nil
}

// Fetches all tasks
func (s *TaskServiceServer) ListTasks(ctx context.Context, req *api.ListTasksRequest) (*api.ListTasksResponse, error) {
	tasks, err := s.repo.ListTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error fetching tasks: %v", err)
	}

	var grpcTasks []*api.Task
	for _, task := range tasks {
		grpcTasks = append(grpcTasks, &api.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			CreatedAt:   task.CreatedAt.String(),
		})
	}

	return &api.ListTasksResponse{Tasks: grpcTasks}, nil
}

// Handles our UpdateTask RPC call for updating tasks
func (s *TaskServiceServer) UpdateTask(ctx context.Context, req *api.UpdateTaskRequest) (*api.UpdateTaskResponse, error) {
	task, err := s.repo.GetTask(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("Task not found: %w", err)
	}

	task.Title = req.Title
	task.Description = req.Description

	err = s.repo.UpdateTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("Error updating task: %v", err)
	}

	return &api.UpdateTaskResponse{
		Task: &api.Task{
			Id:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			CreatedAt:   task.CreatedAt.String(),
		},
	}, nil
}

// Handles our RPC call for deleting tasks
func (s *TaskServiceServer) DeleteTask(ctx context.Context, req *api.DeleteTaskRequest) (*api.DeleteTaskResponse, error) {
	err := s.repo.DeleteTask(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("Error deleting task: %v", err)
	}

	return &api.DeleteTaskResponse{Success: true}, nil
}

// StartServer starts the gRPC server
func RunGRPCServer(repo *repository.TaskRepository, port string) {
	address := fmt.Sprintf(":%s", port)
	listen, err := net.Listen("tcp", address)

	log.Printf("[gRPC] Attempting to start gRPC server on port %s", port)
	if err != nil {
		log.Fatalf("[gRPC] Failed to start GRPC server on %s: %v", address, err)
	}

	server := grpc.NewServer()
	api.RegisterTaskServiceServer(server, &TaskServiceServer{repo: repo})

	reflection.Register(server)

	fmt.Printf("[gRPC] Starting gRPC server on %s\n", address)
	if err := server.Serve(listen); err != nil {
		fmt.Printf("[gRPC] Failed to start GRPC server: %v", err)
	}
	fmt.Printf("gRPC server started on port %s\n", port)
}
