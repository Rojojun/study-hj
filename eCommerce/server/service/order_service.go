package service

import (
	"context"
	pb "kakao-shopping/proto"
	"kakao-shopping/server/store"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	pb.UnimplementedCartServiceServer
	orderStore *store.OrderStore
	cartStore  *store.CartStore
}

func NewOrderService(orderStore *store.OrderStore, cartStore *store.CartStore) *OrderService {
	return &OrderService{
		orderStore: orderStore,
		cartStore:  cartStore,
	}
}

func (os *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "사용자 ID가 필요합니다")
	}

	cart := os.cartStore.GetCart(req.UserId)
	if len(cart.Items) == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "장바구니가 비어있습니다")
	}

	order := os.orderStore.CreateOrder(req.UserId, cart.Items, cart.TotalPrice)

	os.cartStore.ClearCart(req.UserId)

	return order, nil
}

func (os *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	if req.OrderId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "주문 ID가 필요합니다")
	}

	order, exists := os.orderStore.GetOrder(req.OrderId)
	if !exists {
		return nil, status.Errorf(codes.NotFound, "주문을 찾을 수 없습니다: %s", req.OrderId)
	}

	return order, nil
}

func (os *OrderService) ListUserOrders(ctx context.Context, req *pb.ListUserOrdersRequest) (*pb.ListOrdersResponse, error) {
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "사용자 ID가 필요합니다")
	}

	orders, totalCount := os.orderStore.ListUserOrders(req.UserId, req.PageSize, req.Page)

	return &pb.ListOrdersResponse{
		Orders:     orders,
		TotalCount: totalCount,
	}, nil
}
