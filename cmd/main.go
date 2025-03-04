package main

import (
	"github.com/coolorvi/parallel_web_calc/internal/agent"
	"github.com/coolorvi/parallel_web_calc/internal/orchestrator"
)

func main() {
	orchestrator.Start()
	agent.StartWorker()
}
