package main

import (
	"context"

	prod "github.com/TyrellJing/Hermes/go-kit-service/gorpc/server/proto"
)

type ProdService struct {

}

func (p *ProdService) GetProdStock(ctx context.Context, request *prod.ProdRequest) (*prod.ProdResponse, error) {
	return &prod.ProdResponse{ProdStock: 20}, nil;
}
