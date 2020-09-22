package main

import (
	"context"
	"fmt"
	"github.com/TyrellJing/Hermes/go-kit-service/grpc-server/pb"
	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
)

type BookServiceServer struct {
	BookListHandler kitgrpc.Handler
	BookInfoHandler kitgrpc.Handler
}

func (b *BookServiceServer) GetBookList(ctx context.Context, req *pb.BookListRequest) (*pb.BookListResponse, error) {
	_, resp, err := b.BookListHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.BookListResponse), nil
}

func (b *BookServiceServer) GetBookInfo(ctx context.Context, req *pb.BookDetailRequest) (*pb.BookDetailResponse, error) {
	_, resp, err := b.BookInfoHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.BookDetailResponse), nil
}

func BookListEndpoint () endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.BookListRequest)
		var bookInfo []*pb.BookInfo
		for i := 0; i < int(req.PerPage); i++ {
			bookInfo = append(bookInfo, &pb.BookInfo{
				BookId: int32(i),
				BookName: fmt.Sprintf("第%d版语文教材", i),
				BookPrice: 10.2,
				BoolClass: pb.BookClass_b_class_B,
			})
		}
		return &pb.BookListResponse{BookInfo: bookInfo}, nil
	}
}

func BookInfoEndPoint () endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.BookDetailRequest)
		var bookDetail pb.BookDetail
		switch req.BookId {
		case 1:
			bookDetail.BookId = 1001
			bookDetail.BookName = "zui"
			bookDetail.BookAuthor = "Guo"
			bookDetail.BookPage = 100
		case 2:
			bookDetail.BookId = 1002
			bookDetail.BookName = "han"
			bookDetail.BookAuthor = "Han"
			bookDetail.BookPage = 101
		default:
			bookDetail.BookId = 1000
			bookDetail.BookName = "history"
			bookDetail.BookAuthor = "Yu"
			bookDetail.BookPage = 100
		}
		return &pb.BookDetailResponse{BookDetail: &bookDetail}, nil
	}
}

func dencodeServiceRequest (ctx context.Context, req interface{}) (request interface{}, err error) {
	return req, nil
}

func encodeServiceResponse (ctx context.Context, resp interface{}) (response interface{}, err error) {
	return resp, nil
}

