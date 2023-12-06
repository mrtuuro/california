package charge_stationsvc

import (
	"context"

	"california/pkg/model"
	"california/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StationService interface {
	StationRegister(context.Context, *model.Station) (insertedStation *model.Station, err error)
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

func NewStationService(store repository.Store) StationService {
	return &chargeStationService{
		store: store,
	}
}
