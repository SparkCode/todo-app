import React from 'react';
import ReactDOM from 'react-dom/client';
import { App } from './App';

// Use hydrateRoot for SSR support
ReactDOM.hydrateRoot(
  document.getElementById('app')!,
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

