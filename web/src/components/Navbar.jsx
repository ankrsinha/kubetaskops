import { Layers, Plus, FolderPlus } from 'lucide-react';
import { ThemeToggle } from './ThemeToggle.jsx';

export function Navbar({
  projects,
  selectedProjectId,
  onSelectProject,
  onOpenCreateTask,
  onOpenCreateProject,
  theme,
  onToggleTheme,
}) {
  return (
    <header className="navbar">
      <div className="brand-section">
        <div className="brand-logo">
          <Layers size={24} />
          <span>KubeTaskOps</span>
        </div>
      </div>

      <div className="nav-actions">
        <select
          id="project-filter-select"
          className="project-filter-select"
          value={selectedProjectId}
          onChange={(e) => onSelectProject(e.target.value)}
        >
          <option value="">All Projects</option>
          {projects.map((p) => (
            <option key={p.id} value={p.id}>
              📁 {p.name}
            </option>
          ))}
        </select>

        <button 
          id="btn-create-project"
          className="btn btn-secondary" 
          onClick={onOpenCreateProject}
        >
          <FolderPlus size={16} />
          <span>New Project</span>
        </button>

        <button 
          id="btn-create-task"
          className="btn btn-primary" 
          onClick={onOpenCreateTask}
        >
          <Plus size={16} />
          <span>New Task</span>
        </button>

        <ThemeToggle theme={theme} onToggle={onToggleTheme} />
      </div>
    </header>
  );
}
