package navigationsvc

import (
	"context"
	"math"

	"california/pkg/model"
	"california/pkg/repository"
)

const (
	patrolPrice = 33.02
	dieselPrice = 35.47
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
		totalPrice := calculateTotalPrice(avgConsumption, userVehicle.EngineType)
		tripInfo = append(tripInfo, &model.TripInfo{
			Speed:              speed,
			AverageConsumption: roundResult(avgConsumption),
			Distance:           req.Distance,
			TotalPrice:         roundResult(totalPrice),
		})
	}
	return tripInfo, nil
}

func calculateTotalPrice(consumption float64, engineType model.EngineType) float64 {
	switch engineType {
	case 1:
		return consumption * patrolPrice
	case 2:
		return consumption * dieselPrice
	default:
		return consumption * patrolPrice
	}
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

func roundResult(result float64) float64 {
	return math.Round(result*100) / 100
}
