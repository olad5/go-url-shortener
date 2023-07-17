package storage

import (
	"github.com/olad5/go-url-shortener/entity"
)

type Repository interface {
	CreateUrl(shortUrl entity.ShortenUrl) error
	UpdateUrl(shortUrl entity.ShortenUrl) error
	FetchUrl(id string) (entity.ShortenUrl, error)
}
