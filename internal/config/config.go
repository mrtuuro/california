package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBUri     string
	DatabaseName   string
	CollectionName string
	HttpAddr       string
}

func NewConfig() *Config {
	err := godotenv.Load("/Users/tozay/go/src/california/dev.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		MongoDBUri:     os.Getenv("MONGO_DB_CONNECTION_URI"),
		DatabaseName:   os.Getenv("MONGO_DATABASE_NAME"),
		CollectionName: os.Getenv("MONGO_USERS_COLLECTION_NAME"),
		HttpAddr:       os.Getenv("HTTP_ADDRESS"),
	}

}
