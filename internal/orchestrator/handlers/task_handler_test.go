package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTask(t *testing.T) {
	orchestrator := &Orchestrator{
		tasks: make(chan Task, 1),
	}

	task := Task{
		ID:        1,
		Arg1:      5,
		Arg2:      3,
		Operation: "+",
	}

	orchestrator.tasks <- task

	req := httptest.NewRequest(http.MethodGet, "/task", nil)
	w := httptest.NewRecorder()

	orchestrator.GetTask(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]Task
	err := json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)

	assert.Equal(t, task, result["task"])
}

func TestGetTask_NoTasks(t *testing.T) {
	orchestrator := &Orchestrator{
		tasks: make(chan Task, 1),
	}

	req := httptest.NewRequest(http.MethodGet, "/task", nil)
	w := httptest.NewRecorder()

	orchestrator.GetTask(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestPostResult(t *testing.T) {
	orchestrator := &Orchestrator{
		results: make(map[int]float64),
	}

	requestBody := []byte(`{
		"id": 1,
		"result": 42
	}`)

	req := httptest.NewRequest(http.MethodPost, "/result", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()

	orchestrator.PostResult(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	orchestrator.mu.Lock()
	defer orchestrator.mu.Unlock()
	assert.Equal(t, 42.0, orchestrator.results[1])
}

func TestPostResult_InvalidData(t *testing.T) {
	orchestrator := &Orchestrator{
		results: make(map[int]float64),
	}

	requestBody := []byte(`{
		"id": 1
	}`)

	req := httptest.NewRequest(http.MethodPost, "/result", bytes.NewReader(requestBody))
	w := httptest.NewRecorder()

	orchestrator.PostResult(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}
