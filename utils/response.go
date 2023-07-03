package utils

import (
	"encoding/json"
	"net/http"
)

func SuccessResponse(w http.ResponseWriter, data interface{}) {
	type SuccessResponse struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}
	json.NewEncoder(w).Encode(SuccessResponse{Status: "ok", Data: data})
}

func ErrorResponse(w http.ResponseWriter, message string) {
	type SuccessResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	json.NewEncoder(w).Encode(SuccessResponse{Status: "ok", Message: message})
}
