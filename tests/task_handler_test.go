package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/handlers"
)

func TestExpressions(t *testing.T) {
	tests := []struct {
		name           string
		setupData      func()
		expectedStatus int
		expectedBody   ExpressionsResponse
	}{
		{
			name: "Successfully get expressions",
			setupData: func() {
				zero := 0.0
				handlers.Expressions = map[string]*handlers.Expression{
					"1": {ID: "1", Status: "completed", Result: nil},
					"2": {ID: "2", Status: "in progress", Result: nil},
					"3": {ID: "3", Status: "completed", Result: &zero},
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody: ExpressionsResponse{
				Expressions: []ExpressionResponse{
					{ID: "1", Status: "completed"},
					{ID: "2", Status: "in progress"},
					{ID: "3", Status: "completed", Result: new(float64)},
				},
			},
		},
		{
			name: "Empty expressions",
			setupData: func() {
				handlers.Expressions = make(map[string]*handlers.Expression)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   ExpressionsResponse{Expressions: []ExpressionResponse{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupData()

			req, err := http.NewRequest(http.MethodGet, "/expressions", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.ExpressionsHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			var response ExpressionsResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}

			if len(response.Expressions) != len(tt.expectedBody.Expressions) {
				t.Errorf("unexpected number of expressions: got %d, want %d",
					len(response.Expressions), len(tt.expectedBody.Expressions))
			}
		})
	}
}
