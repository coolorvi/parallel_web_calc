package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coolorvi/parallel_web_calc/internal/agent"
)

type MockTask struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

func TestGetTask(t *testing.T) {
	mockTask := MockTask{
		ID:            "1",
		Arg1:          2,
		Arg2:          3,
		Operation:     "+",
		OperationTime: 1000,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockTask)
	}))
	defer server.Close()

	task, err := agent.GetTask()
	if err != nil {
		t.Fatalf("Error fetching task: %v", err)
	}

	if task.ID != mockTask.ID || task.Operation != mockTask.Operation {
		t.Errorf("Incorrect task received: %+v", task)
	}
}

func TestCompute(t *testing.T) {
	task := agent.Task{
		ID:        "1",
		Arg1:      5,
		Arg2:      3,
		Operation: "-",
	}

	expectedResult := 2.0
	actualResult := agent.Compute(&task)

	if actualResult != expectedResult {
		t.Errorf("Expected result %f, got %f", expectedResult, actualResult)
	}
}

func TestSendResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		var result agent.Result
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			t.Fatalf("Error parsing JSON: %v", err)
		}
		if result.Result != 8.0 {
			t.Errorf("Expected result 8.0, got %v", result.Result)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := agent.SendResult(agent.Result{ID: "1", Result: 8.0})
	if err != nil {
		t.Fatalf("Error sending result: %v", err)
	}
}

func TestAgentLifecycle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			mockTask := MockTask{
				ID:        "1",
				Arg1:      4,
				Arg2:      2,
				Operation: "/",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockTask)
		} else if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	go agent.StartAgent()
	time.Sleep(2 * time.Second)
}
