package handlers

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	ID string `json:"id"`
}

type Task struct {
	ID            string  `json:"id"`
	ExpressionID  string  `json:"expression_id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type Expression struct {
	ID          string             `json:"id"`
	Status      string             `json:"status"`
	Result      *float64           `json:"result,omitempty"`
	Tasks       []string           `json:"tasks"`
	TaskResults map[string]float64 `json:"task_results"`
}

var (
	Tasks       = make(map[string]*Task)
	Expressions = make(map[string]*Expression)
	mutex       sync.Mutex
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var req CalcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
		return
	}

	expr := strings.TrimSpace(req.Expression)
	if expr == "" {
		http.Error(w, "Expression not valid", http.StatusUnprocessableEntity)
		return
	}

	node, err := parser.ParseExpr(expr)
	if err != nil {
		http.Error(w, "Fail to parse", http.StatusInternalServerError)
		return
	}

	expressionID := uuid.New().String()

	expression := &Expression{
		ID:          expressionID,
		Status:      "in_progress",
		Tasks:       []string{},
		TaskResults: make(map[string]float64),
	}

	mutex.Lock()
	Expressions[expressionID] = expression
	mutex.Unlock()

	switch v := node.(type) {
	case *ast.BinaryExpr:
		arg1 := extractValue(v.X)
		arg2 := extractValue(v.Y)
		taskId := uuid.New().String()

		task := &Task{
			ID:           taskId,
			ExpressionID: expressionID,
			Arg1:         arg1,
			Arg2:         arg2,
			Operation:    v.Op.String(),
		}

		mutex.Lock()
		Tasks[taskId] = task

		Expressions[expressionID].Tasks = append(Expressions[expressionID].Tasks, taskId)
		mutex.Unlock()
	}

	response := CalcResponse{ID: expressionID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func extractValue(n ast.Expr) float64 {
	if lit, ok := n.(*ast.BasicLit); ok {
		val, _ := strconv.ParseFloat(lit.Value, 64)
		return val
	}
	return 0
}
