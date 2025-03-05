package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/handlers"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TestExpressionHandler_Success(t *testing.T) {
	exprID := uuid.New().String()

	handlers.Expressions = map[string]*handlers.Expression{
		exprID: {
			ID:     exprID,
			Status: "completed",
			Result: nil,
		},
	}

	req, _ := http.NewRequest("GET", "/api/v1/expressions/"+exprID, nil)
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/expressions/{id}", handlers.ExpressionHandler)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался 200 OK, но получен %d", rr.Code)
	}

	var response map[string]map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Ошибка при разборе JSON: %v", err)
	}

	if response["expression"]["id"] != exprID {
		t.Errorf("Ожидался ID '%s', но получен '%s'", exprID, response["expression"]["id"])
	}

	if response["expression"]["status"] != "completed" {
		t.Errorf("Ожидался статус 'completed', но получен '%s'", response["expression"]["status"])
	}
}

func TestExpressionHandler_NotFound(t *testing.T) {
	nonExistentID := uuid.New().String()

	req, _ := http.NewRequest("GET", "/api/v1/expressions/"+nonExistentID, nil)
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/expressions/{id}", handlers.ExpressionHandler)
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("Ожидался 404 Not Found, но получен %d", rr.Code)
	}
}

func TestExpression_NotFound(t *testing.T) {
	exprID := uuid.New().String()
	handlers.Expressions = map[string]*handlers.Expression{
		exprID: nil,
	}

	req, _ := http.NewRequest("GET", "/api/v1/expressions/"+exprID, nil)
	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/expressions/{id}", handlers.ExpressionHandler)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Ожидался 404 Not Found, но получен %d", rr.Code)
	}
}
