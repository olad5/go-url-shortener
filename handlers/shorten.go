package handlers

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"net/http"

	"github.com/olad5/go-url-shortener/utils"
)

func isLinkValid() {}
func isLinkDead()  {}

func Shorten(w http.ResponseWriter, r *http.Request) {
	type shortenLinkResponse struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	type requestUrl struct {
		Url string `json:"url"`
	}

	var originalUrl requestUrl
	err := json.NewDecoder(r.Body).Decode(&originalUrl)
	if err != nil {
		// TODO: handle this error
		log.Println("handle this error")
	}
	randomUniqueId := generateUniqueId()
	convertedString := convertToBase62(randomUniqueId)
	shortUrl := ShortenUrl{ShortUrl: convertedString, OriginalUrl: originalUrl.Url, UniqueId: randomUniqueId}
	utils.SuccessResponse(w, shortUrl)
}

var MAX_INT = 7935425686241

type ShortenUrl struct {
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
	UniqueId    int
}

// Naive Unique Id Generator
func generateUniqueId() int {
	MAX_INT := 7935425686241
	b := new(big.Int).SetInt64(int64(MAX_INT))
	randomBigInt, _ := rand.Int(rand.Reader, b)
	randomeNewInt := int(randomBigInt.Int64())
	return randomeNewInt
}

func convertToBase62(decimal int) string {
	base62Chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base := len(base62Chars)
	converted := ""

	for decimal > 0 {
		remainder := decimal % base
		converted = string(base62Chars[remainder]) + converted
		decimal /= base
	}

	return converted
}
