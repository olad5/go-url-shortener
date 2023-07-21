package repository_adapter

import (
	"log"

	"github.com/olad5/go-url-shortener/entity"
	"github.com/olad5/go-url-shortener/storage"
	"github.com/olad5/go-url-shortener/storage/mongo"
	"github.com/olad5/go-url-shortener/storage/redis"
	"github.com/olad5/go-url-shortener/utils"
)

type RepositoryAdapter struct {
	database mongo.MongoRepository
	cache    redis.RedisCache
}

func NewRespositoryAdapter(mongoConnectionString, redisConnectionString string) (*RepositoryAdapter, error) {
	mongo, err := mongo.New(utils.TodoBackground, mongoConnectionString)
	if err != nil {
		return nil, err
	}

	redisCache, err := redis.New(utils.TodoBackground, redisConnectionString)
	if err != nil {
		log.Println("Redis Failed to Initialize", err)
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

func (r *RepositoryAdapter) Ping() storage.DataSourceHealth {
	if mongoHealth := r.database.Ping(); mongoHealth != storage.OK {
		return storage.DOWN
	}
	if redisHealth := r.cache.Ping(); redisHealth != storage.OK {
		return storage.DEGRADED
	}
	return storage.OK
}
