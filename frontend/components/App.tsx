import { useState } from 'react';
import './App.css';
import { TaskForm } from './TaskForm';
import { TaskList } from './TaskList';

interface Task {
  id: number;
  task_header: string;
  task_description: string;
  completed: boolean;
}

interface AppProps {
  initialTasks?: Task[];
  initialError?: string | null;
}

export function App({ initialTasks = [], initialError = null }: AppProps) {
  const [tasks, setTasks] = useState<Task[]>(initialTasks);

  const handleDeleteTask = async (taskId: number) => {
    try {
      const response = await fetch(`/api/delete-task/${taskId}`, {
        method: 'DELETE',
      });

      if (response.ok) {
        setTasks(tasks.filter(task => task.id !== taskId));
      } else {
        console.error('Failed to delete task');
      }
    } catch (error) {
      console.error('Error deleting task:', error);
    }
  };

  return (
    <div className="app-container">
      <header className="app-header">
        <h1 className="app-title">
          Todo App
        </h1>
      </header>
      {initialError && <div className="error-message">Error: {initialError}</div>}
      {!initialError && (
        <>
          <TaskForm onTaskCreated={setTasks} />
          <TaskList tasks={tasks} onDeleteTask={handleDeleteTask} />
        </>
      )}
    </div>
  );
}


