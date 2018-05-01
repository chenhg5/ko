package ucenter

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

type getuserRequest struct {
	UID string `json:"s"`
}

type getuserResponse struct {
	Data  map[string]interface{}   `json:"data"`
	Err   error                    `json:"err,omitempty"`
}

func makeGetUserEndpoint(s UcenterService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getuserRequest)
		name, err := s.GetUser(req.UID)
		data := map[string]interface{} {
			"user" : name,
		}
		return getuserResponse{Data: data, Err: err}, nil
	}
}