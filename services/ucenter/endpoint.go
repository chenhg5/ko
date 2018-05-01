package ucenter

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

type getuserRequest struct {
	UID string `json:"s"`
}

type getuserResponse struct {
	Name   string `json:"v"`
	Err error `json:"err,omitempty"`
}

func makeGetUserEndpoint(s UcenterService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getuserRequest)
		name, err := s.GetUser(req.UID)
		return getuserResponse{Name: name, Err: err}, nil
	}
}