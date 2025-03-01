package handlers

import (
	"encoding/json"
	"net/http"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for _, task := range Tasks {
		resp := map[string]*Task{"task": task}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		delete(Tasks, task.ID)
		return
	}

	http.Error(w, "No tasks available", http.StatusNotFound)
}
