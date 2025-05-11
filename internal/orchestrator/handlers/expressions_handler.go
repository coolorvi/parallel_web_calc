package handlers

import (
	"encoding/json"
	"net/http"
)

func (o *Orchestrator) GetExpressions(w http.ResponseWriter, r *http.Request) {
	o.mu.Lock()
	defer o.mu.Unlock()

	resp := struct {
		Expressions []*Expression `json:"expressions"`
	}{Expressions: make([]*Expression, 0, len(o.expressions))}
	for _, expr := range o.expressions {
		resp.Expressions = append(resp.Expressions, expr)
	}
	json.NewEncoder(w).Encode(resp)
}
