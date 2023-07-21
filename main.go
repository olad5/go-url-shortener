package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/olad5/go-url-shortener/config"
	"github.com/olad5/go-url-shortener/router"
	"github.com/olad5/go-url-shortener/storage/repository_adapter"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error Loading .env file")
	}

	mongoConnectionString := os.Getenv("MONGO_CONNECTION_STRING")
	redisConnectionString := os.Getenv("REDIS_ADDRESS")
	port := os.Getenv("PORT")

	repository, err := repository_adapter.NewRespositoryAdapter(mongoConnectionString, redisConnectionString)
	config.RepositoryAdapter = repository

	if err != nil {
		log.Fatal("Failed to Initialize Database")
	}

	log.Printf("Server Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router.Initialize()))
}
