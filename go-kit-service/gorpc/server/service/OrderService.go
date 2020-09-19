package service

import (
	"context"
	prod "github.com/TyrellJing/Hermes/go-kit-service/gorpc/server/proto"
)

type OrderService struct {

}

func (o *OrderService) NewOrder(context.Context, *prod.OrderMain) (*prod.OrderResponse, error) {
	return &prod.OrderResponse{
		Status: "OK",
		Message: "success",
	}, nil
}