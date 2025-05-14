package service

import (
	"context"

	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/infra/grpc/pb"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	ListOrdersUseCase usecase.ListOrdersUseCase
}

func NewOrderService(ListOrdersUseCase usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{
		ListOrdersUseCase: ListOrdersUseCase,
	}
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	output, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}
	var orders []*pb.Order
	for _, order := range output {
		orders = append(orders, &pb.Order{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		})
	}
	return &pb.ListOrdersResponse{
		Orders: orders,
	}, nil
}
