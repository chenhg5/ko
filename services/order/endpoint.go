package order

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

type getorderRequest struct {
	Param map[string]interface{} `json:"param"`
}

type getorderResponse struct {
	Data  map[string]interface{}   `json:"data"`
	Err   error                    `json:"err,omitempty"`
}

func makeGetOrderEndpoint(s OrderService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getorderRequest)
		name, err := s.GetOrder(req.Param["orderId"].(string))
		data := map[string]interface{} {
			"order" : name,
		}
		return getorderResponse{Data: data, Err: err}, nil
	}
}