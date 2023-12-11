package authsvc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeAuthHTTPHandler(c context.Context, s AuthService, log log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(c, s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST /authenticate authenticates a user and returns a token.
	r.Methods("POST").Path("/authenticate").Handler(httptransport.NewServer(
		e.AuthenticateEndpoint,
		decodeAuthenticateRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodeAuthenticateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		return nil, ErrNoAuthTokenHeader
	}
	ctx = context.WithValue(r.Context(), "jwt", jwtToken)
	//c := context.WithValue(r.Context(), "jwt", jwtToken)
	var req authenticateRequest
	req.Context = ctx
	return req, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(codeFrom(err))

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": err.Error(),
		"data":    nil,
	})
}

func codeFrom(err error) int {
	switch {
	case errors.Is(err, ErrAuthentication):
		return http.StatusUnauthorized // 401
	case errors.Is(err, ErrNoAuthTokenHeader):
		return http.StatusUnauthorized
	case errors.Is(err, ErrUnexpectedSigningMethod):
		return http.StatusUnauthorized
	case errors.Is(err, ErrInvalidToken):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError // 500
	}
}
