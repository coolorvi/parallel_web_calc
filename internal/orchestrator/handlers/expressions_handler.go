package handlers

import (
	"encoding/json"
	"net/http"
)

var Expressions = make(map[string]*Expression)

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Unsupported method", http.StatusUnprocessableEntity)
		return
	}

	response := map[string][]Expression{"expressions": {}}
	for _, expr := range Expressions {
		response["expressions"] = append(response["expressions"], *expr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
