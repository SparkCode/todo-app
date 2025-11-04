import { useState, type ChangeEvent, type FormEvent } from 'react';
import './TaskForm.css';

interface Task {
  id: number;
  task_header: string;
  task_description: string;
  completed: boolean;
}

interface TaskFormProps {
  onTaskCreated: (tasks: Task[]) => void;
}

export function TaskForm({ onTaskCreated }: TaskFormProps) {
  const [formData, setFormData] = useState({
    task_header: '',
    task_description: '',
  });
  const [formError, setFormError] = useState<string | null>(null);
  const [statusMessage, setStatusMessage] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleInputChange = (event: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = event.target;
    setFormData(prev => ({
      ...prev,
      [name]: value,
    }));
  };

  const fetchTasks = async () => {
    try {
      const response = await fetch('/api/get-tasks');
      if (!response.ok) {
        throw new Error('Failed to refresh tasks');
      }
      const updatedTasks = (await response.json()) as Task[];
      onTaskCreated(updatedTasks);
    } catch (error) {
      setFormError(error instanceof Error ? error.message : 'Unable to refresh tasks');
    }
  };

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    if (!formData.task_header.trim()) {
      setFormError('Task title is required.');
      return;
    }

    setIsSubmitting(true);
    setFormError(null);
    setStatusMessage(null);

    try {
      const response = await fetch('/api/create-task', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          task_header: formData.task_header.trim(),
          task_description: formData.task_description.trim(),
          completed: false,
        }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Failed to create task');
      }

      setStatusMessage('Task created successfully.');
      setFormData({ task_header: '', task_description: '' });
      await fetchTasks();
    } catch (error) {
      setFormError(error instanceof Error ? error.message : 'Failed to create task');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <>
      <form className="task-form" onSubmit={handleSubmit}>
        <div className="form-field">
          <label className="form-label" htmlFor="task_header">
            Task Title
          </label>
          <input
            id="task_header"
            name="task_header"
            type="text"
            value={formData.task_header}
            onChange={handleInputChange}
            disabled={isSubmitting}
            className="form-input"
            placeholder="Enter task title"
            required
          />
        </div>
        <div className="form-field">
          <label className="form-label" htmlFor="task_description">
            Description
          </label>
          <textarea
            id="task_description"
            name="task_description"
            value={formData.task_description}
            onChange={handleInputChange}
            disabled={isSubmitting}
            className="form-textarea"
            placeholder="Optional task description"
            rows={3}
          />
        </div>
        <button type="submit" className="submit-button" disabled={isSubmitting}>
          {isSubmitting ? 'Creating...' : 'Create Task'}
        </button>
      </form>
      {formError && <div className="submit-error">{formError}</div>}
      {statusMessage && <div className="submit-success">{statusMessage}</div>}
    </>
  );
}
