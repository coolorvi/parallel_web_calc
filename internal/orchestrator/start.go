package orchestrator

import (
	"log"
	"net/http"

	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/handlers"
	"github.com/gorilla/mux"
)

func Start() {
	r := mux.NewRouter()
	r.HandleFunc("/api/vi/calculate", handlers.CalculateHandler).Methods("Post")
	http.HandleFunc("/api/v1/expressions", handlers.ExpressionsHandler)
	http.HandleFunc("/api/v1/expressions/", handlers.ExpressionHandler)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
