package model

import (
	"context"
)

type TripInfo struct {
	Speed              float64 `json:"speed"`
	AverageConsumption float64 `json:"average_consumption"` // 4.5/100km
	Distance           float64 `json:"distance"`
	TotalPrice         float64 `json:"total_price"` // 432.54 TL
}

type Advice struct {
	Stops []Stop `json:"stops"`
}

type RecommendRequest struct {
	Context context.Context
	Stops   []Stop `json:"stops"`
}

type Stop struct {
	Name  string `json:"name"`
	Lat   string `json:"lat"`
	Long  string `json:"long"`
	Color string `json:"color"`
}
