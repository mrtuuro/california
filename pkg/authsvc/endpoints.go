package authsvc

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	AuthenticateEndpoint endpoint.Endpoint
}

type BaseResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func MakeServerEndpoints(c context.Context, s AuthService) Endpoints {
	return Endpoints{
		AuthenticateEndpoint: MakeAuthenticateEndpoint(c, s),
	}
}

func MakeAuthenticateEndpoint(ctx context.Context, s AuthService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		req := request.(authenticateRequest)
		jwt := req.Context.Value("jwt")
		ctx = context.WithValue(ctx, "Authorization", jwt)
		e := s.Authenticate(ctx)
		if e != nil {
			return authenticateResponse{
				Err: e,
			}, e
		}
		return BaseResponse{
			Message: "success",
			Data: authenticateResponse{
				Token: jwt.(string),
				Err:   e,
			},
		}, nil
	}
}

type authenticateRequest struct {
	Context context.Context
	Token   string `json:"token"`
}

type authenticateResponse struct {
	*BaseResponse
	Token string `json:"token,omitempty"`
	Err   error  `json:"err,omitempty"`
}

func (e authenticateResponse) error() error { return e.Err }
