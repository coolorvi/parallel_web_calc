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
	t.Run("Valid request", func(t *testing.T) {
		requestBody := `{"expression": "2 + 2 * 2"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		handlers.CalculateHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", res.StatusCode)
		}

		var resp handlers.CalcResponse
		if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if resp.ID == "" {
			t.Errorf("Expected a valid UUID, got an empty string")
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(`{invalid json`))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		handlers.CalculateHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Expected status 422, got %d", res.StatusCode)
		}
	})

	t.Run("Empty expression", func(t *testing.T) {
		requestBody := `{"expression": ""}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		handlers.CalculateHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Expected status 422, got %d", res.StatusCode)
		}
	})

	t.Run("Invalid method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", nil)

		rec := httptest.NewRecorder()
		handlers.CalculateHandler(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Expected status 422, got %d", res.StatusCode)
		}
	})
}
