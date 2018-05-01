package order

import (
	"strings"
	"errors"
)

type OrderServiceInterface interface {
	GetUser(string) (string, error)
}

type OrderService struct{}

func (OrderService) GetOrder(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

var ErrEmpty = errors.New("empty input")

type ServiceMiddleware func(OrderServiceInterface) OrderServiceInterface