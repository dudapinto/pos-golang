package service

import (
	"context"

	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/infra/grpc/pb"
	"github.com/dudapinto/pos-golang/20-CleanArch-Desafio-3/internal/usecase"
)

type ListOrdersService struct {
	pb.UnimplementedListOrdersServiceServer
	ListOrdersUseCase usecase.ListOrdersUseCase
}

func NewListOrdersService(ListOrdersUseCase usecase.ListOrdersUseCase) *ListOrdersService {
	return &ListOrdersService{
		ListOrdersUseCase: ListOrdersUseCase,
	}
}

func (s *ListOrdersService) ListOrders(ctx context.Context, in *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	output, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}
	return &pb.ListOrdersResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}
