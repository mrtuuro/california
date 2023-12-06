package charge_stationsvc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"california/pkg/usersvc"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeStationHTTPHandlers(c context.Context, s StationService, log log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(c, s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST /station adds a new station to the database.
	r.Methods("POST").Path("/station").Handler(httptransport.NewServer(
		e.StationRegisterEndpoint,
		decodeStationRegisterRequest,
		encodeResponse,
		options...,
	))
	return r
}

type errorer interface {
	error() error
}

func decodeStationRegisterRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		return nil, usersvc.ErrNoAuthTokenHeader
	}

	ctx = context.WithValue(r.Context(), "jwt", jwtToken)
	c := context.WithValue(r.Context(), "jwt", jwtToken)
	var req registerStationRequest
	req.Context = c
	if err := json.NewDecoder(r.Body).Decode(&req.Station); err != nil {
		return nil, err
	}
	return req, nil
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
	case errors.Is(err, usersvc.ErrNotFound):
		return http.StatusNotFound // 404
	case errors.Is(err, usersvc.ErrAlreadyExists), errors.Is(err, usersvc.ErrInconsistentIDs):
		return http.StatusBadRequest // 400
	case errors.Is(err, usersvc.ErrAuthentication):
		return http.StatusUnauthorized // 401
	case errors.Is(err, usersvc.ErrPasswordEmailDoesNotMatch):
		return http.StatusUnauthorized // 401
	case errors.Is(err, usersvc.ErrNoAuthTokenHeader):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError // 500
	}
}
