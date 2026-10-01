const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1/traditional';

export const api = {
  // Projects
  async getProjects() {
    const res = await fetch(`${API_BASE}/projects`);
    if (!res.ok) throw new Error(`Failed to fetch projects: ${res.statusText}`);
    const body = await res.json();
    return body.data || [];
  },

  async createProject(payload) {
    const res = await fetch(`${API_BASE}/projects`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.error || `Failed to create project: ${res.statusText}`);
    }
    const body = await res.json();
    return body.data;
  },

  async deleteProject(id) {
    const res = await fetch(`${API_BASE}/projects/${id}`, { method: 'DELETE' });
    if (!res.ok) throw new Error(`Failed to delete project: ${res.statusText}`);
  },

  // Tasks
  async getTasks(projectId, status) {
    const params = new URLSearchParams();
    if (projectId) params.append('project_id', projectId);
    if (status) params.append('status', status);

    const url = `${API_BASE}/tasks?${params.toString()}`;
    const res = await fetch(url);
    if (!res.ok) throw new Error(`Failed to fetch tasks: ${res.statusText}`);
    const body = await res.json();
    return body.data || [];
  },

  async createTask(payload) {
    const res = await fetch(`${API_BASE}/tasks`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.error || `Failed to create task: ${res.statusText}`);
    }
    const body = await res.json();
    return body.data;
  },

  async updateTask(id, payload) {
    const res = await fetch(`${API_BASE}/tasks/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.error || `Failed to update task: ${res.statusText}`);
    }
    const body = await res.json();
    return body.data;
  },

  async deleteTask(id) {
    const res = await fetch(`${API_BASE}/tasks/${id}`, { method: 'DELETE' });
    if (!res.ok) throw new Error(`Failed to delete task: ${res.statusText}`);
  },
};
