package http

import (
	"encoding/json"
	"net/http"
)

// WriteJSON encodes the data into JSON, sets the content-type header, and writes the response.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader((status))
	json.NewEncoder(w).Encode(data)
}

// WriteError sends a structured JSON error message
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}
