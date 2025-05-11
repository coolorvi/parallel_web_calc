package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (o *Orchestrator) GetExpression(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	o.mu.Lock()
	defer o.mu.Unlock()

	expr, ok := o.expressions[id]
	if !ok {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]*Expression{"expression": expr})
}
