package handlers

import (
	"encoding/json"
	"net/http"
)

var Expressions = make(map[string]*Expression)

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string][]Expression{"expressions": {}}
	for _, expr := range Expressions {
		response["expressions"] = append(response["expressions"], *expr)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
