package storage

import (
	"github.com/olad5/go-url-shortener/entity"
)

type Repository interface {
	CreateUrl(shortUrl entity.ShortenUrl) error
	UpdateUrl(shortUrl entity.ShortenUrl) error
	FetchUrlByShortCode(shortCode string) (entity.ShortenUrl, error)
	FetchUrlByOriginalUrl(originalUrl string) (entity.ShortenUrl, error)
	Ping() DataSourceHealth
}
