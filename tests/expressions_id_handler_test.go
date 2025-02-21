package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/handlers"
)

func TestExpressionHandler_Success(t *testing.T) {
	setupTestData()

	req, _ := http.NewRequest("GET", "/api/v1/expressions/1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ExpressionHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался 200 OK, но получен %d", rr.Code)
	}

	var response map[string]handlers.Expression
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Ошибка при разборе JSON: %v", err)
	}

	if response["expression"].ID != "1" {
		t.Errorf("Ожидалось ID '1', но получено '%s'", response["expression"].ID)
	}
}

func TestExpressionHandler_NotFound(t *testing.T) {
	setupTestData()

	req, _ := http.NewRequest("GET", "/api/v1/expressions/999", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ExpressionHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Ожидался 404 Not Found, но получен %d", rr.Code)
	}
}

func TestExpressionHandler_InternalServerError(t *testing.T) {
	handlers.Expressions["1"] = nil

	req, _ := http.NewRequest("GET", "/api/v1/expressions/1", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ExpressionHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался 500 Internal Server Error, но получен %d", rr.Code)
	}
}
