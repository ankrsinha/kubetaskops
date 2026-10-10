import { TaskCard } from './TaskCard.jsx';
import { CircleDot, Clock, CheckCircle2 } from 'lucide-react';

const COLUMNS = [
  {
    status: 'Pending',
    title: 'Pending',
    icon: <CircleDot size={16} className="text-amber-400" />,
    color: '#f59e0b',
  },
  {
    status: 'InProgress',
    title: 'In Progress',
    icon: <Clock size={16} className="text-sky-400" />,
    color: '#38bdf8',
  },
  {
    status: 'Completed',
    title: 'Completed',
    icon: <CheckCircle2 size={16} className="text-emerald-400" />,
    color: '#10b981',
  },
];

export function TaskBoard({ tasks, onUpdateStatus, onDelete }) {
  return (
    <div className="kanban-grid">
      {COLUMNS.map((col) => {
        const columnTasks = tasks.filter((t) => t.status === col.status);

        return (
          <div key={col.status} className="kanban-column" id={`column-${col.status}`} data-status={col.status}>
            <div className="column-header">
              <div className="column-title-group">
                {col.icon}
                <span>{col.title}</span>
              </div>
              <span className="column-count">{columnTasks.length}</span>
            </div>

            <div className="task-list">
              {columnTasks.length === 0 ? (
                <div className="empty-state">
                  <span>No tasks in {col.title.toLowerCase()}</span>
                </div>
              ) : (
                columnTasks.map((task) => (
                  <TaskCard
                    key={task.id}
                    task={task}
                    onUpdateStatus={onUpdateStatus}
                    onDelete={onDelete}
                  />
                ))
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}
