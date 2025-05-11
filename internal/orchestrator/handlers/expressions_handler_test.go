package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrchestrator_GetExpressions(t *testing.T) {
	orch := NewOrchestrator(nil)

	expr1 := &Expression{ID: "1", Status: "pending"}
	expr2 := &Expression{ID: "2", Status: "completed"}

	orch.mu.Lock()
	orch.expressions["1"] = expr1
	orch.expressions["2"] = expr2
	orch.mu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/expressions", nil)
	rr := httptest.NewRecorder()

	orch.GetExpressions(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d", rr.Code, http.StatusOK)
	}

	var resp struct {
		Expressions []*Expression `json:"expressions"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(resp.Expressions) != 2 {
		t.Errorf("expected 2 expressions, got %d", len(resp.Expressions))
	}

	found := map[string]bool{}
	for _, e := range resp.Expressions {
		found[e.ID] = true
	}

	if !found["1"] || !found["2"] {
		t.Errorf("missing expected expressions in response: got %+v", resp.Expressions)
	}
}
