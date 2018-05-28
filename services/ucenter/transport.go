package ucenter

import (
	"context"
	"net/http"
	"encoding/json"
	kithttp "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"errors"
	"fmt"
)

func decodeGetuserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["UID"]
	if !ok {
		return nil, errBadRoute
	}
	fmt.Println("request: ", id)
	return getuserRequest{UID: id}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// MakeHandler returns a handler for the booking service.
func MakeHandler(bs UcenterServiceInterface, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	getUserHandler := kithttp.NewServer(
		makeGetUserEndpoint(bs),
		decodeGetuserRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	// 接口路由
	r.Handle("/svc/ucenter/v1/user/{UID}", getUserHandler).Methods("GET")

	return r
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//switch err {
	//	case ErrUnknown:
	//		w.WriteHeader(http.StatusNotFound)
	//	case ErrInvalidArgument:
	//		w.WriteHeader(http.StatusBadRequest)
	//	default:
	//		w.WriteHeader(http.StatusInternalServerError)
	//}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code": http.StatusNotFound,
		"msg": "from ucenter error: " + err.Error(),
	})
}

//var ErrUnknown = errors.New("unknown cargo")
//var ErrInvalidArgument = errors.New("error argument")
var errBadRoute = errors.New("error argument")