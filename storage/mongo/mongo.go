package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/olad5/go-url-shortener/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	db            *mongo.Database
	urlCollection *mongo.Collection
}

func New(ctx context.Context, connectionString string) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}
	db := client.Database("url-shortener")
	urlCollection := db.Collection("urls")
	repo := MongoRepository{
		db:            db,
		urlCollection: urlCollection,
	}

	return &repo, nil
}

func (u *MongoRepository) CreateUrl(shortUrl entity.ShortenUrl) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := u.urlCollection.InsertOne(ctx, shortUrl)
	if err != nil {
		return err
	}

	return nil
}

func (u *MongoRepository) UpdateUrl(shortUrl entity.ShortenUrl) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"unique_id": shortUrl.UniqueId}
	updatedDoc := bson.M{
		"$set": shortUrl,
	}

	_, err := u.urlCollection.UpdateOne(ctx, filter, updatedDoc)
	if err != nil {
		return err
	}

	return nil
}

func (u *MongoRepository) FetchUrl(shortCode string) (entity.ShortenUrl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := u.urlCollection.FindOne(ctx, bson.M{"short_url": shortCode})
	shortUrl := entity.ShortenUrl{}
	err := result.Decode(&shortUrl)
	if err != nil {
		return entity.ShortenUrl{}, errors.New("Record was not found")
	}

	return shortUrl, nil
}
