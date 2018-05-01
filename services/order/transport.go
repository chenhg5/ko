package order

import (
	"context"
	"net/http"
	"encoding/json"
	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"errors"
)

func decodeGetorderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getorderRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil

}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// MakeHandler returns a handler for the booking service.
func MakeHandler(bs OrderService, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getOrderHandler := kithttp.NewServer(
		makeGetOrderEndpoint(bs),
		decodeGetorderRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/svc/order/v1/order", getOrderHandler).Methods("POST")

	return r
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
		case ErrUnknown:
			w.WriteHeader(http.StatusNotFound)
		case ErrInvalidArgument:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

var ErrUnknown = errors.New("unknown cargo")
var ErrInvalidArgument = errors.New("error argument")
