package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func SuccessResponse(w http.ResponseWriter, data interface{}) {
	type SuccessResponse struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(SuccessResponse{Status: "ok", Data: data}); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	type SuccessResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(SuccessResponse{Status: "error", Message: message}); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}
