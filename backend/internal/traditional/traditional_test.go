package traditional

import (
	"context"
	"testing"
	"time"

	"github.com/ankrsinha/kubetaskops/backend/pkg/store"
	"github.com/google/uuid"
)

func TestTraditionalCRUD_Integration(t *testing.T) {
	ctx := context.Background()
	cfg := store.DefaultConfig()

	pool, err := store.NewPostgresPool(ctx, cfg)
	if err != nil {
		t.Skipf("Skipping integration test: PostgreSQL not reachable: %v", err)
	}
	defer pool.Close()

	repo := NewRepository(pool)
	service := NewService(repo)

	// 1. Create a Project
	targetDate := FlexibleDate(time.Now().Add(7 * 24 * time.Hour))
	project, err := service.CreateProject(ctx, CreateProjectRequest{
		Name:        "Test Integration Project",
		Description: "Testing traditional CRUD flow",
		Lead:        "Jane Doe",
		TargetDate:  &targetDate,
	})
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}
	if project.ID == uuid.Nil {
		t.Fatal("expected non-nil project ID")
	}

	// 2. Create a Task under the project
	dueDate := FlexibleDate(time.Now().Add(3 * 24 * time.Hour))
	task, err := service.CreateTask(ctx, CreateTaskRequest{
		ProjectID:   project.ID,
		Title:       "Test Integration Task",
		Description: "Verify task insertion and foreign key reference",
		Priority:    "High",
		Status:      "Pending",
		DueDate:     &dueDate,
	})
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	if task.ID == uuid.Nil || task.ProjectID != project.ID {
		t.Fatalf("task data mismatch: %+v", task)
	}

	// 3. Update Task Status
	updatedStatus := "Completed"
	updatedTask, err := service.UpdateTask(ctx, task.ID, UpdateTaskRequest{
		Status: &updatedStatus,
	})
	if err != nil {
		t.Fatalf("failed to update task: %v", err)
	}
	if updatedTask.Status != "Completed" {
		t.Fatalf("expected status 'Completed', got '%s'", updatedTask.Status)
	}

	// 4. List Tasks for project
	tasks, err := service.ListTasks(ctx, &project.ID, "")
	if err != nil {
		t.Fatalf("failed to list tasks: %v", err)
	}
	if len(tasks) == 0 {
		t.Fatal("expected at least 1 task")
	}

	// 5. Cleanup: Delete Project (should cascade delete task)
	if err := service.DeleteProject(ctx, project.ID); err != nil {
		t.Fatalf("failed to delete project: %v", err)
	}

	// Verify task was cascade deleted
	_, err = service.GetTask(ctx, task.ID)
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound for cascade deleted task, got: %v", err)
	}
}
