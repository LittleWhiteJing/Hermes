package service

import (
	"context"
	prod "github.com/TyrellJing/Hermes/go-kit-service/gorpc/server/proto"
)

type ProdService struct {

}

func (p *ProdService) GetProdStock(ctx context.Context, request *prod.ProdRequest) (*prod.ProdResponse, error) {
	var stock int32 = 0
	if request.ProdArea == prod.ProdAreas_A {
		stock = 30
	} else if request.ProdArea == prod.ProdAreas_B {
		stock = 40
	} else {
		stock = 50
	}
	return &prod.ProdResponse{ProdStock: stock}, nil
}

func (p *ProdService) GetProdStocks(context.Context, *prod.ListRequest) (*prod.ListResponse, error) {
	prodRes := []*prod.ProdResponse {
		&prod.ProdResponse{ ProdStock: 31},
		&prod.ProdResponse{ ProdStock: 32},
		&prod.ProdResponse{ ProdStock: 33},
		&prod.ProdResponse{ ProdStock: 34},
	}
	return &prod.ListResponse{
		ProdRes: prodRes,
	}, nil
}

func (p * ProdService) GetProdInfo(ctx context.Context, in *prod.ProdRequest) (*prod.ProdModel, error) {
	ret := prod.ProdModel{
		ProdId: 101,
		ProdName: "测试商品",
		ProdPrice: 20.5,
	}
	return &ret, nil
}

