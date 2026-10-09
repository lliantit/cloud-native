package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
}

var (
	movies    = make(map[int]Movie)
	currentID = 1
	mu        sync.Mutex
)

func main() {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// Create та Read All
	mux.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			mu.Lock()
			list := make([]Movie, 0, len(movies))
			for _, m := range movies {
				list = append(list, m)
			}
			mu.Unlock()
			json.NewEncoder(w).Encode(list)
			return
		}

		if r.Method == http.MethodPost {
			var m Movie
			if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			mu.Lock()
			m.ID = currentID
			movies[currentID] = m
			currentID++
			mu.Unlock()
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(m)
			return
		}
	})

	// Read One, Update, Delete
	mux.HandleFunc("/movies/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Отримуємо ID з URL
		idStr := strings.TrimPrefix(r.URL.Path, "/movies/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "Некоректний ID"}`, http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		if r.Method == http.MethodGet {
			if m, exists := movies[id]; exists {
				json.NewEncoder(w).Encode(m)
			} else {
				http.Error(w, `{"error": "Фільм не знайдено"}`, http.StatusNotFound)
			}
			return
		}

		if r.Method == http.MethodPut {
			if _, exists := movies[id]; !exists {
				http.Error(w, `{"error": "Фільм не знайдено"}`, http.StatusNotFound)
				return
			}
			var updated Movie
			if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			updated.ID = id
			movies[id] = updated
			json.NewEncoder(w).Encode(updated)
			return
		}

		if r.Method == http.MethodDelete {
			if _, exists := movies[id]; !exists {
				http.Error(w, `{"error": "Фільм не знайдено"}`, http.StatusNotFound)
				return
			}
			delete(movies, id)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	})

	println("Go сервер запущено на http://localhost:8081")
	http.ListenAndServe(":8081", mux)
}
