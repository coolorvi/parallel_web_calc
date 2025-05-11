package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/coolorvi/parallel_web_calc/internal/agent"
	"github.com/stretchr/testify/assert"
)

func TestCompute_Addition(t *testing.T) {
	task := agent.Task{
		Arg1:          5,
		Arg2:          3,
		Operation:     "+",
		OperationTime: 0,
	}
	result := agent.Compute(task)
	assert.Equal(t, 8.0, result, "Expected result is 8")
}

func TestCompute_Subtraction(t *testing.T) {
	task := agent.Task{
		Arg1:          5,
		Arg2:          3,
		Operation:     "-",
		OperationTime: 0,
	}
	result := agent.Compute(task)
	assert.Equal(t, 2.0, result, "Expected result is 2")
}

func TestCompute_Multiplication(t *testing.T) {
	task := agent.Task{
		Arg1:          5,
		Arg2:          3,
		Operation:     "*",
		OperationTime: 0,
	}
	result := agent.Compute(task)
	assert.Equal(t, 15.0, result, "Expected result is 15")
}

func TestCompute_Division(t *testing.T) {
	task := agent.Task{
		Arg1:          6,
		Arg2:          2,
		Operation:     "/",
		OperationTime: 0,
	}
	result := agent.Compute(task)
	assert.Equal(t, 3.0, result, "Expected result is 3")
}

func TestCompute_DivideByZero(t *testing.T) {
	task := agent.Task{
		Arg1:          1,
		Arg2:          0,
		Operation:     "/",
		OperationTime: 0,
	}
	result := agent.Compute(task)
	assert.Equal(t, 0.0, result, "Expected result is 0 when dividing by zero")
}

func TestAgent_Worker(t *testing.T) {
	task := agent.Task{
		ID:        1,
		Arg1:      3,
		Arg2:      2,
		Operation: "+",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/internal/task" {
			resp := struct {
				Task agent.Task `json:"task"`
			}{
				Task: task,
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		if r.Method == "POST" && r.URL.Path == "/internal/task" {
			var requestData map[string]interface{}
			json.NewDecoder(r.Body).Decode(&requestData)
			assert.Equal(t, 1.0, requestData["id"])
			assert.Equal(t, 5.0, requestData["result"])
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))

	defer server.Close()

	client := &http.Client{}
	resp, err := client.Get(server.URL + "/internal/task")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var taskResponse struct {
		Task agent.Task `json:"task"`
	}
	err = json.NewDecoder(resp.Body).Decode(&taskResponse)
	assert.Nil(t, err)

	result := agent.Compute(taskResponse.Task)

	resultData := map[string]interface{}{
		"id":     taskResponse.Task.ID,
		"result": result,
	}
	reqBody, err := json.Marshal(resultData)
	assert.Nil(t, err)

	resp, err = client.Post(server.URL+"/internal/task", "application/json", bytes.NewBuffer(reqBody))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
