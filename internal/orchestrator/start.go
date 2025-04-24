package orchestrator

import (
	"database/sql"
	"log"
	"net/http"

	handler "github.com/coolorvi/parallel_web_calc/internal/auth"
	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/handlers"
	"github.com/coolorvi/parallel_web_calc/internal/storage"
	"github.com/gorilla/mux"
)

var DB *sql.DB

func Start() {
	r := mux.NewRouter()

	auth := &handler.Server{DB: storage.DB}
	r.HandleFunc("/api/v1/register", auth.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/v1/login", auth.LoginHandler).Methods("POST")

	r.HandleFunc("/api/v1/calculate", handlers.CalculateHandler).Methods("POST")
	r.HandleFunc("/api/v1/expressions", handlers.ExpressionsHandler).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", handlers.ExpressionHandler).Methods("GET")
	r.HandleFunc("/internal/task", handlers.GetTaskHandler).Methods("GET", "POST")

	log.Println("Server started on :8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
