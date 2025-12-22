package service

import (
	"context"
	pb "kakao-shopping/proto"
	"kakao-shopping/server/store"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartService struct {
	pb.UnimplementedCartServiceServer
	cartStore    *store.CartStore
	productStore *store.ProductStore
}

func NewCartService(cartStore *store.CartStore, productStore *store.ProductStore) *CartService {
	return &CartService{
		cartStore:    cartStore,
		productStore: productStore,
	}
}

func (s *CartService) AddToCart(ctx context.Context, req *pb.AddToCartRequest) (*pb.CartResponse, error) {
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "사용자 ID가 필요합니다")
	}

	if req.Quantity <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "수량은 1개 이상이어야 합니다")
	}

	product, exists := s.productStore.GetProduct(req.ProductId)
	if !exists {
		return nil, status.Errorf(codes.NotFound, "상품을 찾을 수 없습니다: ID %d", req.ProductId)
	}

	if product.Stock < req.Quantity {
		return nil, status.Errorf(codes.FailedPrecondition, "재고가 부족합니다. 현재 재고: %d개", product.Stock)
	}

	cart := s.cartStore.AddToCart(req.UserId, product, req.Quantity)

	return &pb.CartResponse{
		Cart:    cart,
		Message: "삼품이 장바구니에 추가되었습니다",
	}, nil
}

func (s *CartService) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.CartResponse, error) {
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "사용자 ID가 필요합니다")
	}

	cart := s.cartStore.GetCart(req.UserId)

	return &pb.CartResponse{
		Cart:    cart,
		Message: "장바구니 조회 완료",
	}, nil
}

func (s *CartService) RemoveFromCart(ctx context.Context, req *pb.RemoveFromCartRequest) (*pb.CartResponse, error) {
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "사용자 ID가 필요합니다")
	}

	cart := s.cartStore.RemoveFromCart(req.UserId, req.ProductId)

	return &pb.CartResponse{
		Cart:    cart,
		Message: "삼품이 장바구니에서 제거되었습니다",
	}, nil
}

func (s *CartService) ClearCart(ctx context.Context, req *pb.ClearCartRequest) (*pb.CartResponse, error) {
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "사용자 ID가 필요합니다")
	}

	cart := s.cartStore.ClearCart(req.UserId)

	return &pb.CartResponse{
		Cart:    cart,
		Message: "장바구니가 비워졌습니다",
	}, nil
}
