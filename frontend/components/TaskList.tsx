import './TaskList.css';

interface Task {
  id: number;
  task_header: string;
  task_description: string;
  completed: boolean;
}

interface TaskListProps {
  tasks: Task[];
  onDeleteTask: (taskId: number) => void;
  onToggleTask: (taskId: number, completed: boolean) => void;
}

export function TaskList({ tasks, onDeleteTask, onToggleTask }: TaskListProps) {
  if (tasks.length === 0) {
    return <div className="empty-message">No tasks yet. Create your first task!</div>;
  }

  return (
    <div className="task-list-container">
      <div className="task-list">
        {tasks.map(task => (
          <div
            key={task.id}
            className={`task-card ${task.completed ? 'completed' : ''}`}
          >
            <div className="task-header-row">
              <input
                type="checkbox"
                checked={task.completed}
                onChange={() => onToggleTask(task.id, !task.completed)}
                className="task-checkbox"
              />
              <h3 className={`task-title ${task.completed ? 'completed' : ''}`}>
                {task.task_header}
              </h3>
              <button
                className="delete-button"
                onClick={() => onDeleteTask(task.id)}
                aria-label="Delete task"
              >
                ✕
              </button>
            </div>
            {task.task_description && (
              <p className="task-description">
                {task.task_description}
              </p>
            )}
            <div className="task-id">
              ID: {task.id}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

