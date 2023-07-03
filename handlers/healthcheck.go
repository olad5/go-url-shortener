package handlers

import (
	"net/http"

	"github.com/olad5/go-url-shortener/utils"
)

func Healthcheck(w http.ResponseWriter, r *http.Request) {
	type healthCheckResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	data := map[string]interface{}{"message": "healthcheck complete"}
	utils.SuccessResponse(w, data)
}
