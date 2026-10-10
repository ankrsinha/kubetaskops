package traditional

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes registers all Phase 1 traditional REST endpoints
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/traditional/projects", h.ListProjects)
	mux.HandleFunc("POST /api/v1/traditional/projects", h.CreateProject)
	mux.HandleFunc("GET /api/v1/traditional/projects/{id}", h.GetProject)
	mux.HandleFunc("DELETE /api/v1/traditional/projects/{id}", h.DeleteProject)

	mux.HandleFunc("GET /api/v1/traditional/tasks", h.ListTasks)
	mux.HandleFunc("POST /api/v1/traditional/tasks", h.CreateTask)
	mux.HandleFunc("GET /api/v1/traditional/tasks/{id}", h.GetTask)
	mux.HandleFunc("PUT /api/v1/traditional/tasks/{id}", h.UpdateTask)
	mux.HandleFunc("DELETE /api/v1/traditional/tasks/{id}", h.DeleteTask)
}

// Helper JSON responders
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// -------------------------------------------------------------
// Project Handlers
// -------------------------------------------------------------

func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.ListProjects(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if projects == nil {
		projects = []Project{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": projects})
}

func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid UUID format")
		return
	}

	p, err := h.service.GetProject(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": p})
}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body: "+err.Error())
		return
	}

	p, err := h.service.CreateProject(r.Context(), req)
	if errors.Is(err, ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"data": p})
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid UUID format")
		return
	}

	if err := h.service.DeleteProject(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "project not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// -------------------------------------------------------------
// Task Handlers
// -------------------------------------------------------------

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	var projectID *uuid.UUID
	projectIDStr := r.URL.Query().Get("project_id")
	if projectIDStr != "" {
		parsed, err := uuid.Parse(projectIDStr)
		if err == nil {
			projectID = &parsed
		}
	}

	status := strings.TrimSpace(r.URL.Query().Get("status"))

	tasks, err := h.service.ListTasks(r.Context(), projectID, status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tasks == nil {
		tasks = []Task{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": tasks})
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid UUID format")
		return
	}

	t, err := h.service.GetTask(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": t})
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body: "+err.Error())
		return
	}

	t, err := h.service.CreateTask(r.Context(), req)
	if errors.Is(err, ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"data": t})
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid UUID format")
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON request body: "+err.Error())
		return
	}

	t, err := h.service.UpdateTask(r.Context(), id, req)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	if errors.Is(err, ErrInvalidInput) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": t})
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid UUID format")
		return
	}

	if err := h.service.DeleteTask(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
