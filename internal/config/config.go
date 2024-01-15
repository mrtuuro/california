package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBUri   string
	DatabaseName string

	UsersCollectionName    string
	StationsCollectionName string
	SocketsCollectionName  string

	UsersHttpAddr      string
	StationsHttpAddr   string
	NavigationHttpAddr string
	AuthHttpAddr       string
}

func NewConfig() *Config {
	err := godotenv.Load("./dev.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	return &Config{
		MongoDBUri:   os.Getenv("MONGO_DB_CONNECTION_URI"),
		DatabaseName: os.Getenv("MONGO_DATABASE_NAME"),

		UsersCollectionName:    os.Getenv("MONGO_USERS_COLLECTION_NAME"),
		StationsCollectionName: os.Getenv("MONGO_STATIONS_COLLECTION_NAME"),
		SocketsCollectionName:  os.Getenv("MONGO_SOCKETS_COLLECTION_NAME"),

		UsersHttpAddr:      os.Getenv("USER_HTTP_ADDRESS"),
		StationsHttpAddr:   os.Getenv("STATIONS_HTTP_ADDRESS"),
		NavigationHttpAddr: os.Getenv("NAVIGATION_HTTP_ADDRESS"),
		AuthHttpAddr:       os.Getenv("AUTH_HTTP_ADDRESS"),
	}

}
