package main

import (
	"log"

	"github.com/coolorvi/parallel_web_calc/internal/agent"
	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/start"
	"github.com/coolorvi/parallel_web_calc/internal/storage"
)

func main() {
	if err := storage.InitDB("data.db"); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	go func() {
		log.Println("Start orchestrator")
		start.StartServer()
	}()

	log.Println("Start agent")
	agent.StartAgent()
}
