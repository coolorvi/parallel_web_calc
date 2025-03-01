package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/handlers"
)

func TestTaskHandler(t *testing.T) {
	tasks := make(map[string]*handlers.Task)
	mutex := &sync.Mutex{}
	tasks["1"] = &handlers.Task{ID: "1", Arg1: 1, Arg2: 2, Operation: "+", OperationTime: 100}

	handler := handlers.NewTaskHandler(tasks, mutex)

	r := httptest.NewRequest(http.MethodGet, "/internal/task", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTaskHandler_NoTask(t *testing.T) {
	tasks := make(map[string]*handlers.Task)
	mutex := &sync.Mutex{}

	handler := handlers.NewTaskHandler(tasks, mutex)

	r := httptest.NewRequest(http.MethodGet, "/internal/task", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestTaskResultHandler(t *testing.T) {
	tasks := make(map[string]*handlers.Task)
	results := make(map[string]float64)
	mutex := sync.Mutex{}

	tasks["1"] = &handlers.Task{ID: "1", Arg1: 1, Arg2: 2, Operation: "+", OperationTime: 100}

	body := strings.NewReader(`{"id":"1", "result":3}`)
	r := httptest.NewRequest(http.MethodPost, "/internal/task", body)
	w := httptest.NewRecorder()

	handler := handlers.NewTaskResultHandler(tasks, results, &mutex)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTaskResultHandler_TaskNotFound(t *testing.T) {
	tasks := make(map[string]*handlers.Task)
	results := make(map[string]float64)
	mutex := sync.Mutex{}

	body := strings.NewReader(`{"id":"999", "result":3}`)
	r := httptest.NewRequest(http.MethodPost, "/internal/task", body)
	w := httptest.NewRecorder()

	handler := handlers.NewTaskResultHandler(tasks, results, &mutex)
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
