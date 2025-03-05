package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/handlers"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedID     string
	}{
		{
			name:           "Valid expression",
			requestBody:    `{"expression": "3 + 4"}`,
			expectedStatus: http.StatusCreated,
			expectedID:     "valid_uuid",
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{expression: "3 + 4"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedID:     "",
		},
		{
			name:           "Empty expression",
			requestBody:    `{"expression": ""}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedID:     "",
		},
		{
			name:           "Fail to parse expression",
			requestBody:    `{"expression": "3 ?? 4"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedID:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/calculate", bytes.NewBufferString(tt.requestBody))
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.CalculateHandler)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, status)
			}

			if tt.expectedStatus == http.StatusCreated {
				var resp handlers.CalcResponse
				if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
					t.Fatalf("could not decode response body: %v", err)
				}

				if resp.ID == "" {
					t.Errorf("expected non-empty ID, got empty")
				}

			}
		})
	}
}
