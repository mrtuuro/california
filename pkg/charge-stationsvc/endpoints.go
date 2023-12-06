package charge_stationsvc

import (
	"context"

	"california/pkg/model"
	"github.com/go-kit/kit/endpoint"
)

type BaseResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type StationEndpoints struct {
	StationRegisterEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(c context.Context, s StationService) StationEndpoints {
	return StationEndpoints{
		StationRegisterEndpoint: MakeRegisterStationEndpoint(c, s),
	}
}

type registerStationRequest struct {
	Context context.Context
	Station *model.Station
}

type registerStationResponse struct {
	*BaseResponse
	Station *model.Station `json:"insertedStation,omitempty"`
	Err     error          `json:"err,omitempty"`
}

func (r registerStationResponse) Failed() error { return r.Err }

func MakeRegisterStationEndpoint(c context.Context, s StationService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerStationRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)
		insertedStation, e := s.StationRegister(c, req.Station)
		if e != nil {
			return registerStationResponse{
				Err: e,
			}, e
		}
		return BaseResponse{
			Message: "success",
			Data: registerStationResponse{
				Station: insertedStation,
				Err:     e,
			},
		}, nil
	}
}
