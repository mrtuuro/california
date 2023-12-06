package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Station struct {
	ID          primitive.ObjectID `bson:"id,omitempty" json:"id,omitempty"` // This id is created by mongo and stored as 'id'
	Name        string             `bson:"Name" json:"name"`
	Latitude    float64            `bson:"Latitude" json:"latitude"`
	Longitude   float64            `bson:"Longitude" json:"longitude"`
	Status      int                `bson:"Status" json:"status"`            // Silinebilir
	CurrentType int                `bson:"CurrentType" json:"current_type"` // Silinebilir
	Distance    float64            `bson:"Distance" json:"distance"`        // Silinebilir
	Address     string             `bson:"Address" json:"address"`          // Silinebilir
	Sockets     []Socket           `bson:"Sockets" json:"sockets"`
}

type Socket struct {
	KW          float64 `bson:"KW" json:"kw,"`                    // Silinebilir
	CurrentType int     `bson:"CurrentType" json:"current_type,"` // Silinebilir
	Price       float64 `bson:"Price" json:"price,omitempty"`     // Silinebilir
	SocketType  int     `bson:"SocketType" json:"socket_type,"`
	Status      int     `bson:"Status" json:"status,"`
}
