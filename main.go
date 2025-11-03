package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"

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
		task_description TEXT
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://localhost:3000")
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

	log.Fatal(http.ListenAndServe(":8080", nil))
}
