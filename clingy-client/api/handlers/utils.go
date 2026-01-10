package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func decodeJSONBody(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		endpoint := r.Method + " " + r.URL.Path
		log.Printf("API: %s - Error decoding JSON: %v", endpoint, err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return false
	}
	return true
}
