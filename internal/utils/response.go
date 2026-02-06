package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"todo-api/internal/domain"
)

func DTOResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	DTOResponse(w, statusCode, map[string]string{"error": message})
}

func ParseError(w http.ResponseWriter, err error) {
	if appErr, ok := err.(*domain.AppError); ok {
		ErrorResponse(w, appErr.StatusCode, appErr.Message)
		return
	}
	log.Printf("Unhandled error: %v", err)
	ErrorResponse(w, http.StatusInternalServerError, "An unexpected error occurred")
}
