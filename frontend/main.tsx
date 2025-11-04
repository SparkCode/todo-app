import React from 'react';
import ReactDOM from 'react-dom/client';
import { App } from './components/App';

// Get initial data from server-side rendering
declare global {
  interface Window {
    __INITIAL_DATA__?: {
      tasks: any[];
      error: string | null;
    };
  }
}

const initialData = window.__INITIAL_DATA__ || { tasks: [], error: null };

// Use hydrateRoot for SSR support
ReactDOM.hydrateRoot(
  document.getElementById('app')!,
  <React.StrictMode>
    <App initialTasks={initialData.tasks} initialError={initialData.error} />
  </React.StrictMode>
);

