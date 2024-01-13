package main

import (
	"cache-go/application/cache"
	routes "cache-go/application/router"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		errorMessage := fmt.Sprintf("Error loading .env file: %v", err)
		log.Fatal(errorMessage)
	}

	// Initialize Redis
	cache.InitRedis()
	defer func() {
		if err := cache.RedisClient.Close(); err != nil {
			log.Println("Error closing Redis:", err)
		}
	}()

	// Run server
	routes.RunServer()
	fmt.Println("Server started on port: 9000")
}
