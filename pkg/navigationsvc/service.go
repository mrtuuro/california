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

	earthRadius = 6371
)

type NavigationService interface {
	CalculateTrip(ctx context.Context, req calculateTripRequest) (tripInfo []*model.TripInfo, err error)
	Recommend(ctx context.Context, req *model.RecommendRequest) (advices []*model.Advice, err error)
}

type navigationService struct {
	store repository.Store
}

func NewNavigationService(store repository.Store) NavigationService {
	return &navigationService{
		store: store,
	}
}

func (s *navigationService) Recommend(ctx context.Context, rec *model.RecommendRequest) (advices []*model.Advice, err error) {
	var (
		allStops         = rec.Stops
		startPoint       = rec.StartPoint
		arrivalPoint     = rec.ArrivalPoint
		_                = haversineDistance(startPoint.Lat, startPoint.Long, arrivalPoint.Lat, arrivalPoint.Long)
		realDistance     = rec.Distance
		startStopDistMap = make(map[string]float64)
		totalAdviceCount = 3
	)

	for _, stop := range allStops {
		distBetweenStartAndStop := haversineDistance(startPoint.Lat, startPoint.Long, stop.Lat, stop.Long)
		startStopDistMap[stop.Name] = distBetweenStartAndStop
		//startStopDistList = append(startStopDistList, distBetweenStartAndStop)
	}
	//sort.Float64s(startStopDistList)

	// Kaç tane durak önerilecek?
	// Yaklaşık her 300 km'de 1 tane durak önerilecek.
	// Eğer toplam yol 600 km ise 1 durak önerilecek. Ve 300 300 artarak gidecek
	// İlk kontrol edilecek durak 200 km'den uzak olmalı. 200 km'den uzak olan duraklar arasından en yakın olanı seçilecek.
	// Eğer toplam yol 0-600 km arasındaysa 1 durak önerilecek. Ve bu durak yaklaşık %50lik kısımda olacak. +- 5km.
	totalStopCount := int(float64(realDistance / 300))
	if realDistance <= 650 {
		totalStopCount = 1
	}

	if totalStopCount == 1 {
		fmt.Println(startStopDistMap)
		for i := 0; i < totalAdviceCount; i++ {
			var advice model.Advice
			advice.Number = i + 1
			stopPoint := realDistance / 2
			increment := 10

			//for stopName, dist := range startStopDistMap {
			//	if dist > float64(stopPoint-15) && dist < float64(stopPoint+15) {
			//		for _, stop := range allStops {
			//			if stop.Name == stopName {
			//				advice.Stops = append(advice.Stops, model.Stop{
			//					Name:  stopName,
			//					Lat:   stop.Lat,
			//					Long:  stop.Long,
			//					Color: "red",
			//				})
			//			}
			//		}
			//		break
			//	}
			//}
			for {
				found := false
				for stopName, dist := range startStopDistMap {
					if dist > float64(stopPoint-increment) && dist < float64(stopPoint+increment) {
						for _, stop := range allStops {
							if stop.Name == stopName {
								advice.Stops = append(advice.Stops, model.Stop{
									Name:  stopName,
									Lat:   stop.Lat,
									Long:  stop.Long,
									Color: "red",
								})
								found = true // Durak bulundu.
								break        // İç döngüyü kır.
							}
						}
					}
					if found {
						break // Dış döngüyü kır.
					}
				}

				if found {
					break // Durak bulundu, ana döngüyü kır.
				} else {
					increment += 10     // Durak bulunamadı, aralığı genişlet.
					if increment > 50 { // Maksimum aralığa ulaştıysa döngüyü sonlandır.
						fmt.Println("Maksimum aralık aşıldı, durak bulunamadı.")
						break
					}
				}
			}
			advices = append(advices, &advice)
		}
		return advices, nil
	}

	for i := 1; i <= totalStopCount; i++ {
		stopPoint := 300 * i
		for j := 1; j <= totalAdviceCount; j++ {
			var advice model.Advice
			advice.Number = j
			increment := 10

			//for stopName, dist := range startStopDistMap {
			//	if dist > float64(stopPoint-20) && dist < float64(stopPoint+20) {
			//		fmt.Println(stopName, dist)
			//		for _, stop := range allStops {
			//			if stop.Name == stopName {
			//				advice.Stops = append(advice.Stops, model.Stop{
			//					Name:  stopName,
			//					Lat:   stop.Lat,
			//					Long:  stop.Long,
			//					Color: "red",
			//				})
			//			}
			//		}
			//		break
			//	}
			//}
			for {
				found := false
				for stopName, dist := range startStopDistMap {
					if dist > float64(stopPoint-increment) && dist < float64(stopPoint+increment) {
						for _, stop := range allStops {
							if stop.Name == stopName {
								advice.Stops = append(advice.Stops, model.Stop{
									Name:  stopName,
									Lat:   stop.Lat,
									Long:  stop.Long,
									Color: "red",
								})
								found = true // Durak bulundu.
								break        // İç döngüyü kır.
							}
						}
					}
					if found {
						break // Dış döngüyü kır.
					}
				}

				if found {
					break // Durak bulundu, ana döngüyü kır.
				} else {
					increment += 10     // Durak bulunamadı, aralığı genişlet.
					if increment > 50 { // Maksimum aralığa ulaştıysa döngüyü sonlandır.
						fmt.Println("Maksimum aralık aşıldı, durak bulunamadı.")
						break
					}
				}
			}
			advices = append(advices, &advice)
		}
	}
	return advices, nil
}

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	lat1 = degreesToRadians(lat1)
	lon1 = degreesToRadians(lon1)
	lat2 = degreesToRadians(lat2)
	lon2 = degreesToRadians(lon2)

	// calculate haversine
	lat := lat2 - lat1
	lon := lon2 - lon1
	a := math.Pow(math.Sin(lat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(lon/2), 2)
	c := 2 * math.Asin(math.Sqrt(a))
	return earthRadius * c
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
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
