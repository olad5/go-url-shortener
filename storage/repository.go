package storage

import (
	"context"
	"log"
	"os"

	"github.com/olad5/go-url-shortener/entity"
	"github.com/olad5/go-url-shortener/storage/mongo"
	"github.com/olad5/go-url-shortener/storage/redis"
)

type Repository interface {
	CreateUrl(shortUrl entity.ShortenUrl) error
	UpdateUrl(shortUrl entity.ShortenUrl) error
	FetchUrlByShortCode(shortCode string) (entity.ShortenUrl, error)
	FetchUrlByOriginalUrl(originalUrl string) (entity.ShortenUrl, error)
}

type RepositoryAdapter struct {
	database mongo.MongoRepository
	cache    redis.RedisCache
}

func NewRespositoryAdapter() (*RepositoryAdapter, error) {
	mongo, err := mongo.New(context.Background(), os.Getenv("MONGO_CONNECTION_STRING"))
	if err != nil {
		return nil, err
	}

	redisCache, err := redis.New(context.Background(), os.Getenv("REDIS_ADDRESS"))
	if err != nil {
		log.Println(err)
	}

	return &RepositoryAdapter{database: *mongo, cache: *redisCache}, nil
}

func (r *RepositoryAdapter) CreateUrl(shortUrl entity.ShortenUrl) error {
	err := r.database.CreateUrl(shortUrl)
	if err != nil {
		return err
	}

	err = r.cache.CreateUrl(shortUrl)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (r *RepositoryAdapter) UpdateUrl(shortUrl entity.ShortenUrl) error {
	err := r.database.UpdateUrl(shortUrl)
	if err != nil {
		return err
	}

	err = r.cache.UpdateUrl(shortUrl)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (r *RepositoryAdapter) FetchUrlByShortCode(shortCode string) (entity.ShortenUrl, error) {
	shortUrlFromCache, err := r.cache.FetchUrlByShortCode(shortCode)
	if err == nil {
		return shortUrlFromCache, nil
	}

	shortUrlFromDB, err := r.database.FetchUrlByShortCode(shortCode)
	if err != nil {
		log.Println(err)
		return entity.ShortenUrl{}, err
	}

	return shortUrlFromDB, nil
}

func (r *RepositoryAdapter) FetchUrlByOriginalUrl(originalUrl string) (entity.ShortenUrl, error) {
	shortUrlFromCache, err := r.cache.FetchUrlByOriginalUrl(originalUrl)
	if err == nil {
		return shortUrlFromCache, nil
	}

	shortUrlFromDB, err := r.database.FetchUrlByOriginalUrl(originalUrl)
	if err != nil {
		return entity.ShortenUrl{}, err
	}

	return shortUrlFromDB, nil
}
