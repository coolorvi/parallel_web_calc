package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/handlers"
	"github.com/stretchr/testify/assert"
)

func TestCalculateHandler_ValidExpressions(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		statusCode int
	}{
		{"Simple addition", "3+5", http.StatusCreated},
		{"Simple subtraction", "10-7", http.StatusCreated},
		{"Multiplication", "4*6", http.StatusCreated},
		{"Division", "8/2", http.StatusCreated},
		{"Complex expression", "(3+5)*2-4/2", http.StatusCreated},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(handlers.CalcRequest{Expression: tt.expression})
			req, err := http.NewRequest("POST", "/calculate", bytes.NewBuffer(requestBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			http.HandlerFunc(handlers.CalculateHandler).ServeHTTP(recorder, req)

			assert.Equal(t, tt.statusCode, recorder.Code)
			if recorder.Code == http.StatusCreated {
				var resp handlers.CalcResponse
				err = json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.ID)
			}
		})
	}
}

func TestCalculateHandler_InvalidExpressions(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		statusCode int
	}{
		{"Empty expression", "", http.StatusUnprocessableEntity},
		{"Invalid characters", "3&5", http.StatusInternalServerError},
		{"Unbalanced parentheses", "(3+5*2", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, _ := json.Marshal(handlers.CalcRequest{Expression: tt.expression})
			req, err := http.NewRequest("POST", "/calculate", bytes.NewBuffer(requestBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			http.HandlerFunc(handlers.CalculateHandler).ServeHTTP(recorder, req)

			assert.Equal(t, tt.statusCode, recorder.Code)
		})
	}
}
