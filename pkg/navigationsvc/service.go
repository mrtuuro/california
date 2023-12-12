package navigationsvc

import (
	"context"
	"math"

	"california/pkg/model"
	"california/pkg/repository"
)

type NavigationService interface {
	CalculateTrip(ctx context.Context, req calculateTripRequest) (tripInfo []*model.TripInfo, err error)
}

type navigationService struct {
	store repository.Store
}

func NewNavigationService(store repository.Store) NavigationService {
	return &navigationService{
		store: store,
	}
}

func (s *navigationService) CalculateTrip(ctx context.Context, req calculateTripRequest) (tripInfo []*model.TripInfo, err error) {
	userEmail := ctx.Value("email").(string)
	user, err := s.store.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	var userVehicle *model.Vehicle
	userVehicle = &user.Vehicle

	speeds := []float64{60, 80, 90, 100, 110, 120, 150, 200}
	for _, speed := range speeds {
		avgConsumption := calculateFuelConsumption(userVehicle.EngineType, userVehicle.EngineSize, userVehicle.AverageConsumption, req.Distance, speed)
		tripInfo = append(tripInfo, &model.TripInfo{
			Speed:              speed,
			AverageConsumption: avgConsumption,
			Distance:           req.Distance,
		})
	}
	return tripInfo, nil
}

func calculateFuelConsumption(engineType model.EngineType, engineSize, averageConsumption, distance, speed float64) float64 {
	// Base speed factor for all engines
	speedFactor := math.Pow(speed/80, 1.2)

	// Engine size factor for petrol and diesel engines
	engineSizeFactor := engineSize / 2.0 // Assuming 2.0 liters as a baseline for comparison

	switch engineType {
	case 1:
		return averageConsumption * engineSizeFactor * speedFactor * distance / 100
	case 2:
		dieselEfficiencyModifier := 0.85
		return averageConsumption * dieselEfficiencyModifier * engineSizeFactor * speedFactor * distance / 100
	case 3:
		hybridEfficiencyModifier := 0.75
		return averageConsumption * hybridEfficiencyModifier * engineSizeFactor * speedFactor * distance / 100
	default:
		return averageConsumption * engineSizeFactor * speedFactor * distance / 100
	}

}
