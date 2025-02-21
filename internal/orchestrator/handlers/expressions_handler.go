package handlers

import (
	"encoding/json"
	"net/http"
)

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	expressionsList := make([]Expression, 0, len(expressions))

	for _, expr := range expressions {
		expressionsList = append(expressionsList, *expr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{"expressions": expressionsList}); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
