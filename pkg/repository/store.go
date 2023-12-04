package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"california/internal/config"
	"california/internal/helpers"
	"california/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	InsertUser(ctx context.Context, user *model.User) (*model.User, error)
	UserExists(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	InsertVehicleToUser(ctx context.Context, user *model.User, vehicle *model.Vehicle) error
	FindUsersByFilter(ctx context.Context, filter bson.M) ([]*model.User, error)
	UpdateUser(ctx context.Context, reqUser *model.User) error
	UpdateVehicle(ctx context.Context, reqVehicle *model.Vehicle) error
	GetAllUsers(ctx context.Context) ([]*model.User, error)
}

type MongoStore struct {
	Client       *mongo.Client
	UsersColl    *mongo.Collection
	StationsColl *mongo.Collection
}

func NewMongoStore(cfg *config.Config) *MongoStore {
	client := ConnectDB(cfg.MongoDBUri)
	userColl := GetCollection(client, cfg.DatabaseName, cfg.UsersCollectionName)
	stationsColl := GetCollection(client, cfg.DatabaseName, cfg.StationsCollectionName)
	return &MongoStore{
		Client:       client,
		UsersColl:    userColl,
		StationsColl: stationsColl,
	}
}

func (s *MongoStore) InsertUser(_ context.Context, user *model.User) (*model.User, error) {
	var insertedUser *model.User
	insertRes, err := s.UsersColl.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	insertedIdStr := insertRes.InsertedID.(primitive.ObjectID)
	if err = s.UsersColl.FindOne(context.Background(), bson.M{"_id": insertedIdStr}).Decode(&insertedUser); err != nil {
		return nil, err
	}
	return insertedUser, nil
}

func (s *MongoStore) UserExists(_ context.Context, email string) (bool, error) {
	filter := bson.M{"$or": []bson.M{{"Email": email}}}
	count, err := s.UsersColl.CountDocuments(context.Background(), filter)
	return count > 0, err

}

func (s *MongoStore) GetUserByEmail(_ context.Context, email string) (*model.User, error) {
	var user model.User
	filter := bson.M{"Email": email}
	err := s.UsersColl.FindOne(context.Background(), filter).Decode(&user)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return nil, mongo.ErrNoDocuments
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoStore) InsertVehicleToUser(_ context.Context, user *model.User, vehicle *model.Vehicle) error {
	filter := bson.M{"Email": user.Email}
	update := bson.M{"$set": bson.M{"Vehicle": vehicle}}
	_, err := s.UsersColl.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoStore) UpdateUser(ctx context.Context, reqUser *model.User) error {
	userId := ctx.Value("userId").(string)
	oid, _ := primitive.ObjectIDFromHex(userId)
	fmt.Println(oid)
	fmt.Println(userId)

	filter := bson.M{"id": oid}
	update := bson.M{"$set": bson.M{"Name": reqUser.Name}}
	if reqUser.Password != "" {
		newHashedPass, err := helpers.HashRegisterPassword(reqUser.Password)
		if err != nil {
			return err
		}
		update = bson.M{"$set": bson.M{"Name": reqUser.Name, "Password": newHashedPass}}
	}

	_, err := s.UsersColl.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoStore) UpdateVehicle(ctx context.Context, reqVehicle *model.Vehicle) error {
	userId := ctx.Value("userId").(string)
	oid, _ := primitive.ObjectIDFromHex(userId)

	filter := bson.M{"id": oid}
	update := bson.M{"$set": bson.M{"Vehicle": reqVehicle}}
	_, err := s.UsersColl.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoStore) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	cursor, err := s.UsersColl.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (s *MongoStore) FindUsersByFilter(ctx context.Context, filter bson.M) ([]*model.User, error) {
	var users []*model.User
	cursor, err := s.UsersColl.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func ConnectDB(dbUri string) *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println("Connected to Database!")
	return client
}

func GetCollection(client *mongo.Client, dbName, collectionName string) *mongo.Collection {
	coll := client.Database(dbName).Collection(collectionName)
	return coll
}
