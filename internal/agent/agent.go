package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
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

func getTask() (*Task, error) {
	resp, err := http.Get("http://localhost/internal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func compute(task *Task) float64 {
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)
	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			return 0
		}
		return task.Arg1 / task.Arg2
	default:
		return 0
	}
}

func sendResult(result Result) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost/internal/task", "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task, err := getTask()
		if err != nil {
			continue
		}

		result := compute(task)
		sendResult(Result{ID: task.ID, Result: result})
	}
}

func StartAgent() {
	computingPower, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if computingPower <= 0 {
		computingPower = 1
	}

	var wg sync.WaitGroup
	for i := 0; i < computingPower; i++ {
		wg.Add(1)
		go worker(&wg)
	}
	wg.Wait()
}
