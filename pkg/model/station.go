package model

type Station struct {
	Brand       string   `bson:"Brand" json:"brand"`
	Latitude    float64  `bson:"Latitude" json:"latitude"`
	Longitude   float64  `bson:"Longitude" json:"longitude"`
	Status      int      `bson:"Status" json:"status"`
	CurrentType int      `bson:"CurrentType" json:"current_type"`
	Distance    float64  `bson:"Distance" json:"distance"`
	Address     string   `bson:"Address" json:"address"`
	Sockets     []Socket `bson:"Sockets" json:"sockets"`
}

type Socket struct {
	Name        string       `bson:"Name" json:"name"` // Bu field bağlı olduğu istastonun Brand'ine eşit.
	KW          float64      `bson:"KW" json:"kw"`
	CurrentType CurrentType  `bson:"CurrentType" json:"current_type"`
	Price       float64      `bson:"Price" json:"price"`
	SocketType  string       `bson:"SocketType" json:"socket_type"`
	Status      SocketStatus `bson:"Status" json:"status"`
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
