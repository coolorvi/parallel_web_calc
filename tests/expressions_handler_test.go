package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/handlers"
)

func setupTestData() {
	handlers.Expressions = map[string]*handlers.Expression{
		"1": {ID: "1", Status: "done", Result: floatPtr(6)},
		"2": {ID: "2", Status: "pending", Result: nil},
	}
}

func floatPtr(v float64) *float64 {
	return &v
}

func TestExpressionsHandler(t *testing.T) {
	setupTestData()

	req, _ := http.NewRequest("GET", "/api/v1/expressions", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ExpressionsHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался 200 OK, но получен %d", rr.Code)
	}

	var response map[string][]handlers.Expression
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Ошибка при разборе JSON: %v", err)
	}

	if len(response["expressions"]) != len(handlers.Expressions) {
		t.Errorf("Ожидаемое количество выражений %d, но получено %d", len(handlers.Expressions), len(response["expressions"]))
	}
}

func TestExpressionsHandler_EmptyList(t *testing.T) {
	handlers.Expressions = map[string]*handlers.Expression{}

	req, _ := http.NewRequest("GET", "/api/v1/expressions", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ExpressionsHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался 200 OK, но получен %d", rr.Code)
	}

	var response map[string][]handlers.Expression
	json.Unmarshal(rr.Body.Bytes(), &response)

	if len(response["expressions"]) != 0 {
		t.Errorf("Ожидался пустой массив, но получено %v", response["expressions"])
	}
}

func TestExpressionsHandler_WrongMethod(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v1/expressions", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ExpressionsHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался 500, но получен %d", rr.Code)
	}
}
