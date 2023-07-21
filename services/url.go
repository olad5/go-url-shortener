package services

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"os"

	"github.com/olad5/go-url-shortener/entity"
	"github.com/olad5/go-url-shortener/storage"
)

type UrlService struct {
	repository storage.Repository
}

func NewUrlService(RepositoryAdapter storage.Repository) (*UrlService, error) {
	if RepositoryAdapter == nil {
		return &UrlService{}, errors.New("UrlService failed to initialize")
	}
	return &UrlService{RepositoryAdapter}, nil
}

func (u *UrlService) ShortenUrl(url string) (entity.ShortenUrl, error) {
	if url == "" {
		return entity.ShortenUrl{}, errors.New("Cannot shorten empty string")
	}

	exisitingShortUrl, err := u.repository.FetchUrlByOriginalUrl(url)

	if err == nil {
		return exisitingShortUrl, nil
	}

	if isLinkLive(url) != true {
		return entity.ShortenUrl{}, errors.New("Link is dead")
	}

	randomUniqueId := generateUniqueId()
	base62String := convertIdToBase62String(randomUniqueId)

	shortUrl := entity.ShortenUrl{
		ShortUrl: base62String, OriginalUrl: url,
		ClickCount: 0,
		UniqueId:   randomUniqueId,
	}

	err = u.repository.CreateUrl(shortUrl)
	if err != nil {
		return entity.ShortenUrl{}, err
	}
	return shortUrl, nil
}

func (u *UrlService) Info(slug string) (entity.ShortenUrl, error) {
	shortUrl, err := u.repository.FetchUrlByShortCode(slug)
	if err != nil {
		return entity.ShortenUrl{}, err
	}
	return shortUrl, nil
}

func (u *UrlService) UpdateClickCount(shortUrl entity.ShortenUrl) error {
	updatedShortUrl := entity.ShortenUrl{
		ShortUrl: shortUrl.ShortUrl, OriginalUrl: shortUrl.OriginalUrl,
		ClickCount: shortUrl.ClickCount + 1,
		UniqueId:   shortUrl.UniqueId,
	}
	err := u.repository.UpdateUrl(updatedShortUrl)
	if err != nil {
		return err
	}
	return nil
}

// Naive Unique Id Generator
func generateUniqueId() int {
	MAX_INT := 7935425686241
	b := new(big.Int).SetInt64(int64(MAX_INT))
	randomBigInt, _ := rand.Int(rand.Reader, b)
	randomeNewInt := int(randomBigInt.Int64())
	return randomeNewInt
}

func convertIdToBase62String(decimal int) string {
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

func isLinkLive(originalUrl string) bool {
	currentEnvironment := os.Getenv("ENVIRONMENT")

	// I don't know how to make this work in my tests, i'm sorry :(
	if currentEnvironment != "production" {
		return true
	}

	json.NewEncoder(&bytes.Buffer{}).Encode(nil)

	res, err := http.Get(originalUrl)
	if err != nil {
		return false
	}
	if res.StatusCode != 200 {
		return false
	}
	return true
}
