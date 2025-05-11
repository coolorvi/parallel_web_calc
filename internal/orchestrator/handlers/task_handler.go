package handlers

import (
	"encoding/json"
	"net/http"
)

func (o *Orchestrator) GetTask(w http.ResponseWriter, r *http.Request) {
	select {
	case task := <-o.tasks:
		json.NewEncoder(w).Encode(map[string]Task{"task": task})
	default:
		http.Error(w, "No tasks available", http.StatusNotFound)
	}
}

func (o *Orchestrator) PostResult(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data", http.StatusUnprocessableEntity)
		return
	}

	o.mu.Lock()
	o.results[req.ID] = req.Result
	o.mu.Unlock()

	w.WriteHeader(http.StatusOK)
}
