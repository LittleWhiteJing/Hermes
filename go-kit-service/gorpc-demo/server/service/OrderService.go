package service

import (
	"context"
	prod "github.com/TyrellJing/Hermes/go-kit-service/gorpc/server/proto"
)

type OrderService struct {

}

func (o *OrderService) NewOrder(ctx context.Context, r *prod.OrderRequest) (*prod.OrderResponse, error) {
	err := r.OrderMain.Validate()
	if err != nil {
		return &prod.OrderResponse{
			Status: "err",
			Message: err.Error(),
		}, nil
	}
	return &prod.OrderResponse{
		Status: "OK",
		Message: "success",
	}, nil
}