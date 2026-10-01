package traditional

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrInvalidInput = errors.New("invalid input data")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListProjects(ctx context.Context) ([]Project, error) {
	return s.repo.ListProjects(ctx)
}

func (s *Service) GetProject(ctx context.Context, id uuid.UUID) (*Project, error) {
	return s.repo.GetProjectByID(ctx, id)
}

func (s *Service) CreateProject(ctx context.Context, req CreateProjectRequest) (*Project, error) {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, fmt.Errorf("%w: project name cannot be empty", ErrInvalidInput)
	}
	return s.repo.CreateProject(ctx, req)
}

func (s *Service) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteProject(ctx, id)
}

func (s *Service) ListTasks(ctx context.Context, projectID *uuid.UUID, status string) ([]Task, error) {
	return s.repo.ListTasks(ctx, projectID, status)
}

func (s *Service) GetTask(ctx context.Context, id uuid.UUID) (*Task, error) {
	return s.repo.GetTaskByID(ctx, id)
}

func (s *Service) CreateTask(ctx context.Context, req CreateTaskRequest) (*Task, error) {
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		return nil, fmt.Errorf("%w: task title cannot be empty", ErrInvalidInput)
	}
	if req.ProjectID == uuid.Nil {
		return nil, fmt.Errorf("%w: project_id is required", ErrInvalidInput)
	}

	// Verify project exists
	if _, err := s.repo.GetProjectByID(ctx, req.ProjectID); err != nil {
		return nil, fmt.Errorf("referenced project does not exist: %w", err)
	}

	// Validate Priority
	if req.Priority == "" {
		req.Priority = "Medium"
	}
	validPriorities := map[string]bool{"Low": true, "Medium": true, "High": true, "Critical": true}
	if !validPriorities[req.Priority] {
		return nil, fmt.Errorf("%w: invalid priority '%s'", ErrInvalidInput, req.Priority)
	}

	// Validate Status
	if req.Status == "" {
		req.Status = "Pending"
	}
	validStatuses := map[string]bool{"Pending": true, "InProgress": true, "Completed": true}
	if !validStatuses[req.Status] {
		return nil, fmt.Errorf("%w: invalid status '%s'", ErrInvalidInput, req.Status)
	}

	return s.repo.CreateTask(ctx, req)
}

func (s *Service) UpdateTask(ctx context.Context, id uuid.UUID, req UpdateTaskRequest) (*Task, error) {
	if req.Priority != nil {
		validPriorities := map[string]bool{"Low": true, "Medium": true, "High": true, "Critical": true}
		if !validPriorities[*req.Priority] {
			return nil, fmt.Errorf("%w: invalid priority '%s'", ErrInvalidInput, *req.Priority)
		}
	}
	if req.Status != nil {
		validStatuses := map[string]bool{"Pending": true, "InProgress": true, "Completed": true}
		if !validStatuses[*req.Status] {
			return nil, fmt.Errorf("%w: invalid status '%s'", ErrInvalidInput, *req.Status)
		}
	}
	if req.ProjectID != nil {
		if _, err := s.repo.GetProjectByID(ctx, *req.ProjectID); err != nil {
			return nil, fmt.Errorf("referenced project does not exist: %w", err)
		}
	}

	return s.repo.UpdateTask(ctx, id, req)
}

func (s *Service) DeleteTask(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteTask(ctx, id)
}
