package traditional

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// -------------------------------------------------------------
// Project Repository Methods
// -------------------------------------------------------------

func (r *Repository) ListProjects(ctx context.Context) ([]Project, error) {
	query := `
		SELECT id, name, description, lead, target_date, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Lead, &p.TargetDate, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *Repository) GetProjectByID(ctx context.Context, id uuid.UUID) (*Project, error) {
	query := `
		SELECT id, name, description, lead, target_date, created_at, updated_at
		FROM projects
		WHERE id = $1
	`
	var p Project
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Lead, &p.TargetDate, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	return &p, nil
}

func (r *Repository) CreateProject(ctx context.Context, req CreateProjectRequest) (*Project, error) {
	query := `
		INSERT INTO projects (name, description, lead, target_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, lead, target_date, created_at, updated_at
	`
	var p Project
	err := r.pool.QueryRow(ctx, query, req.Name, req.Description, req.Lead, req.TargetDate.TimePtr()).Scan(
		&p.ID, &p.Name, &p.Description, &p.Lead, &p.TargetDate, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert project: %w", err)
	}
	return &p, nil
}

func (r *Repository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// -------------------------------------------------------------
// Task Repository Methods
// -------------------------------------------------------------

func (r *Repository) ListTasks(ctx context.Context, projectID *uuid.UUID, status string) ([]Task, error) {
	query := `
		SELECT id, project_id, title, description, priority, status, due_date, created_at, updated_at
		FROM tasks
		WHERE ($1::uuid IS NULL OR project_id = $1)
		  AND ($2 = '' OR status = $2)
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, projectID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Priority, &t.Status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan task row: %w", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *Repository) GetTaskByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	query := `
		SELECT id, project_id, title, description, priority, status, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`
	var t Task
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Priority, &t.Status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return &t, nil
}

func (r *Repository) CreateTask(ctx context.Context, req CreateTaskRequest) (*Task, error) {
	query := `
		INSERT INTO tasks (project_id, title, description, priority, status, due_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, project_id, title, description, priority, status, due_date, created_at, updated_at
	`
	var t Task
	err := r.pool.QueryRow(ctx, query, req.ProjectID, req.Title, req.Description, req.Priority, req.Status, req.DueDate.TimePtr()).Scan(
		&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Priority, &t.Status, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert task: %w", err)
	}
	return &t, nil
}

func (r *Repository) UpdateTask(ctx context.Context, id uuid.UUID, req UpdateTaskRequest) (*Task, error) {
	// First fetch existing task
	existing, err := r.GetTaskByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.ProjectID != nil {
		existing.ProjectID = *req.ProjectID
	}
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Priority != nil {
		existing.Priority = *req.Priority
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.DueDate != nil {
		existing.DueDate = req.DueDate.TimePtr()
	}

	query := `
		UPDATE tasks
		SET project_id = $2, title = $3, description = $4, priority = $5, status = $6, due_date = $7
		WHERE id = $1
		RETURNING id, project_id, title, description, priority, status, due_date, created_at, updated_at
	`
	var updated Task
	err = r.pool.QueryRow(ctx, query, id, existing.ProjectID, existing.Title, existing.Description, existing.Priority, existing.Status, existing.DueDate).Scan(
		&updated.ID, &updated.ProjectID, &updated.Title, &updated.Description, &updated.Priority, &updated.Status, &updated.DueDate, &updated.CreatedAt, &updated.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	return &updated, nil
}

func (r *Repository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
