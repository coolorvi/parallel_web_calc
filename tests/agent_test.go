package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/coolorvi/parallel_web_calc/internal/agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestWorker_SimpleCalculation(t *testing.T) {
	task := agent.Task{
		ID:           "1",
		ExpressionID: "exp_1",
		Arg1:         5,
		Arg2:         3,
		Operation:    "+",
	}

	var wg sync.WaitGroup
	jobs := make(chan agent.Task)
	results := make(chan agent.Result, 1)

	wg.Add(1)
	go agent.Worker(jobs, results, &wg)

	jobs <- task
	close(jobs)

	wg.Wait()

	result := <-results
	assert.Equal(t, result.ID, "1")
	assert.Equal(t, result.Result, 8.0)
}

func TestSendResult_Success(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodPost)

		var result agent.Result
		err := json.NewDecoder(r.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, result.ID, "1")
		assert.Equal(t, result.Result, 8.0)

		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	agent.RESULT_URL = server.URL + "/internal/task"

	result := agent.Result{
		ID:           "1",
		ExpressionID: "exp_1",
		Result:       8.0,
	}

	err := agent.SendResult(result)
	assert.NoError(t, err)
}

func TestSendResult_Failure(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(agent.RESULT_URL, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	agent.RESULT_URL = server.URL

	result := agent.Result{
		ID:           "1",
		ExpressionID: "exp_1",
		Result:       8.0,
	}

	err := agent.SendResult(result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "received non-OK status")
}
