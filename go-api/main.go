package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
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

	mux.HandleFunc("/io", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "I/O operation completed"}`))
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

	// Важка операція: паралельне виконання в окремих goroutines
	mux.HandleFunc("/cpu-async", func(w http.ResponseWriter, r *http.Request) {
		iterations := 500_000_000
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			count1 := 0
			for i := 0; i < iterations; i++ {
				count1++
			}
		}()

		go func() {
			defer wg.Done()
			count2 := 0
			for i := 0; i < iterations; i++ {
				count2++
			}
		}()

		wg.Wait()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "CPU async completed"}`))
	})
	// Важка операція: послідовне виконання двох циклів в одному потоці
	mux.HandleFunc("/cpu-seq", func(w http.ResponseWriter, r *http.Request) {
		iterations := 500_000_000

		// Перший важкий цикл
		count1 := 0
		for i := 0; i < iterations; i++ {
			count1++
		}

		// Другий важкий цикл (послідовно після першого)
		count2 := 0
		for i := 0; i < iterations; i++ {
			count2++
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "CPU sequential completed"}`))
	})



	println("Go сервер запущено на http://localhost:8081")
	http.ListenAndServe(":8081", mux)
}
