package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserType int

const (
	Admin UserType = iota + 1
	Normal
	Premium
)

type User struct {
	ID           primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"` // This id is created by mongo and stored as 'id'
	Name         string             `bson:"Name" json:"name"`
	Email        string             `bson:"Email" json:"email"`
	Password     string             `bson:"Password" json:"password"` // Store the password as a hash
	UserType     UserType           `bson:"UserType" json:"user_type"`
	Vehicle      Vehicle            `bson:"Vehicle" json:"vehicle"`
	RefreshToken string             `bson:"RefreshToken" json:"refresh_token",omitempty`
}
