package repository

import (
	"context"
	"fmt"
	"log"

	"california/internal/config"
	"california/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	InsertUser(ctx context.Context, user *model.User) error
	UserExists(ctx context.Context, email, phoneNumber string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type MongoStore struct {
	Client *mongo.Client
	Coll   *mongo.Collection
}

func NewMongoStore(cfg *config.Config) *MongoStore {
	client, coll := ConnectDB(cfg.MongoDBUri, cfg.DatabaseName, cfg.CollectionName)
	return &MongoStore{
		Client: client,
		Coll:   coll,
	}
}

func (s *MongoStore) InsertUser(ctx context.Context, user *model.User) error {
	_, err := s.Coll.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoStore) UserExists(ctx context.Context, email, phoneNumber string) (bool, error) {
	filter := bson.M{"$or": []bson.M{{"Email": email}, {"PhoneNumber": phoneNumber}}}
	count, err := s.Coll.CountDocuments(context.Background(), filter)
	return count > 0, err

}

func (s *MongoStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	filter := bson.M{"Email": email}
	err := s.Coll.FindOne(context.Background(), filter).Decode(&user)
	if err != nil && err == mongo.ErrNoDocuments {
		return nil, mongo.ErrNoDocuments
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func ConnectDB(dbUri, dbName, collectionName string) (*mongo.Client, *mongo.Collection) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}
	coll := client.Database(dbName).Collection(collectionName)
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
		return nil, nil
	}
	fmt.Println("Connected to Database!")
	return client, coll
}
