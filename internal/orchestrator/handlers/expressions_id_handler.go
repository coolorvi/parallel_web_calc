package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	vars := mux.Vars(r)
	id := vars["id"]

	expr, exists := Expressions[id]
	if !exists {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"id":     expr.ID,
		"status": expr.Status,
	}

	if expr.Result != nil {
		response["result"] = *expr.Result
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expression": response,
	})
}
