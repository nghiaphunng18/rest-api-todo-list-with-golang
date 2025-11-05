package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func DTOResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode JSON response: %v", err)
	}
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	DTOResponse(w, statusCode, map[string]string{"error": message})
}
