package services

import (
	"context"
	"crypto/rand"
	"math/big"
	"os"

	"github.com/olad5/go-url-shortener/entity"
	"github.com/olad5/go-url-shortener/storage"
	"github.com/olad5/go-url-shortener/storage/mongo"
)

type UrlService struct {
	repository storage.Repository
}

func NewUrlService() (*UrlService, error) {
	mongo, err := mongo.New(context.Background(), os.Getenv("MONGO_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}
	service := UrlService{mongo}
	return &service, nil
}

func (u *UrlService) ShortenUrl(url string) (entity.ShortenUrl, error) {
	randomUniqueId := generateUniqueId()
	base62String := convertIdToBase62String(randomUniqueId)
	shortUrl := entity.ShortenUrl{
		ShortUrl: base62String, OriginalUrl: url,
		ClickCount: 0,
		UniqueId:   randomUniqueId,
	}
	err := u.repository.SaveUrl(shortUrl)
	if err != nil {
		return entity.ShortenUrl{}, err
	}
	return shortUrl, nil
}

func (u *UrlService) Info(slug string) (entity.ShortenUrl, error) {
	shortUrl, err := u.repository.FetchUrl(slug)
	if err != nil {
		return entity.ShortenUrl{}, err
	}
	return shortUrl, nil
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
func isLinkValid() {}
func isLinkDead()  {}
