// Package httputil provides HTTP utility functions for request and response handling.
package httputil

import (
	"encoding/json"
	"net/http"
)

// MakeJSONResponse writes a JSON response with the given status code and data
func MakeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
