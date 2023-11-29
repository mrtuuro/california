package usersvc

import (
	"context"

	"california/pkg/model"
	"github.com/go-kit/kit/endpoint"
)

type BaseResponse struct {
	Message string `json:"message,omitempty"`
	Data    Data   `json:"data,omitempty"`
}

type Data struct {
	registerResponse
}

type EndPoints struct {
	RegisterEndpoint        endpoint.Endpoint
	LoginEndpoint           endpoint.Endpoint
	VehicleRegisterEndpoint endpoint.Endpoint
	GetMeEndpoint           endpoint.Endpoint
	UpdateUserEndpoint      endpoint.Endpoint
	UpdateVehicleEndpoint   endpoint.Endpoint
	GetUsersEndpoint        endpoint.Endpoint
}

func MakeServerEndpoints(c context.Context, s UserService) EndPoints {
	return EndPoints{
		RegisterEndpoint:        MakeRegisterEndpoint(s),
		LoginEndpoint:           MakeLoginEndpoint(s),
		VehicleRegisterEndpoint: MakeVehicleRegisterEndpoint(c, s),
		GetMeEndpoint:           MakeGetMeEndpoint(c, s),
		UpdateUserEndpoint:      MakeUpdateUserEndpoint(c, s),
		UpdateVehicleEndpoint:   MakeUpdateVehicleEndpoint(c, s),
		GetUsersEndpoint:        MakeListAllUsersEndpoint(c, s),
	}
}

//
//func (e EndPoints) Register(ctx context.Context, user *model.User) error {
//	request := registerRequest{User: user}
//	response, err := e.RegisterEndpoint(ctx, request)
//	if err != nil {
//		return err
//	}
//	resp := response.(registerResponse)
//	return resp.Err
//}

func MakeRegisterEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerRequest)
		insertedUser, e := s.Register(ctx, req.User)
		if e != nil {
			baseErrResponse := &BaseResponse{
				Message: e.Error(),
				Data:    Data{},
			}
			return registerResponse{
				BaseResponse: baseErrResponse,
			}, nil
		}

		// Generate data struct only for registerResponse.
		data := Data{
			registerResponse: registerResponse{
				UserType: insertedUser.UserType,
				Token:    req.User.RefreshToken,
				Err:      e,
			},
		}
		return BaseResponse{
			Message: "success",
			Data:    data,
		}, nil

		//return registerResponse{
		//	UserType: insertedUser.UserType,
		//	Token:    req.User.RefreshToken,
		//	Err:      e,
		//}, nil
	}
}

// registerRequest is used to decode json request body of register endpoint's.
type registerRequest struct {
	User *model.User
}

// registerResponse is used to encode json response body of register endpoint's.
type registerResponse struct {
	*BaseResponse
	UserType model.UserType `json:"user_type,omitempty"`
	Token    string         `json:"token,omitempty"`
	Err      error          `json:"err,omitempty"`
}

func (e registerResponse) error() error { return e.Err }

func MakeLoginEndpoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(loginRequest)
		user, e := s.Login(ctx, req.Email, req.Password)
		if e != nil {
			return loginResponse{
				Err: e,
			}, nil
		}
		return loginResponse{
			UserType: user.UserType,
			Token:    user.RefreshToken,
			Err:      e,
		}, nil
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	UserType model.UserType `json:"user_type,omitempty"`
	Token    string         `json:"token,omitempty"`
	Err      error          `json:"err,omitempty"`
}

func (e loginResponse) error() error { return e.Err }

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

type vehicleRegisterRequest struct {
	Context context.Context
	Vehicle *model.Vehicle
}

type vehicleRegisterResponse struct {
	Err error `json:"err,omitempty"`
}

func (e vehicleRegisterResponse) error() error { return e.Err }

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

type getMeRequest struct {
	Context context.Context
}

type getMeResponse struct {
	User *model.User `json:"user,omitempty"`
	Err  error       `json:"err,omitempty"`
}

func (e getMeResponse) error() error { return e.Err }

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

type updateUserRequest struct {
	Context context.Context
	User    *model.User
}

type updateUserResponse struct {
	Err error `json:"err,omitempty"`
}

func (e updateUserResponse) error() error { return e.Err }

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

func MakeListAllUsersEndpoint(c context.Context, s UserService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listAllUsersRequest)
		jwt := req.Context.Value("jwt")
		c = context.WithValue(c, "Authorization", jwt)
		users, e := s.ListAllUsers(c)
		return listAllUsersResponse{
			Users: users,
			Err:   e,
		}, nil
	}
}

type listAllUsersRequest struct {
	Context context.Context
}

type listAllUsersResponse struct {
	Users []*model.User `json:"users,omitempty"`
	Err   error         `json:"err,omitempty"`
}

func (e listAllUsersResponse) error() error { return e.Err }
