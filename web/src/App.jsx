import { useState, useEffect, useCallback } from 'react';
import { api } from './services/api.js';
import { Navbar } from './components/Navbar.jsx';
import { TaskBoard } from './components/TaskBoard.jsx';
import { CreateTaskModal } from './components/CreateTaskModal.jsx';
import { CreateProjectModal } from './components/CreateProjectModal.jsx';
import { AlertCircle, RefreshCw, CheckCircle2, Clock, ListTodo, User, Calendar } from 'lucide-react';

export function App() {
  // ── Theme ──────────────────────────────────────────────────────────────────
  const [theme, setTheme] = useState(() => {
    // Persist choice across sessions
    return localStorage.getItem('kubetasks-theme') || 'dark';
  });

  useEffect(() => {
    document.documentElement.setAttribute('data-theme', theme);
    localStorage.setItem('kubetasks-theme', theme);
  }, [theme]);

  const handleToggleTheme = useCallback(() => {
    setTheme((prev) => (prev === 'dark' ? 'light' : 'dark'));
  }, []);

  // ── Data ───────────────────────────────────────────────────────────────────
  const [projects, setProjects] = useState([]);
  const [selectedProjectId, setSelectedProjectId] = useState('');
  const [tasks, setTasks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Modals
  const [isTaskModalOpen, setIsTaskModalOpen] = useState(false);
  const [isProjectModalOpen, setIsProjectModalOpen] = useState(false);

  // Fetch projects
  const fetchProjects = useCallback(async () => {
    try {
      const data = await api.getProjects();
      setProjects(data);
    } catch (err) {
      setError(err.message || 'Failed to connect to backend');
    }
  }, []);

  // Fetch tasks
  const fetchTasks = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await api.getTasks(selectedProjectId || undefined);
      setTasks(data);
    } catch (err) {
      setError(err.message || 'Failed to connect to backend');
    } finally {
      setLoading(false);
    }
  }, [selectedProjectId]);

  useEffect(() => {
    fetchProjects();
  }, [fetchProjects]);

  useEffect(() => {
    fetchTasks();
  }, [fetchTasks]);

  // ── Handlers ───────────────────────────────────────────────────────────────

  const handleUpdateStatus = async (taskId, newStatus) => {
    try {
      // Optimistic update
      setTasks((prev) =>
        prev.map((t) => (t.id === taskId ? { ...t, status: newStatus } : t))
      );
      await api.updateTask(taskId, { status: newStatus });
    } catch (err) {
      alert(`Status update failed: ${err.message}`);
      fetchTasks();
    }
  };

  const handleDeleteTask = async (taskId) => {
    if (!confirm('Are you sure you want to delete this task?')) return;
    try {
      setTasks((prev) => prev.filter((t) => t.id !== taskId));
      await api.deleteTask(taskId);
    } catch (err) {
      alert(`Delete failed: ${err.message}`);
      fetchTasks();
    }
  };

  const handleCreateProject = async (payload) => {
    const newProj = await api.createProject(payload);
    setProjects((prev) => [newProj, ...prev]);
    setSelectedProjectId(newProj.id);
  };

  const handleCreateTask = async (payload) => {
    const newTask = await api.createTask(payload);
    setTasks((prev) => [newTask, ...prev]);
  };

  // ── Derived state ──────────────────────────────────────────────────────────
  const selectedProject = projects.find((p) => p.id === selectedProjectId);
  const totalCount      = tasks.length;
  const inProgressCount = tasks.filter((t) => t.status === 'InProgress').length;
  const completedCount  = tasks.filter((t) => t.status === 'Completed').length;

  // ── Render ─────────────────────────────────────────────────────────────────
  return (
    <div className="app-container">
      <Navbar
        projects={projects}
        selectedProjectId={selectedProjectId}
        onSelectProject={setSelectedProjectId}
        onOpenCreateTask={() => setIsTaskModalOpen(true)}
        onOpenCreateProject={() => setIsProjectModalOpen(true)}
        theme={theme}
        onToggleTheme={handleToggleTheme}
      />

      {/* Project / Workspace Overview Header */}
      <section className="overview-header">
        <div className="overview-title-section">
          <h2 className="overview-title">
            {selectedProject ? selectedProject.name : 'All Projects'}
          </h2>
          {selectedProject?.description && (
            <p className="overview-description">{selectedProject.description}</p>
          )}
          {selectedProject && (
            <div className="overview-meta">
              {selectedProject.lead && (
                <span className="meta-item">
                  <User size={14} />
                  <span>Lead: {selectedProject.lead}</span>
                </span>
              )}
              {selectedProject.target_date && (
                <span className="meta-item">
                  <Calendar size={14} />
                  <span>Due: {new Date(selectedProject.target_date).toLocaleDateString()}</span>
                </span>
              )}
            </div>
          )}
        </div>

        <div className="metrics-group">
          <div className="metric-chip">
            <ListTodo size={15} className="text-muted" />
            <span className="metric-label">Total</span>
            <span className="metric-val">{totalCount}</span>
          </div>
          <div className="metric-chip">
            <Clock size={15} className="text-sky" />
            <span className="metric-label">In Progress</span>
            <span className="metric-val">{inProgressCount}</span>
          </div>
          <div className="metric-chip">
            <CheckCircle2 size={15} className="text-emerald" />
            <span className="metric-label">Completed</span>
            <span className="metric-val">{completedCount}</span>
          </div>
        </div>
      </section>

      {/* Error alert */}
      {error && (
        <div className="error-alert">
          <div style={{ display: 'flex', alignItems: 'center', gap: '0.75rem' }}>
            <AlertCircle size={20} />
            <span>{error}</span>
          </div>
          <button className="btn btn-secondary" onClick={fetchTasks}>
            <RefreshCw size={14} />
            <span>Retry</span>
          </button>
        </div>
      )}

      {/* Main Kanban Content */}
      <main>
        {loading && tasks.length === 0 ? (
          <div className="empty-state">
            <RefreshCw size={24} className="animate-spin" />
            <span>Loading tasks…</span>
          </div>
        ) : (
          <TaskBoard
            tasks={tasks}
            onUpdateStatus={handleUpdateStatus}
            onDelete={handleDeleteTask}
          />
        )}
      </main>

      {/* Modals */}
      <CreateTaskModal
        isOpen={isTaskModalOpen}
        projects={projects}
        defaultProjectId={selectedProjectId}
        onClose={() => setIsTaskModalOpen(false)}
        onSubmit={handleCreateTask}
      />

      <CreateProjectModal
        isOpen={isProjectModalOpen}
        onClose={() => setIsProjectModalOpen(false)}
        onSubmit={handleCreateProject}
      />
    </div>
  );
}

export default App;
