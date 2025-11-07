package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("sqlite", "./todos.db")
	if err != nil {
		return err
	}

	// Create tasks table
	createTableSQL := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task_header TEXT NOT NULL,
		task_description TEXT,
		completed BOOLEAN NOT NULL DEFAULT FALSE
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}

func main() {
	// Initialize database
	if err := initDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()
	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://localhost:3000" + r.URL.Path)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// Copy headers from upstream response
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(body)
	})

	http.HandleFunc("/api/create-task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var task struct {
				TaskHeader      string `json:"task_header"`
				TaskDescription string `json:"task_description"`
				Completed       bool   `json:"completed"`
			}
			err = json.Unmarshal(body, &task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if task.TaskHeader == "" {
				http.Error(w, "task_header is required", http.StatusBadRequest)
				return
			}
			_, err = db.Exec("INSERT INTO tasks (task_header, task_description, completed) VALUES (?, ?, ?)", task.TaskHeader, task.TaskDescription, task.Completed)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("Task created successfully"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	http.HandleFunc("/api/delete-task/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "DELETE":
			parts := strings.Split(r.URL.Path, "/")
			if len(parts) < 4 || parts[3] == "" {
				http.Error(w, "Task ID is required", http.StatusBadRequest)
				return
			}

			idStr := parts[3]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid task ID", http.StatusBadRequest)
				return
			}

			result, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if rowsAffected == 0 {
				http.Error(w, "Task not found", http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Task deleted successfully"))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/toggle-task/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PATCH":
			parts := strings.Split(r.URL.Path, "/")
			if len(parts) < 4 || parts[3] == "" {
				http.Error(w, "Task ID is required", http.StatusBadRequest)
				return
			}

			idStr := parts[3]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Invalid task ID", http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var requestBody struct {
				Completed bool `json:"completed"`
			}
			err = json.Unmarshal(body, &requestBody)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			result, err := db.Exec("UPDATE tasks SET completed = ? WHERE id = ?", requestBody.Completed, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if rowsAffected == 0 {
				http.Error(w, "Task not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "Task updated successfully"})
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/get-tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			tasksRows, err := db.Query("SELECT * FROM tasks")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			type Task struct {
				ID              int    `json:"id"`
				TaskHeader      string `json:"task_header"`
				TaskDescription string `json:"task_description"`
				Completed       bool   `json:"completed"`
			}

			var tasks []Task

			for tasksRows.Next() {
				var task Task
				err := tasksRows.Scan(&task.ID, &task.TaskHeader, &task.TaskDescription, &task.Completed)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				tasks = append(tasks, task)
			}

			if err = tasksRows.Err(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tasksRows.Close()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Fetch tasks from database
		tasksRows, err := db.Query("SELECT * FROM tasks")

		type Task struct {
			ID              int    `json:"id"`
			TaskHeader      string `json:"task_header"`
			TaskDescription string `json:"task_description"`
			Completed       bool   `json:"completed"`
		}

		var tasks []Task
		var errorMsg string

		if err != nil {
			errorMsg = err.Error()
		} else {
			for tasksRows.Next() {
				var task Task
				err := tasksRows.Scan(&task.ID, &task.TaskHeader, &task.TaskDescription, &task.Completed)
				if err != nil {
					errorMsg = err.Error()
					break
				}
				tasks = append(tasks, task)
			}

			if err = tasksRows.Err(); err != nil {
				errorMsg = err.Error()
			}

			tasksRows.Close()
		}

		// Prepare data to send to Node.js server
		requestData := map[string]interface{}{
			"tasks": tasks,
			"error": nil,
		}
		if errorMsg != "" {
			requestData["error"] = errorMsg
		}

		jsonData, err := json.Marshal(requestData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// POST tasks to Node.js server for SSR
		resp, err := http.Post("http://localhost:3000/", "application/json", bytes.NewReader(jsonData))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Copy headers from upstream response
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(body)
	})

	fmt.Println("server started on the port :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
