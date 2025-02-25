package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	orchestrator "github.com/coolorvi/parallel_web_calc/internal/parser"
	"github.com/google/uuid"
)

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	ID string `json:"id"`
}

type Result struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type Expression struct {
	ID     string   `json:"id"`
	Status string   `json:"status"`
	Result *float64 `json:"result,omitempty"`
	Tasks  []string `json:"tasks"`
}

var (
	Tasks       = make(map[string]*Task)
	Expressions = make(map[string]*Expression)
	mutex       sync.Mutex
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var req CalcRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	parsed, err := orchestrator.ParseExpression([]byte(`{"expression": "` + req.Expression + `"}`))
	if err != nil {
		http.Error(w, "Failed to parse expression", http.StatusUnprocessableEntity)
		return
	}

	var output orchestrator.OutputJSON
	if err := json.Unmarshal(parsed, &output); err != nil {
		http.Error(w, "Failed to process parsed expression", http.StatusInternalServerError)
		return
	}

	exprID := uuid.New().String()
	var taskIDs []string

	mutex.Lock()
	for _, expr := range output.SubExpressions {
		taskID := uuid.New().String()

		Tasks[taskID] = &Task{
			ID:            taskID,
			Arg1:          atof(expr.LeftOperand),
			Arg2:          atof(expr.RightOperand),
			Operation:     expr.Operator,
			OperationTime: 100,
		}
		taskIDs = append(taskIDs, taskID)
	}

	Expressions[exprID] = &Expression{
		ID:     exprID,
		Status: "pending",
		Tasks:  taskIDs,
	}
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CalcResponse{ID: exprID})
}

func atof(s string) float64 {
	var f float64
	json.Unmarshal([]byte(s), &f)
	return f
}
