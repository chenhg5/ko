package ucenter

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

type getuserRequest struct {
	UID string `json:"s"`
}

type commonResponse struct {
	Code  int                      `json:"code"`
	Msg   string                   `json:"msg"`
	Data  map[string]interface{}   `json:"data"`
	Err   string                    `json:"err,omitempty"`
}

func makeGetUserEndpoint(s UcenterServiceInterface) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getuserRequest)
		name, err := s.GetUser(req.UID)
		data := map[string]interface{} {
			"user" : name,
		}
		var errmsg string
		if err != nil {
			errmsg = err.Error()
		} else {
			errmsg = ""
		}
		return commonResponse{Code: 0, Msg: "ok", Data: data, Err: errmsg}, nil
	}
}