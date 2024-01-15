package navigationsvc

import (
	"context"

	"california/pkg/model"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CalculateTripEndpoint endpoint.Endpoint
	RecommendEndpoint     endpoint.Endpoint
}

func MakeServerEndpoints(c context.Context, s NavigationService) Endpoints {
	return Endpoints{
		CalculateTripEndpoint: MakeCalculateTripEndpoint(c, s),
		RecommendEndpoint:     MakeRecommendEndpoint(c, s),
	}
}

type BaseResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func MakeRecommendEndpoint(c context.Context, s NavigationService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.RecommendRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)
		advice, e := s.Recommend(c, &req)
		if e != nil {
			return BaseResponse{
				Message: "failed",
				Data:    nil,
			}, e
		}
		return BaseResponse{
			Message: "success",
			Data: recommendResponse{
				Advice: advice,
				Err:    e,
			},
		}, nil
	}
}

type recommendRequest struct {
	Context context.Context
	Stops   []model.Stop
}

type recommendResponse struct {
	*BaseResponse
	Advice []*model.Advice `json:"recommendations,omitempty"`
	Err    error
}

func (r recommendResponse) Failed() error { return r.Err }

func MakeCalculateTripEndpoint(c context.Context, s NavigationService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(calculateTripRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)
		tripInfo, e := s.CalculateTrip(c, req)
		if e != nil {
			return calculateTripResponse{
				Err: e,
			}, e
		}
		return BaseResponse{
			Message: "success",
			Data: calculateTripResponse{
				TripInfo: tripInfo,
				Err:      e,
			},
		}, nil
	}
}

type calculateTripRequest struct {
	Context  context.Context
	Distance float64
}

type calculateTripResponse struct {
	*BaseResponse
	TripInfo []*model.TripInfo
	Err      error
}

func (r calculateTripResponse) Failed() error { return r.Err }
