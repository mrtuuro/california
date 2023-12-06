package charge_stationsvc

import (
	"context"

	"california/pkg/model"
	"california/pkg/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StationService interface {
	StationRegister(context.Context, *model.Station) (insertedStation *model.Station, err error)
	GetStations(context.Context) (stations []*model.Station, err error)
	UpdateStation(ctx context.Context, station *model.Station, stationId string) (err error)
	RemoveStation(ctx context.Context, stationId string) (err error)
	SearchStation(ctx context.Context, brandName string) (stations []*model.Station, err error)
}

type chargeStationService struct {
	store repository.Store
}

func (s *chargeStationService) StationRegister(ctx context.Context, station *model.Station) (*model.Station, error) {
	oidStr := primitive.NewObjectID().Hex()
	station.ID, _ = primitive.ObjectIDFromHex(oidStr)
	insertedStation, err := s.store.InsertStation(ctx, station)
	if err != nil {
		return nil, err
	}
	return insertedStation, nil
}

func (s *chargeStationService) GetStations(ctx context.Context) ([]*model.Station, error) {
	stations, err := s.store.GetAllStations(ctx)
	if err != nil {
		return nil, err
	}
	return stations, nil
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

func NewStationService(store repository.Store) StationService {
	return &chargeStationService{
		store: store,
	}
}
