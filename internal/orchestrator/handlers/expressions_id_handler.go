package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	id := parts[4]

	expr, exists := Expressions[id]
	if !exists {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if expr == nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]*Expression{"expression": expr}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
