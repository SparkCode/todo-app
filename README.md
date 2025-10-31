# Go Todo App

A full-stack todo application built with Go backend and TypeScript/React frontend.

## Features

- RESTful API built with Go
- Modern React frontend with TypeScript
- Vite for fast development and building

## Prerequisites

- Go 1.x or higher
- Node.js 18+ and npm

## Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd go-todo-app
```

2. Install frontend dependencies:
```bash
cd frontend
npm install
```

## Running the Application

1. Start the Go backend:
```bash
go run main.go
```

2. In a separate terminal, start the frontend development server:
```bash
cd frontend
npm run dev
```

## Building for Production

Build the frontend:
```bash
cd frontend
npm run build
```

## Project Structure

- `main.go` - Go backend server
- `frontend/` - React TypeScript frontend
  - `main.tsx` - Frontend entry point
  - `vite.config.ts` - Vite configuration

