package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/olad5/go-url-shortener/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	port := os.Getenv("PORT")
	fmt.Printf("Server Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router.Initialize()))
}
