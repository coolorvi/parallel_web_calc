package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var TASK_URL = "http://localhost:8080/internal/task"
var RESULT_URL = "http://localhost:8080/internal/task"

type Task struct {
	ID            string  `json:"id"`
	ExpressionID  string  `json:"expression_id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type Result struct {
	ID           string  `json:"id"`
	ExpressionID string  `json:"expression_id"`
	Result       float64 `json:"result"`
}

func Worker(jobs <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	time_add, _ := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	time_sub, _ := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	time_mul, _ := strconv.Atoi(os.Getenv("TIME_MULTPIPLICATIONS_MS"))
	time_div, _ := strconv.Atoi(os.Getenv("TIME_DIVISION_MS"))
	defer wg.Done()
	for task := range jobs {
		if task.Arg1 == 0 || task.Arg2 == 0 {
			log.Printf("Error: one or both arguments are zero (Task ID: %s)", task.ID)
			continue
		}

		var result float64
		switch task.Operation {
		case "+":
			time.Sleep(time.Duration(time_add) * time.Millisecond)
			result = task.Arg1 + task.Arg2
		case "-":
			time.Sleep(time.Duration(time_sub) * time.Millisecond)
			result = task.Arg1 - task.Arg2
		case "*":
			time.Sleep(time.Duration(time_mul) * time.Millisecond)
			result = task.Arg1 * task.Arg2
		case "/":
			time.Sleep(time.Duration(time_div) * time.Millisecond)
			if task.Arg2 == 0 {
				log.Printf("Error: division by zero (Task ID: %s)", task.ID)
				continue
			}
			result = task.Arg1 / task.Arg2
		default:
			log.Printf("Unsupported operation: %s (Task ID: %s)", task.Operation, task.ID)
			continue
		}

		results <- Result{
			ID:           task.ID,
			ExpressionID: task.ExpressionID,
			Result:       result,
		}
	}
}

func SendResult(result Result) error {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Printf("Ошибка сериализации результата: %v", err)
		return err
	}

	log.Printf("Отправка результата: %s", string(resultJSON))
	resp, err := http.Post(RESULT_URL, "application/json", bytes.NewBuffer(resultJSON))
	if err != nil {
		log.Printf("Ошибка отправки результата: %v", err)
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Сервер вернул ошибку: %d - %s", resp.StatusCode, string(body))
		return fmt.Errorf("received non-OK status: %d", resp.StatusCode)
	}

	log.Println("Результат успешно отправлен!")
	return nil
}

func StartWorker() {
	log.Println("StartWorker запущен")
	url := "http://localhost:8080/internal/task"
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading global variables")
		return
	}

	var wg sync.WaitGroup
	jobs := make(chan Task)
	results := make(chan Result, 100)

	comp := os.Getenv("COMPUTING_POWER")
	comp_pow, err := strconv.Atoi(comp)
	if err != nil {
		log.Fatalf("Error converting COMPUTING_POWER to int: %v", err)
		return
	}

	go func() {
		for {
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Ошибка запроса задач: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Printf("Ошибка чтения ответа: %v", err)
				time.Sleep(10 * time.Second)
				continue
			}

			log.Printf("Ответ от сервера перед разбором JSON: %s", string(body))

			var data struct {
				Task Task `json:"task"`
			}
			if err := json.Unmarshal(body, &data); err != nil {
				log.Printf("Ошибка разбора JSON: %v. Ответ сервера: %s", err, string(body))
				time.Sleep(10 * time.Second)
				continue
			}

			task := data.Task

			if task.ID == "" {
				log.Println("Нет доступных задач.")
				time.Sleep(10 * time.Second)
				continue
			}

			log.Printf("Получена задача: %+v", task)
			jobs <- task
		}
	}()

	go func() {
		for w := 1; w <= comp_pow; w++ {
			wg.Add(1)
			go Worker(jobs, results, &wg)
		}
	}()

	go func() {
		wg.Wait()
		close(jobs)
		close(results)
	}()

	go func() {
		for result := range results {
			err := SendResult(result)
			if err != nil {
				log.Printf("Failed to send result: %v", err)
			}
		}
	}()
}
