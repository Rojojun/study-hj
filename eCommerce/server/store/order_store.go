package store

import (
	"fmt"
	pb "kakao-shopping/proto"
	"sync"
	"time"
)

type OrderStore struct {
	mu     sync.RWMutex
	orders map[string]*pb.Order
	nextID int32
}

func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: make(map[string]*pb.Order),
		nextID: 1,
	}
}

func (s *OrderStore) CreateOrder(userID string, items []*pb.CartItem, totalPrice int32) *pb.Order {
	s.mu.Lock()
	defer s.mu.Unlock()

	orderID := fmt.Sprintf("ORDER-%d", s.nextID)
	s.nextID++

	order := &pb.Order{
		OrderId:    orderID,
		UserId:     userID,
		Items:      items,
		TotalPrice: totalPrice,
		Status:     "pending",
		CreatedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

	s.orders[orderID] = order
	return order
}

func (s *OrderStore) GetOrder(orderID string) (*pb.Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.orders[orderID]
	return order, exists
}

func (s *OrderStore) ListUserOrders(userID string, page, pageSize int32) ([]*pb.Order, int32) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var userOrders []*pb.Order

	for _, order := range s.orders {
		if order.UserId == userID {
			userOrders = append(userOrders, order)
		}
	}

	totalPage := int32(len(userOrders))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= totalPage {
		return []*pb.Order{}, totalPage
	}

	if end > totalPage {
		end = totalPage
	}

	return userOrders[start:end], totalPage
}

func (s *OrderStore) UpdateOrderStatus(orderID string, status string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	order, exists := s.orders[orderID]
	if !exists {
		return false
	}

	order.Status = status
	return true
}
