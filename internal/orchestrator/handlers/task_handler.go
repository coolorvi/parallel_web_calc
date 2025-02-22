package handlers

import (
	"encoding/json"
	"net/http"
)

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type Result struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

var (
	tasks   = make(map[string]*Task)
	results = make(map[string]float64)
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		mutex.Lock()
		defer mutex.Unlock()
		for _, task := range tasks {
			json.NewEncoder(w).Encode(task)
			return
		}
		http.Error(w, "no tasks available", http.StatusNotFound)
		return
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func TaskResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var res Result
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, "invalid data", http.StatusUnprocessableEntity)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()
	if _, exists := tasks[res.ID]; !exists {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	results[res.ID] = res.Result
	w.WriteHeader(http.StatusOK)
}
