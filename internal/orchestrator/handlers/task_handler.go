package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Result struct {
	ID           string  `json:"id"`
	ExpressionID string  `json:"expressionid"`
	Result       float64 `json:"result"`
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetTask(w)
	case http.MethodPost:
		handlePostResult(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func handleGetTask(w http.ResponseWriter) {
	mutex.Lock()
	defer mutex.Unlock()

	for id, task := range Tasks {
		log.Printf("Sending task to agent: ID=%s, ExpressionID=%s", task.ID, task.ExpressionID)

		response := map[string]*Task{"task": task}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		delete(Tasks, id)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "No tasks available"})
}

func handlePostResult(w http.ResponseWriter, r *http.Request) {
	var result Result
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received result: %+v", result)

	mutex.Lock()
	defer mutex.Unlock()

	expr, exists := Expressions[result.ExpressionID]
	if !exists {
		log.Printf("Error: Expression not found (ID: %s)", result.ExpressionID)
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	expr.TaskResults[result.ID] = result.Result
	log.Printf("Result saved: ExpressionID=%s, TaskID=%s, Result=%f", result.ExpressionID, result.ID, result.Result)

	if len(expr.TaskResults) == len(expr.Tasks) {
		expr.Result = &result.Result
		expr.Status = "completed"
		log.Printf("All tasks completed. Expression status set to 'completed' (ExpressionID: %s)", result.ExpressionID)
	}

	w.WriteHeader(http.StatusOK)
}
