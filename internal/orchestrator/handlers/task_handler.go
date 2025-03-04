package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

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
		log.Printf("Отправляем задачу агенту: ID=%s, ExpressionID=%s, Arg1=%f, Arg2=%f, Operation=%s",
			task.ID, task.ExpressionID, task.Arg1, task.Arg2, task.Operation)

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

	log.Printf("Получен результат: %+v", result)

	mutex.Lock()

	expr, exists := Expressions[result.ExpressionID]
	if !exists {
		mutex.Unlock()
		log.Printf("Ошибка: Expression not found (ID: %s)", result.ExpressionID)
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	expr.TaskResults[result.ID] = result.Result
	log.Printf("Результат сохранён: ExpressionID=%s, TaskID=%s, Result=%f", result.ExpressionID, result.ID, result.Result)

	if len(expr.TaskResults) == len(expr.Tasks) {
		expr.Result = &result.Result
		expr.Status = "completed"
		log.Printf("Все задачи завершены. Статус выражения изменен на 'completed' (ExpressionID: %s)", result.ExpressionID)
	}

	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
}
