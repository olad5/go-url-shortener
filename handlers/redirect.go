package handlers

import (
	"net/http"

	"github.com/olad5/go-url-shortener/config"
	"github.com/olad5/go-url-shortener/services"
	"github.com/olad5/go-url-shortener/utils"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	slug := utils.GetField(r, 0)
	urlService, err := services.NewUrlService(config.RepositoryAdapter)
	if err != nil {
		utils.ErrorResponse(w, utils.ErrSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	shortUrl, err := urlService.Info(slug)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	err = urlService.UpdateClickCount(shortUrl)
	if err != nil {
		utils.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, shortUrl.OriginalUrl, http.StatusTemporaryRedirect)
}
