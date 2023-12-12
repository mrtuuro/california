package model

type TripInfo struct {
	Speed              float64
	AverageConsumption float64 `bson:"AverageConsumption" json:"average_consumption"` // 4.5/100km
	Distance           float64
}
