package charge_stationsvc

import (
	"context"
	"errors"

	"california/pkg/model"
	"california/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StationService interface {
	StationRegister(context.Context, *model.Station) (insertedStation *model.Station, err error)
	GetStations(context.Context) (stations []*model.Station, err error)
	GetStation(ctx context.Context, stationId string) (station *model.Station, err error)
	UpdateStation(ctx context.Context, station *model.Station, stationId string) (err error)
	RemoveStation(ctx context.Context, stationId string) (err error)
	DeleteSocket(ctx context.Context, socketId string) (err error)
	SearchStation(ctx context.Context, brandName string) (stations []*model.Station, err error)
	ListBrands(ctx context.Context) (brands []string, err error)
	ListSockets(ctx context.Context) (sockets []*model.Socket, err error)
	FilterStation(ctx context.Context, brandName []string, socketType []string, currentType int) (stations []*model.Station, err error)
}

type chargeStationService struct {
	store repository.Store
}

func (s *chargeStationService) DeleteSocket(ctx context.Context, socketId string) (err error) {
	err = s.store.DeleteSocket(ctx, socketId)
	if err != nil {
		return err
	}
	return nil
}

func (s *chargeStationService) StationRegister(ctx context.Context, station *model.Station) (*model.Station, error) {
	lat := station.Latitude
	long := station.Longitude

	locationFilter := bson.M{
		"Latitude":  lat,
		"Longitude": long,
	}
	existedStations, err := s.store.FindStationByFilter(ctx, locationFilter)
	if len(existedStations) == 0 || errors.Is(err, mongo.ErrNoDocuments) {
		for i := range station.Sockets {
			station.Sockets[i].ID = primitive.NewObjectID()
		}

		station.ID = primitive.NewObjectID()
		insertedStation, err := s.store.InsertStation(ctx, station)
		if err != nil {
			return nil, err
		}

		for i, _ := range station.Sockets {
			err = s.store.InsertSocket(ctx, &station.Sockets[i])
			if err != nil {
				return nil, err
			}
		}

		return insertedStation, nil
	}

	for i := range station.Sockets {
		station.Sockets[i].ID = primitive.NewObjectID()
		err = s.store.InsertSocket(ctx, &station.Sockets[i])
		if err != nil {
			return nil, err
		}

		err := s.store.PushSocketToStation(ctx, existedStations[0], station.Sockets[i])
		if err != nil {
			return nil, err
		}
	}
	return existedStations[0], nil
}

func (s *chargeStationService) GetStations(ctx context.Context) ([]*model.Station, error) {
	stations, err := s.store.GetAllStations(ctx)
	if err != nil {
		return nil, err
	}
	return stations, nil
}

func (s *chargeStationService) GetStation(ctx context.Context, stationId string) (station *model.Station, err error) {
	station, err = s.store.GetStationById(ctx, stationId)
	if err != nil {
		return nil, err
	}
	return station, nil
}

func (s *chargeStationService) UpdateStation(ctx context.Context, station *model.Station, stationId string) (err error) {
	err = s.store.UpdateStationInfo(ctx, station, stationId)
	if err != nil {
		return err
	}
	return nil
}

func (s *chargeStationService) RemoveStation(ctx context.Context, stationId string) (err error) {
	err = s.store.DeleteStation(ctx, stationId)
	if err != nil {
		return err
	}
	return nil
}

func (s *chargeStationService) SearchStation(ctx context.Context, brandName string) (stations []*model.Station, err error) {
	filter := bson.M{"Brand": bson.M{"$regex": primitive.Regex{Pattern: brandName, Options: "i"}}}
	stations, err = s.store.FindStationByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return stations, nil
}

func (s *chargeStationService) ListBrands(ctx context.Context) (brands []string, err error) {
	allStations, err := s.store.GetAllStations(ctx)
	if err != nil {
		return nil, err
	}

	var allBrandsWithDuplicates []string
	for _, station := range allStations {
		allBrandsWithDuplicates = append(allBrandsWithDuplicates, station.Brand)
	}

	brands = removeDuplicates(allBrandsWithDuplicates)

	return brands, nil
}

func (s *chargeStationService) ListSockets(ctx context.Context) (sockets []*model.Socket, err error) {
	sockets, err = s.store.ListSockets(ctx)
	if err != nil {
		return nil, err
	}
	return sockets, nil
}

func (s *chargeStationService) FilterStation(ctx context.Context, brandNames []string, socketNames []string, currentType int) (stations []*model.Station, err error) {

	filter := bson.M{}

	// Brand name'e göre filtreleme
	if len(brandNames) > 0 {
		filter["Brand"] = bson.M{
			"$in": brandNames,
		}
	}

	// socketNames'e göre filtreleme
	if len(socketNames) > 0 {
		filter["Sockets"] = bson.M{
			"$elemMatch": bson.M{
				"Name": bson.M{
					"$in": socketNames,
				},
			},
		}
	}

	// currentType'a göre filtreleme
	if currentType == 0 || currentType == 1 {
		if elemMatch, ok := filter["Sockets"].(bson.M); ok {
			elemMatch["$elemMatch"].(bson.M)["CurrentType"] = currentType
		} else {
			filter["Sockets"] = bson.M{
				"$elemMatch": bson.M{
					"CurrentType": currentType,
				},
			}
		}
	}

	//filter := bson.M{
	//	"Brand": bson.M{
	//		"$in": brandNames, // birden fazla brandName kabul et
	//	},
	//	"Sockets": bson.M{
	//		"$elemMatch": bson.M{
	//			"Name": bson.M{
	//				"$in": socketNames, // birden fazla socketName kabul et
	//			},
	//			"CurrentType": currentType,
	//		},
	//	},
	//}
	stations, err = s.store.FindStationByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}
	return stations, nil
}

func removeDuplicates(duplicates []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range duplicates {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func NewStationService(store repository.Store) StationService {
	return &chargeStationService{
		store: store,
	}
}
