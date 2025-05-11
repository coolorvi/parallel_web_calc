package start

import (
	"fmt"
	"log"
	"net/http"

	"github.com/coolorvi/parallel_web_calc/internal/auth"
	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/handlers"
	"github.com/coolorvi/parallel_web_calc/internal/storage"
	"github.com/gorilla/mux"
)

func Web(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

func StartServer() {
	err := storage.InitDB("data.db")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	o := handlers.NewOrchestrator(storage.DB)
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/register", auth.Register(storage.DB)).Methods("POST")
	r.HandleFunc("/api/v1/login", auth.Login(storage.DB)).Methods("POST")

	r.HandleFunc("/api/v1/calculate", o.AddExpression).Methods("POST")
	r.HandleFunc("/api/v1/expressions", o.GetExpressions).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", o.GetExpression).Methods("GET")
	r.HandleFunc("/internal/task", o.GetTask).Methods("GET")
	r.HandleFunc("/internal/task", o.PostResult).Methods("POST")
	r.HandleFunc("/", Web).Methods("GET")

	er := http.ListenAndServe(":8080", r)
	if er != nil {
		fmt.Println(err)
	}
}
