package navigationsvc

import (
	"context"
	"fmt"
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
	Recommend(ctx context.Context, req *model.RecommendRequest) (recommendation []*model.Advice, err error)
}

type navigationService struct {
	store repository.Store
}

func NewNavigationService(store repository.Store) NavigationService {
	return &navigationService{
		store: store,
	}
}

func (s *navigationService) Recommend(ctx context.Context, rec *model.RecommendRequest) (advice []*model.Advice, err error) {
	allStops := rec.Stops
	for _, stop := range allStops {
		fmt.Printf("Stop Name: %s\n", stop.Name)
		fmt.Printf("Stop Long: %s\n", stop.Long)
		fmt.Printf("Stop Lat: %s\n", stop.Lat)
	}
	advice = append(advice, &model.Advice{
		Stops: []model.Stop{
			{
				Name:  allStops[0].Name,
				Lat:   allStops[0].Lat,
				Long:  allStops[0].Long,
				Color: "red",
			},
			{
				Name:  allStops[3].Name,
				Lat:   allStops[3].Lat,
				Long:  allStops[3].Long,
				Color: "green",
			},
			{
				Name:  allStops[5].Name,
				Lat:   allStops[5].Lat,
				Long:  allStops[5].Long,
				Color: "blue",
			},
		},
	})
	return advice, nil
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
	//speedFactor := math.Pow(speed/80, 1.2)
	kilometers := distance / 1000

	// Engine size factor for petrol and diesel engines
	//engineSizeFactor := engineSize / 2.0 // Assuming 2.0 liters as a baseline for comparison

	//switch engineType {
	//case 1:
	//	return averageConsumption * engineSizeFactor * speedFactor * kilometers / 100
	//case 2:
	//	dieselEfficiencyModifier := 0.85
	//	return averageConsumption * dieselEfficiencyModifier * engineSizeFactor * speedFactor * kilometers / 100
	//case 3:
	//	hybridEfficiencyModifier := 0.75
	//	return averageConsumption * hybridEfficiencyModifier * engineSizeFactor * speedFactor * kilometers / 100
	//default:
	//	return averageConsumption * engineSizeFactor * speedFactor * kilometers / 100
	//}

	switch engineType {
	case 1:
		return averageConsumption * kilometers / 100
	case 2:
		dieselEfficiencyModifier := 0.85
		return averageConsumption * dieselEfficiencyModifier * kilometers / 100
	case 3:
		hybridEfficiencyModifier := 0.75
		return averageConsumption * hybridEfficiencyModifier * kilometers / 100
	default:
		return averageConsumption * kilometers / 100
	}

}

func roundResult(result float64) float64 {
	return math.Round(result*100) / 100
}
