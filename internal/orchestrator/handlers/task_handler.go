package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
)

func NewTaskHandler(tasks map[string]*Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		for _, task := range tasks {
			json.NewEncoder(w).Encode(task)
			return
		}

		http.Error(w, "no tasks available", http.StatusNotFound)
	}
}

func NewTaskResultHandler(tasks map[string]*Task, results map[string]float64, mutex *sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var result Result
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		if _, exists := tasks[result.ID]; !exists {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}

		results[result.ID] = result.Result
		w.WriteHeader(http.StatusOK)
	}
}
