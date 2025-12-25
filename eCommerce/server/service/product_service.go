package service

import (
	"context"
	pb "kakao-shopping/proto"
	"kakao-shopping/server/store"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductService struct {
	pb.UnimplementedProductServiceServer
	store *store.ProductStore
}

func NewProductService(store *store.ProductStore) *ProductService {
	return &ProductService{
		store: store,
	}
}

func (s *ProductService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, totalCount := s.store.ListProducts(req.Category, req.Page, req.PageSize)

	return &pb.ListProductsResponse{
		Products:   products,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}

func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	product, exists := s.store.GetProduct(req.ProductId)
	if !exists {
		return nil, status.Errorf(codes.NotFound, "상품을 찾을 수 없습니다: ID %d", req.ProductId)
	}

	return product, nil
}

func (s *ProductService) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.ListProductsResponse, error) {
	if req.Query == "" {
		return nil, status.Errorf(codes.InvalidArgument, "검색어를 입력해주세요")
	}

	products, totalCount := s.store.SearchProducts(req.Query, req.Page, req.PageSize)

	return &pb.ListProductsResponse{
		Products:   products,
		TotalCount: totalCount,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}
