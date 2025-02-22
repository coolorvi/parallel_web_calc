package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/handlers"
)

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type Result struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

var (
	tasks   = make(map[string]*Task)
	results = make(map[string]float64)
	mutex   = sync.Mutex{}
)

func TestTaskHandler(t *testing.T) {
	tasks["1"] = &Task{ID: "1", Arg1: 1, Arg2: 2, Operation: "+", OperationTime: 100}
	r := httptest.NewRequest(http.MethodGet, "/internal/task", nil)
	w := httptest.NewRecorder()

	handlers.TaskHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTaskHandler_NoTask(t *testing.T) {
	tasks = make(map[string]*Task)
	r := httptest.NewRequest(http.MethodGet, "/internal/task", nil)
	w := httptest.NewRecorder()

	handlers.TaskHandler(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestTaskResultHandler(t *testing.T) {
	tasks["1"] = &Task{ID: "1", Arg1: 1, Arg2: 2, Operation: "+", OperationTime: 100}
	body := strings.NewReader(`{"id":"1", "result":3}`)
	r := httptest.NewRequest(http.MethodPost, "/internal/task", body)
	w := httptest.NewRecorder()

	handlers.TaskResultHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTaskResultHandler_TaskNotFound(t *testing.T) {
	body := strings.NewReader(`{"id":"999", "result":3}`)
	r := httptest.NewRequest(http.MethodPost, "/internal/task", body)
	w := httptest.NewRecorder()

	handlers.TaskResultHandler(w, r)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
