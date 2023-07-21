package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/olad5/go-url-shortener/entity"
	"github.com/olad5/go-url-shortener/storage"
	"github.com/olad5/go-url-shortener/utils"
)

type RedisCache struct {
	Client *redis.Client
}

var (
	shortUrlPrefixKey      = "short-"
	originalUrlPrefixKey   = "original-"
	ttl                    = time.Hour * 24
	contextTimeoutDuration = 3 * time.Second
)

func New(ctx context.Context, address string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		Client: client,
	}, nil
}

func (r *RedisCache) CreateUrl(shortUrl entity.ShortenUrl) error {
	ctx, cancel := context.WithTimeout(utils.TodoContext, contextTimeoutDuration)
	defer cancel()

	json, err := json.Marshal(shortUrl)
	if err != nil {
		return err
	}

	originalUrlHash := getSHA256HashOfOriginalUrl(shortUrl.OriginalUrl)

	_, err = r.Client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, prefixOriginalUrlHash(originalUrlHash), json, ttl)
		pipe.Set(ctx, prefixShortCode(shortUrl.ShortUrl), json, ttl)
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) UpdateUrl(shortUrl entity.ShortenUrl) error {
	ctx, cancel := context.WithTimeout(utils.TodoContext, contextTimeoutDuration)
	defer cancel()
	json, err := json.Marshal(shortUrl)
	if err != nil {
		return err
	}

	originalUrlHash := getSHA256HashOfOriginalUrl(shortUrl.OriginalUrl)

	txf := func(tx *redis.Tx) error {
		_, err = tx.Pipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, prefixOriginalUrlHash(originalUrlHash), json, ttl)
			pipe.Set(ctx, prefixShortCode(shortUrl.ShortUrl), json, ttl)
			return nil
		})

		return err
	}

	err = r.Client.Watch(ctx, txf, prefixOriginalUrlHash(originalUrlHash), prefixShortCode(shortUrl.ShortUrl))
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) FetchUrlByShortCode(shortCode string) (entity.ShortenUrl, error) {
	ctx, cancel := context.WithTimeout(utils.TodoContext, contextTimeoutDuration)
	defer cancel()

	result, err := r.Client.Get(ctx, prefixShortCode(shortCode)).Result()
	shortUrl := entity.ShortenUrl{}
	err = json.Unmarshal([]byte(result), &shortUrl)
	if err != nil {
		return entity.ShortenUrl{}, errors.New(utils.ErrRecordNotFound)
	}

	return shortUrl, nil
}

func (r *RedisCache) FetchUrlByOriginalUrl(originalUrl string) (entity.ShortenUrl, error) {
	ctx, cancel := context.WithTimeout(utils.TodoContext, contextTimeoutDuration)
	defer cancel()

	originalUrlHash := getSHA256HashOfOriginalUrl(originalUrl)
	result, err := r.Client.Get(ctx, prefixOriginalUrlHash(originalUrlHash)).Result()
	shortUrl := entity.ShortenUrl{}
	err = json.Unmarshal([]byte(result), &shortUrl)
	if err != nil {
		return entity.ShortenUrl{}, errors.New(utils.ErrRecordNotFound)
	}

	return shortUrl, nil
}

func (r *RedisCache) Ping() storage.DataSourceHealth {
	if err := r.Client.Ping(utils.TodoContext).Err(); err != nil {
		return storage.DEGRADED
	}
	return storage.OK
}

func getSHA256HashOfOriginalUrl(originalUrl string) string {
	hashString := sha256.New()
	hashString.Write([]byte(originalUrl))
	hashBytes := hashString.Sum(nil)
	encodedString := hex.EncodeToString(hashBytes)
	return encodedString
}

func prefixOriginalUrlHash(hash string) string {
	return originalUrlPrefixKey + hash
}

func prefixShortCode(shortcode string) string {
	return shortUrlPrefixKey + shortcode
}
