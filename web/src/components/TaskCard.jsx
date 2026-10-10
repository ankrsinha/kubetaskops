import { Calendar, Trash2, ArrowRight, ArrowLeft } from 'lucide-react';

export function TaskCard({ task, onUpdateStatus, onDelete }) {
  const getPrevStatus = (current) => {
    if (current === 'Completed') return 'InProgress';
    if (current === 'InProgress') return 'Pending';
    return null;
  };

  const getNextStatus = (current) => {
    if (current === 'Pending') return 'InProgress';
    if (current === 'InProgress') return 'Completed';
    return null;
  };

  const prevStatus = getPrevStatus(task.status);
  const nextStatus = getNextStatus(task.status);

  return (
    <div className="task-card" id={`task-card-${task.id}`}>
      <div className="task-card-header">
        <span className={`priority-badge priority-${task.priority}`}>
          {task.priority}
        </span>
        <button
          className="btn-danger-subtle"
          title="Delete Task"
          onClick={() => onDelete(task.id)}
        >
          <Trash2 size={14} />
        </button>
      </div>

      <h4 className="task-title">{task.title}</h4>
      {task.description && <p className="task-description">{task.description}</p>}

      <div className="task-footer">
        <div className="task-date">
          <Calendar size={13} />
          <span>
            {task.due_date 
              ? new Date(task.due_date).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
              : 'No due date'}
          </span>
        </div>

        <div className="task-actions">
          {prevStatus && (
            <button
              className="btn-move"
              title={`Move to ${prevStatus}`}
              onClick={() => onUpdateStatus(task.id, prevStatus)}
            >
              <ArrowLeft size={12} />
            </button>
          )}
          {nextStatus && (
            <button
              className="btn-move"
              title={`Move to ${nextStatus}`}
              onClick={() => onUpdateStatus(task.id, nextStatus)}
            >
              <ArrowRight size={12} />
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
