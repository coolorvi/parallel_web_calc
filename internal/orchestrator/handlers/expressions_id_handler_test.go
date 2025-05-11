package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetExpression(t *testing.T) {
	orch := NewOrchestrator(nil)

	exprID := "test-id-123"
	expr := &Expression{
		ID:     exprID,
		Status: "completed",
		Result: floatPtr(42),
	}
	orch.mu.Lock()
	orch.expressions[exprID] = expr
	orch.mu.Unlock()

	req := httptest.NewRequest("GET", "/expressions/"+exprID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": exprID})
	w := httptest.NewRecorder()
	orch.GetExpression(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]*Expression
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["expression"].ID != exprID {
		t.Errorf("expected expression ID %q, got %q", exprID, body["expression"].ID)
	}

	req = httptest.NewRequest("GET", "/expressions/bad-id", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "bad-id"})
	w = httptest.NewRecorder()
	orch.GetExpression(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404 for bad id, got %d", resp.StatusCode)
	}
}

func floatPtr(f float64) *float64 {
	return &f
}
