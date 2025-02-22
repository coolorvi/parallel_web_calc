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

var taskQueue = []Task{
	{ID: "1", Arg1: 10, Arg2: 5, Operation: "+", OperationTime: 100},
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if len(taskQueue) == 0 {
			http.Error(w, "no tasks available", http.StatusNotFound)
			return
		}
		task := taskQueue[0]
		taskQueue = taskQueue[1:]
		json.NewEncoder(w).Encode(task)
	} else if r.Method == http.MethodPost {
		var result map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
