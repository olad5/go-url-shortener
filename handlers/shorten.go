package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/olad5/go-url-shortener/services"
	"github.com/olad5/go-url-shortener/utils"
)

func Shorten(w http.ResponseWriter, r *http.Request) {
	var originalUrl utils.RequestBody
	err := json.NewDecoder(r.Body).Decode(&originalUrl)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, "Something wrong with the json passed", http.StatusBadRequest)
		return
	}

	urlService, err := services.NewUrlService()
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, utils.ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}
	shortUrl, err := urlService.ShortenUrl(originalUrl.Url)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.SuccessResponse(w, shortUrl)
}
