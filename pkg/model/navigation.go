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
	Number int    `json:"advice"`
	Stops  []Stop `json:"stops"`
}

type RecommendRequest struct {
	Context      context.Context
	Distance     int        `json:"distance"`
	StartPoint   Coordinate `json:"start_point"`
	ArrivalPoint Coordinate `json:"arrival_point"`
	Stops        []Stop     `json:"stops"`
}

type Coordinate struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type Stop struct {
	Name  string  `json:"name"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
	Color string  `json:"color"`
}

func (s *Stop) DetermineColor(increment int) {
	if increment == 10 {
		s.Color = "green"
	} else if increment == 20 {
		s.Color = "blue"
	} else if increment == 30 {
		s.Color = "red"
	} else if increment == 40 {
		s.Color = "red"
	} else if increment == 50 {
		s.Color = "red"
	} else {
		s.Color = "red"
	}
}
