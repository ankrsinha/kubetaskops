package traditional

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// FlexibleDate handles both YYYY-MM-DD and RFC3339 date JSON formats
type FlexibleDate time.Time

func (d *FlexibleDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" || s == "null" {
		return nil
	}
	// Try YYYY-MM-DD format (standard HTML5 date input)
	if t, err := time.Parse("2006-01-02", s); err == nil {
		*d = FlexibleDate(t)
		return nil
	}
	// Try RFC3339 format
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		*d = FlexibleDate(t)
		return nil
	}
	// Try ISO timestamp without timezone
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		*d = FlexibleDate(t)
		return nil
	}
	return fmt.Errorf("invalid date format '%s' (expected YYYY-MM-DD)", s)
}

func (d *FlexibleDate) TimePtr() *time.Time {
	if d == nil {
		return nil
	}
	t := time.Time(*d)
	return &t
}

// Project represents the traditional project domain model
type Project struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Lead        string     `json:"lead"`
	TargetDate  *time.Time `json:"target_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Task represents the traditional task domain model
type Task struct {
	ID          uuid.UUID  `json:"id"`
	ProjectID   uuid.UUID  `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    string     `json:"priority"` // Low, Medium, High, Critical
	Status      string     `json:"status"`   // Pending, InProgress, Completed
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// DTOs for client requests
type CreateProjectRequest struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Lead        string        `json:"lead"`
	TargetDate  *FlexibleDate `json:"target_date,omitempty"`
}

type UpdateProjectRequest struct {
	Name        *string       `json:"name,omitempty"`
	Description *string       `json:"description,omitempty"`
	Lead        *string       `json:"lead,omitempty"`
	TargetDate  *FlexibleDate `json:"target_date,omitempty"`
}

type CreateTaskRequest struct {
	ProjectID   uuid.UUID     `json:"project_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Priority    string        `json:"priority"`
	Status      string        `json:"status"`
	DueDate     *FlexibleDate `json:"due_date,omitempty"`
}

type UpdateTaskRequest struct {
	ProjectID   *uuid.UUID    `json:"project_id,omitempty"`
	Title       *string       `json:"title,omitempty"`
	Description *string       `json:"description,omitempty"`
	Priority    *string       `json:"priority,omitempty"`
	Status      *string       `json:"status,omitempty"`
	DueDate     *FlexibleDate `json:"due_date,omitempty"`
}
