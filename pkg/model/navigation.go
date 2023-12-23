package model

type TripInfo struct {
	Speed              float64 `json:"speed"`
	AverageConsumption float64 `json:"average_consumption"` // 4.5/100km
	Distance           float64 `json:"distance"`
	TotalPrice         float64 `json:"total_price"` // 432.54 TL
}
