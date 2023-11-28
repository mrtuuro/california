package usersvc

import (
	"context"

	"california/pkg/model"
	"github.com/go-kit/kit/endpoint"
)

type EndPoints struct {
	RegisterEndpoint        endpoint.Endpoint
	LoginEndpoint           endpoint.Endpoint
	VehicleRegisterEndpoint endpoint.Endpoint
	GetMeEndpoint           endpoint.Endpoint
	UpdateUserEndpoint      endpoint.Endpoint
	UpdateVehicleEndpoint   endpoint.Endpoint
}

func MakeServerEndpoints(c context.Context, s UserService) EndPoints {
	return EndPoints{
		RegisterEndpoint:        MakeRegisterEndpoint(s),
		LoginEndpoint:           MakeLoginEndpoint(s),
		VehicleRegisterEndpoint: MakeVehicleRegisterEndpoint(c, s),
		GetMeEndpoint:           MakeGetMeEndpoint(c, s),
		UpdateUserEndpoint:      MakeUpdateUserEndpoint(c, s),
		UpdateVehicleEndpoint:   MakeUpdateVehicleEndpoint(c, s),
	}
}

func (e EndPoints) Register(ctx context.Context, user *model.User) error {
	request := registerRequest{User: user}
	response, err := e.RegisterEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(registerResponse)
	return resp.Err
}

func MakeRegisterEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerRequest)
		e := s.Register(ctx, req.User)
		return registerResponse{
			Token: req.User.RefreshToken,
			Err:   e,
		}, nil
	}
}

func MakeLoginEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loginRequest)
		token, e := s.Login(ctx, req.Email, req.Password)
		return loginResponse{
			Token: token,
			Err:   e,
		}, nil
	}
}

func MakeVehicleRegisterEndpoint(c context.Context, s UserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(vehicleRegisterRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)
		e := s.VehicleRegister(c, req.Vehicle)
		return vehicleRegisterResponse{
			Err: e,
		}, nil
	}
}

func MakeGetMeEndpoint(c context.Context, s UserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getMeRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)
		user, e := s.GetMe(c)
		return getMeResponse{
			User: user,
			Err:  e,
		}, nil
	}
}

func MakeUpdateUserEndpoint(c context.Context, s UserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateUserRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)
		e := s.UpdateUserInfo(c, req.User)
		return updateUserResponse{
			Err: e,
		}, nil
	}
}

func MakeUpdateVehicleEndpoint(c context.Context, s UserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateVehicleRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)
		e := s.UpdateVehicleInfo(c, req.Vehicle)
		return updateVehicleResponse{
			Err: e,
		}, nil
	}
}

type updateVehicleRequest struct {
	Context context.Context
	Vehicle *model.Vehicle
}

type updateVehicleResponse struct {
	Err error `json:"err,omitempty"`
}

func (e updateVehicleResponse) error() error { return e.Err }

type updateUserRequest struct {
	Context context.Context
	User    *model.User
}

type updateUserResponse struct {
	Err error `json:"err,omitempty"`
}

func (e updateUserResponse) error() error { return e.Err }

type getMeRequest struct {
	Context context.Context
}

type getMeResponse struct {
	User *model.User `json:"user,omitempty"`
	Err  error       `json:"err,omitempty"`
}

func (e getMeResponse) error() error { return e.Err }

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token,omitempty"`
	Err   error  `json:"err,omitempty"`
}

func (e loginResponse) error() error { return e.Err }

// registerRequest is used to decode json request body of register endpoint's.
type registerRequest struct {
	User *model.User
}

// registerResponse is used to encode json response body of register endpoint's.
type registerResponse struct {
	Token string `json:"token,omitempty"`
	Err   error  `json:"err,omitempty"`
}

func (e registerResponse) error() error { return e.Err }

type vehicleRegisterRequest struct {
	Context context.Context
	Vehicle *model.Vehicle
}

type vehicleRegisterResponse struct {
	Err error `json:"err,omitempty"`
}

func (e vehicleRegisterResponse) error() error { return e.Err }
