package handlers

import (
	"encoding/json"
	"net/http"
)

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	if len(Expressions) == 0 {
		http.Error(w, "No expressions found", http.StatusNotFound)
		return
	}

	exprList := make([]*Expression, 0, len(Expressions))
	for _, expr := range Expressions {
		exprList = append(exprList, expr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exprList)
}
