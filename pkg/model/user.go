package model

type UserType int

const (
	Admin UserType = iota + 1
	Normal
	Premium
)

type User struct {
	Name         string   `bson:"Name" json:"name"`
	PhoneNumber  string   `bson:"PhoneNumber" json:"phone_number"`
	Email        string   `bson:"Email" json:"email"`
	UserType     UserType `bson:"UserType" json:"user_type"`
	Vehicle      Vehicle  `bson:"Vehicle" json:"vehicle"`
	RefreshToken string   `bson:"RefreshToken" json:"refresh_token",omitempty`
}
