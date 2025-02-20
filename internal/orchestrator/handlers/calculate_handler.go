package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	ID string `json:"id"`
}

type Expression struct {
	ID      string   `json:"id"`
	RawExpr string   `json:"raw_expr"`
	Status  string   `json:"status"`
	Result  *float64 `json:"result,omitempty"`
}

var (
	expressions = make(map[string]*Expression)
	mutex       sync.Mutex
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	var req CalcRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	if req.Expression == "" {
		http.Error(w, "Expression is empty", http.StatusUnprocessableEntity)
		return
	}

	exprID := uuid.New().String()

	mutex.Lock()
	expressions[exprID] = &Expression{
		ID:      exprID,
		RawExpr: req.Expression,
		Status:  "pending",
	}
	mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CalcResponse{ID: exprID})
}
