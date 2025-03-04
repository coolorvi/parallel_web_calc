package orchestrator

import (
	"log"
	"net/http"

	"github.com/coolorvi/parallel_web_calc/internal/agent"
	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/handlers"
	"github.com/gorilla/mux"
)

func Start() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/calculate", handlers.CalculateHandler).Methods("POST")
	r.HandleFunc("/api/v1/expressions", handlers.ExpressionsHandler).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", handlers.ExpressionHandler).Methods("GET")
	r.HandleFunc("/internal/task", handlers.GetTaskHandler).Methods("GET", "POST")

	log.Println("Server started on :8080")

	go agent.StartWorker()

	log.Fatal(http.ListenAndServe(":8080", r))
}
