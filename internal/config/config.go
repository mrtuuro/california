package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBUri   string
	DatabaseName string

	UsersCollectionName    string
	StationsCollectionName string

	UsersHttpAddr    string
	StationsHttpAddr string
}

func NewConfig() *Config {
	err := godotenv.Load("./dev.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		MongoDBUri:             os.Getenv("MONGO_DB_CONNECTION_URI"),
		DatabaseName:           os.Getenv("MONGO_DATABASE_NAME"),
		UsersCollectionName:    os.Getenv("MONGO_USERS_COLLECTION_NAME"),
		StationsCollectionName: os.Getenv("MONGO_STATIONS_COLLECTION_NAME"),
		UsersHttpAddr:          os.Getenv("USER_HTTP_ADDRESS"),
		StationsHttpAddr:       os.Getenv("STATIONS_HTTP_ADDRESS"),
	}

}
