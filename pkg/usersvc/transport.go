package usersvc

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

func MakeHTTPHandler(c context.Context, s UserService, log log.Logger, signingKey string) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(c, s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST /register/ adds a new user to the database.
	// POST /login/ logs in a user and returns a token.
	// POST /vehicle/register/ adds a new vehicle to the database.
	// GET /me/ returns the user's information.
	// PUT /user/ updates the user's information.

	r.Methods("POST").Path("/register/").Handler(httptransport.NewServer(
		e.RegisterEndpoint,
		decodeRegisterRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/login/").Handler(httptransport.NewServer(
		e.LoginEndpoint,
		decodeLoginRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/vehicle/register/").Handler(httptransport.NewServer(
		e.VehicleRegisterEndpoint,
		decodeVehicleRegisterRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/me/").Handler(httptransport.NewServer(
		e.GetMeEndpoint,
		decodeGetMeRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/user/").Handler(httptransport.NewServer(
		e.UpdateUserEndpoint,
		decodeUpdateUserRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodeGetMeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		return nil, ErrNoAuthToken
	}
	ctx = context.WithValue(r.Context(), "jwt", jwtToken)
	c := context.WithValue(r.Context(), "jwt", jwtToken)
	var req getMeRequest
	req.Context = c
	return req, nil
}

func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req registerRequest
	if e := json.NewDecoder(r.Body).Decode(&req.User); e != nil {
		return nil, e
	}
	return req, nil

}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req loginRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeVehicleRegisterRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		return nil, ErrNoAuthToken
	}
	ctx = context.WithValue(ctx, "jwt", jwtToken)
	c := context.WithValue(r.Context(), "jwt", jwtToken)

	var req vehicleRegisterRequest
	req.Context = c

	if e := json.NewDecoder(r.Body).Decode(&req.Vehicle); e != nil {
		return nil, e
	}
	return req, nil

}

func decodeUpdateUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		return nil, ErrNoAuthToken
	}
	ctx = context.WithValue(ctx, "jwt", jwtToken)
	c := context.WithValue(r.Context(), "jwt", jwtToken)

	var req updateUserRequest
	req.Context = c

	if e := json.NewDecoder(r.Body).Decode(&req.User); e != nil {
		return nil, e
	}
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
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound // 404
	case errors.Is(err, ErrAlreadyExists), errors.Is(err, ErrInconsistentIDs):
		return http.StatusBadRequest // 400
	case errors.Is(err, ErrAuthentication):
		return http.StatusUnauthorized // 401
	case errors.Is(err, ErrPasswordEmailDoesNotMatch):
		return http.StatusUnauthorized // 401
	default:
		return http.StatusInternalServerError // 500
	}
}
