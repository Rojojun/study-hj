package service

import (
	pb "kakao-shopping/proto"
	"kakao-shopping/server/store"
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
