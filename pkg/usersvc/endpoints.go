package usersvc

import (
	"context"

	"california/pkg/model"
	"github.com/go-kit/kit/endpoint"
)

type EndPoints struct {
	RegisterEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s UserService) EndPoints {
	return EndPoints{
		RegisterEndpoint: MakeRegisterEndPoint(s),
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

func MakeRegisterEndPoint(s UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(registerRequest)
		e := s.Register(ctx, req.User)
		return registerResponse{
			Token: req.User.RefreshToken,
			Err:   e,
		}, nil
	}
}

type registerRequest struct {
	User *model.User
}

type registerResponse struct {
	Token string `json:"token,omitempty"`
	Err   error  `json:"err,omitempty"`
}

func (e registerResponse) error() error { return e.Err }
