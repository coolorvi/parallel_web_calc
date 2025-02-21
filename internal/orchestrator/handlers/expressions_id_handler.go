package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")

	expr, exists := Expressions[id]
	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]*Expression{"expression": expr}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
