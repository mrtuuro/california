package navigationsvc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"california/pkg/usersvc"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHTTPHandler(c context.Context, s NavigationService, log log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(c, s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(log)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// GET /trip?distance=543 returns the trip information.
	r.Methods("GET").Path("/trip").Handler(httptransport.NewServer(
		e.CalculateTripEndpoint,
		decodeCalculateTripRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodeCalculateTripRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	authHeader := r.Header.Get("Authorization")
	jwtToken := strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" {
		return nil, usersvc.ErrNoAuthTokenHeader
	}

	distanceStr := r.URL.Query().Get("distance")
	distFloat, _ := strconv.ParseFloat(distanceStr, 64)

	ctx = context.WithValue(ctx, "jwt", jwtToken)
	c := context.WithValue(r.Context(), "jwt", jwtToken)
	c = context.WithValue(c, "distance", distanceStr)
	var req calculateTripRequest
	req.Context = c
	req.Distance = distFloat
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
