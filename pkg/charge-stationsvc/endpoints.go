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
	StationRegisterEndpoint   endpoint.Endpoint
	GetAllStationsEndpoint    endpoint.Endpoint
	UpdateStationInfoEndpoint endpoint.Endpoint
	RemoveStationEndpoint     endpoint.Endpoint
	SearchStationEndpoint     endpoint.Endpoint
	ListBrandsEndpoint        endpoint.Endpoint
	ListSocketsEndpoint       endpoint.Endpoint
	FilterStationsEndpoint    endpoint.Endpoint
}

func MakeServerEndpoints(c context.Context, s StationService) StationEndpoints {
	return StationEndpoints{
		StationRegisterEndpoint:   MakeRegisterStationEndpoint(c, s),
		GetAllStationsEndpoint:    MakeGetAllStationsEndpoint(c, s),
		UpdateStationInfoEndpoint: MakeUpdateStationInfoEndpoint(c, s),
		RemoveStationEndpoint:     MakeRemoveStationEndpoint(c, s),
		SearchStationEndpoint:     MakeSearchStationEndpoint(c, s),
		ListBrandsEndpoint:        MakeListBrandsEndpoint(c, s),
		ListSocketsEndpoint:       MakeListSocketsEndpoint(c, s),
		FilterStationsEndpoint:    MakeFilterStationsEndpoint(c, s),
	}
}

func MakeFilterStationsEndpoint(c context.Context, s StationService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(filterStationsRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)

		stations, e := s.FilterStation(c, req.BrandNames, req.SocketNames, req.CurrentType)
		if e != nil {
			return filterStationsResponse{
				Err: e,
			}, e
		}
		return BaseResponse{
			Message: "success",
			Data: filterStationsResponse{
				Stations: stations,
				Err:      e,
			},
		}, nil
	}
}

type filterStationsRequest struct {
	Context     context.Context
	BrandNames  []string
	SocketNames []string
	CurrentType int
}

type filterStationsResponse struct {
	*BaseResponse
	Stations []*model.Station `json:"stations,omitempty"`
	Err      error            `json:"err,omitempty"`
}

func (r filterStationsResponse) Failed() error { return r.Err }

func MakeListSocketsEndpoint(c context.Context, s StationService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listSocketsRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)

		sockets, e := s.ListSockets(c)
		if e != nil {
			return listSocketsResponse{
				Err: e,
			}, e
		}
		return BaseResponse{
			Message: "success",
			Data: listSocketsResponse{
				Sockets: sockets,
				Err:     e,
			},
		}, nil
	}
}

type listSocketsRequest struct {
	Context context.Context
}
type listSocketsResponse struct {
	*BaseResponse
	Sockets []*model.Socket `json:"sockets,omitempty"`
	Err     error           `json:"err,omitempty"`
}

func (r listSocketsResponse) Failed() error { return r.Err }

func MakeListBrandsEndpoint(c context.Context, s StationService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listBrandsRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)

		brands, e := s.ListBrands(c)
		if e != nil {
			return listBrandsResponse{
				Err: e,
			}, e
		}
		return BaseResponse{
			Message: "success",
			Data: listBrandsResponse{
				Brands: brands,
				Err:    e,
			},
		}, nil
	}
}

type listBrandsRequest struct {
	Context context.Context
}

type listBrandsResponse struct {
	*BaseResponse
	Brands []string `json:"brands,omitempty"`
	Err    error    `json:"err,omitempty"`
}

func (r listBrandsResponse) Failed() error { return r.Err }

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

type registerStationRequest struct {
	Context context.Context
	Station *model.Station
	Sockets []*model.Socket
}

type registerStationResponse struct {
	*BaseResponse
	Station *model.Station `json:"insertedStation,omitempty"`
	Err     error          `json:"err,omitempty"`
}

func (r registerStationResponse) Failed() error { return r.Err }

func MakeGetAllStationsEndpoint(c context.Context, s StationService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getAllStationsRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)

		stations, e := s.GetStations(c)
		if e != nil {
			return getAllStationsResponse{
				Err: e,
			}, e
		}
		return BaseResponse{
			Message: "success",
			Data: getAllStationsResponse{
				Stations: stations,
				Err:      e,
			},
		}, nil
	}
}

type getAllStationsRequest struct {
	Context context.Context
}

type getAllStationsResponse struct {
	*BaseResponse
	Stations []*model.Station `json:"stations,omitempty"`
	Err      error            `json:"err,omitempty"`
}

func (r getAllStationsResponse) Failed() error { return r.Err }

func MakeUpdateStationInfoEndpoint(c context.Context, s StationService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateStationInfoRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)

		e := s.UpdateStation(c, req.Station, req.StationID)
		if e != nil {
			return updateStationInfoResponse{
				Err: e,
			}, e
		}
		return BaseResponse{
			Message: "success",
			Data: updateStationInfoResponse{
				Err: e,
			},
		}, nil
	}
}

type updateStationInfoRequest struct {
	Context   context.Context
	StationID string
	Station   *model.Station
}

type updateStationInfoResponse struct {
	*BaseResponse
	Err error `json:"err,omitempty"`
}

func (r updateStationInfoResponse) Failed() error { return r.Err }

func MakeRemoveStationEndpoint(c context.Context, s StationService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(removeStationRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)

		e := s.RemoveStation(c, req.StationID)
		if e != nil {
			return removeStationResponse{
				Err: e,
			}, e
		}
		return BaseResponse{
			Message: "success",
			Data: removeStationResponse{
				Err: e,
			},
		}, nil
	}
}

type removeStationRequest struct {
	Context   context.Context
	StationID string
}

type removeStationResponse struct {
	*BaseResponse
	Err error `json:"err,omitempty"`
}

func (r removeStationResponse) Failed() error { return r.Err }

func MakeSearchStationEndpoint(c context.Context, s StationService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(searchStationRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)

		stations, e := s.SearchStation(c, req.Brand)
		if e != nil {
			return searchStationResponse{
				Err: e,
			}, e
		}
		return BaseResponse{
			Message: "success",
			Data: searchStationResponse{
				Stations: stations,
				Err:      e,
			},
		}, nil
	}
}

type searchStationRequest struct {
	Context context.Context
	Station *model.Station
	Brand   string
}

type searchStationResponse struct {
	*BaseResponse
	Stations []*model.Station `json:"stations,omitempty"`
	Err      error            `json:"err,omitempty"`
}

func (r searchStationResponse) Failed() error { return r.Err }
