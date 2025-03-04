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

const (
	TASK_URL   = "http://localhost:8080/internal/task"
	RESULT_URL = "http://localhost:8080/internal/task"
)

type Task struct {
	ID            string      `json:"id"`
	ExpressionID  string      `json:"expression_id"`
	Arg1          interface{} `json:"arg1"`
	Arg2          interface{} `json:"arg2"`
	Operation     string      `json:"operation"`
	OperationTime int         `json:"operation_time"`
}

type Result struct {
	ID           string  `json:"id"`
	ExpressionID string  `json:"expression_id"`
	Result       float64 `json:"result"`
}

func worker(jobs <-chan Task, results chan<- Result, resultsCache map[string]float64, cacheMutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range jobs {
		var arg1, arg2 float64
		var arg1Ready, arg2Ready bool

		cacheMutex.Lock()
		log.Printf("Кэш перед выполнением задачи %s: %+v", task.ID, resultsCache)
		if id, ok := task.Arg1.(string); ok {
			arg1, arg1Ready = resultsCache[id]
		} else if num, ok := task.Arg1.(float64); ok {
			arg1, arg1Ready = num, true
		}

		if id, ok := task.Arg2.(string); ok {
			arg2, arg2Ready = resultsCache[id]
		} else if num, ok := task.Arg2.(float64); ok {
			arg2, arg2Ready = num, true
		}
		cacheMutex.Unlock()

		if !arg1Ready || !arg2Ready {
			log.Printf("Ожидание аргументов для задачи: %s (Arg1: %v, Arg2: %v)", task.ID, task.Arg1, task.Arg2)
			for {
				time.Sleep(100 * time.Millisecond)

				cacheMutex.Lock()
				if id, ok := task.Arg1.(string); ok {
					arg1, arg1Ready = resultsCache[id]
				}
				if id, ok := task.Arg2.(string); ok {
					arg2, arg2Ready = resultsCache[id]
				}
				cacheMutex.Unlock()

				if arg1Ready && arg2Ready {
					break
				}
			}
		}

		var result float64
		switch task.Operation {
		case "+":
			result = arg1 + arg2
		case "-":
			result = arg1 - arg2
		case "*":
			result = arg1 * arg2
		case "/":
			if arg2 == 0 {
				log.Printf("Error: division by zero (Task ID: %s)", task.ID)
				continue
			}
			result = arg1 / arg2
		default:
			log.Printf("Unsupported operation: %s (Task ID: %s)", task.Operation, task.ID)
			continue
		}

		cacheMutex.Lock()
		resultsCache[task.ID] = result
		cacheMutex.Unlock()

		results <- Result{
			ID:           task.ID,
			ExpressionID: task.ExpressionID,
			Result:       result,
		}
	}
}

func sendResult(result Result) error {
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

	resultsCache := make(map[string]float64)
	cacheMutex := &sync.Mutex{}

	go func() {
		for w := 1; w <= comp_pow; w++ {
			wg.Add(1)
			go worker(jobs, results, resultsCache, cacheMutex, &wg)
		}
	}()

	go func() {
		wg.Wait()
		close(jobs)
		close(results)
	}()

	go func() {
		for result := range results {
			err := sendResult(result)
			if err != nil {
				log.Printf("Failed to send result: %v", err)
			}
		}
	}()
}
