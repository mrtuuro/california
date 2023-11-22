package usersvc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHTTPHandler(s UserService, log log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST /register/ adds a new user to the database.
	// POST /login/ logs in a user and returns a token.

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

	return r
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
	switch err {
	case ErrNotFound:
		return http.StatusNotFound // 404
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest // 400
	case ErrAuthentication:
		return http.StatusUnauthorized // 401
	case ErrPasswordEmailDoesNotMatch:
		return http.StatusUnauthorized // 401
	default:
		return http.StatusInternalServerError // 500
	}
}
