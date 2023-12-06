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
	// GET /stations lists all the stations.
	// PUT /station?id=<stationId> updates the station info.
	// DELETE /station?id=<stationId> deletes the station.
	// GEt /station/search?brand=<brandName> searches for a station by brand name.

	r.Methods("POST").Path("/station").Handler(httptransport.NewServer(
		e.StationRegisterEndpoint,
		decodeStationRegisterRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/stations").Handler(httptransport.NewServer(
		e.GetAllStationsEndpoint,
		decodeGetAllStationsRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/station").Handler(httptransport.NewServer(
		e.UpdateStationInfoEndpoint,
		decodeUpdateStationInfoRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/station").Handler(httptransport.NewServer(
		e.RemoveStationEndpoint,
		decodeRemoveStationRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/station/search").Handler(httptransport.NewServer(
		e.SearchStationEndpoint,
		decodeSearchStationRequest,
		encodeResponse,
		options...,
	))
	return r
}

type errorer interface {
	error() error
}

func decodeSearchStationRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		return nil, usersvc.ErrNoAuthTokenHeader
	}

	brand := r.URL.Query().Get("brand")

	ctx = context.WithValue(r.Context(), "jwt", jwtToken)
	c := context.WithValue(r.Context(), "jwt", jwtToken)
	var req searchStationRequest
	req.Context = c
	req.Brand = brand
	return req, nil
}

func decodeRemoveStationRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		return nil, usersvc.ErrNoAuthTokenHeader
	}

	stationId := r.URL.Query().Get("id")

	ctx = context.WithValue(r.Context(), "jwt", jwtToken)
	c := context.WithValue(r.Context(), "jwt", jwtToken)
	var req removeStationRequest
	req.Context = c
	req.StationID = stationId
	return req, nil
}

func decodeUpdateStationInfoRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		return nil, usersvc.ErrNoAuthTokenHeader
	}

	stationId := r.URL.Query().Get("id")

	ctx = context.WithValue(r.Context(), "jwt", jwtToken)
	c := context.WithValue(r.Context(), "jwt", jwtToken)
	var req updateStationInfoRequest
	req.Context = c
	req.StationID = stationId
	if err := json.NewDecoder(r.Body).Decode(&req.Station); err != nil {
		return nil, err
	}
	return req, nil
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

func decodeGetAllStationsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		return nil, usersvc.ErrNoAuthTokenHeader
	}

	ctx = context.WithValue(r.Context(), "jwt", jwtToken)
	c := context.WithValue(r.Context(), "jwt", jwtToken)
	var req getAllStationsRequest
	req.Context = c
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
