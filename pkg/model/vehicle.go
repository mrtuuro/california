package model

type Vehicle struct {
	Brand              string     `bson:"Brand" json:"brand"`
	Model              string     `bson:"Model" json:"model"`
	EngineType         EngineType `bson:"EngineType" json:"engine_type"`                 // Diesel, Petrol, Hybrid
	EngineSize         float64    `bson:"EngineSize" json:"engine_size"`                 // 1.6L/2.0L
	AverageConsumption float64    `bson:"AverageConsumption" json:"average_consumption"` // 4.5/100km
}

type EngineType int

const (
	Petrol EngineType = iota + 1
	Diesel
	Hybrid
	Electric
)
