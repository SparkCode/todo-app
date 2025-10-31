package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
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
