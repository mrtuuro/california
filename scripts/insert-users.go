package main

import (
	"context"
	"fmt"
	"log"

	"california/internal/config"
	model2 "california/pkg/model"
	"california/pkg/repository"
)

func main() {
	cfg := config.NewConfig()
	mongoStore := repository.NewMongoStore(cfg)
	insertOneResult, err := mongoStore.Coll.InsertOne(context.Background(), &model2.User{
		Name:        "John Doe",
		PhoneNumber: "123456789",
		Email:       "",
		UserType:    model2.Normal,
		Vehicle: model2.Vehicle{
			Brand:              "BMW",
			Model:              "X5",
			EngineType:         model2.Petrol,
			EngineSize:         2.0,
			AverageConsumption: 4.5,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertOneResult.InsertedID)
}
