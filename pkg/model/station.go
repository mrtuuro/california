package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Station struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Brand       string             `bson:"Brand" json:"brand"`
	Latitude    float64            `bson:"Latitude" json:"latitude"`
	Longitude   float64            `bson:"Longitude" json:"longitude"`
	Status      int                `bson:"Status" json:"status"`
	CurrentType int                `bson:"CurrentType" json:"current_type"`
	Distance    float64            `bson:"Distance" json:"distance"`
	Address     string             `bson:"Address" json:"address"`
	Sockets     []Socket           `bson:"Sockets" json:"sockets"`
}

type Socket struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Name        string             `bson:"Name" json:"name"` // Bu field bağlı olduğu istastonun Brand'ine eşit.
	KW          float64            `bson:"KW" json:"kw"`
	CurrentType CurrentType        `bson:"CurrentType" json:"current_type"`
	Price       float64            `bson:"Price" json:"price"`
	SocketType  string             `bson:"SocketType" json:"socket_type"`
	Status      SocketStatus       `bson:"Status" json:"status"`
}

type CurrentType int

const (
	DC CurrentType = iota + 1
	AC
)

type SocketStatus int

const (
	UnAvailable SocketStatus = iota
	Available
)
