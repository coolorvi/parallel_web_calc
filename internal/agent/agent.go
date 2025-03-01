package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"

	"go/constant"
	"go/token"
)

const (
	TASK_URL        = "http://localhost/internal/task"
	RESULT_URL      = "http://localhost/internal/task"
	COMPUTING_POWER = 3
)

type Task struct {
	ID            string `json:"id"`
	Arg1          string `json:"arg1"`
	Arg2          string `json:"arg2"`
	Operation     string `json:"operation"`
	OperationTime int    `json:"operation_time"`
}

type Result struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

func fetchTask() (*Task, error) {
	resp, err := http.Get(TASK_URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var task struct {
		Task Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}

	return &task.Task, nil
}

func sendResult(result Result) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp, err := http.Post(RESULT_URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send result, status code: %d", resp.StatusCode)
	}

	return nil
}

func computeTask(task *Task) (float64, error) {
	var x, y constant.Value

	x = constant.MakeFromLiteral(task.Arg1, token.FLOAT, 0)
	y = constant.MakeFromLiteral(task.Arg2, token.FLOAT, 0)

	if x.Kind() == constant.Unknown || y.Kind() == constant.Unknown {
		return 0, fmt.Errorf("invalid operands")
	}

	var result constant.Value
	switch task.Operation {
	case "+":
		result = constant.BinaryOp(x, token.ADD, y)
	case "-":
		result = constant.BinaryOp(x, token.SUB, y)
	case "*":
		result = constant.BinaryOp(x, token.MUL, y)
	case "/":
		if constant.Sign(y) == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		result = constant.BinaryOp(x, token.QUO, y)
	default:
		return 0, fmt.Errorf("unsupported operation: %s", task.Operation)
	}

	f, _ := new(big.Float).SetString(result.ExactString())
	res, _ := f.Float64()
	return res, nil
}

func worker(id int, jobs <-chan *Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range jobs {
		fmt.Printf("Worker %d processing task %s\n", id, task.ID)
		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

		result, err := computeTask(task)
		if err != nil {
			log.Printf("Worker %d failed to compute task %s: %v", id, task.ID, err)
			continue
		}

		if err := sendResult(Result{ID: task.ID, Result: result}); err != nil {
			log.Printf("Worker %d failed to send result for task %s: %v", id, task.ID, err)
		}
	}
}

func StartWorker() {
	jobs := make(chan *Task, 100)
	var wg sync.WaitGroup

	for w := 1; w <= COMPUTING_POWER; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	for {
		task, err := fetchTask()
		if err != nil {
			log.Printf("Error fetching task: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		if task == nil {
			log.Println("No tasks available, retrying...")
			time.Sleep(2 * time.Second)
			continue
		}

		jobs <- task
	}
}
