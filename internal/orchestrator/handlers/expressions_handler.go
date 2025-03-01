package handlers

import (
	"encoding/json"
	"net/http"
)

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var expressionsList []map[string]interface{}
	for _, expr := range Expressions {
		exprData := map[string]interface{}{
			"id":     expr.ID,
			"status": expr.Status,
		}

		if expr.Result != nil {
			exprData["result"] = *expr.Result
		}

		expressionsList = append(expressionsList, exprData)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expressions": expressionsList,
	})
}
