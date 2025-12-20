package store

import (
	"sync"

	pb "kakao-shopping/proto"
)

type CartStore struct {
	mu    sync.RWMutex
	carts map[string]*pb.Cart
}

func NewCartStore() *CartStore {
	return &CartStore{
		carts: make(map[string]*pb.Cart),
	}
}

func (s *CartStore) GetCart(userID string) *pb.Cart {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cart, exists := s.carts[userID]
	if !exists {
		return &pb.Cart{
			UserId:     userID,
			Items:      []*pb.CartItem{},
			TotalPrice: 0,
		}
	}

	return cart
}

func (s *CartStore) AddToCart(userID string, product *pb.Product, quantity int32) *pb.Cart {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, exists := s.carts[userID]
	if !exists {
		cart = &pb.Cart{
			UserId:     userID,
			Items:      []*pb.CartItem{},
			TotalPrice: 0,
		}
	}

	for _, item := range cart.Items {
		if item.Product.Id == product.Id {
			item.Quantity += quantity
			s.updateCartTotal(cart)
			return cart
		}
	}

	cartItem := &pb.CartItem{
		Product:  product,
		Quantity: quantity,
	}
	cart.Items = append(cart.Items, cartItem)
	s.updateCartTotal(cart)

	return cart
}

func (s *CartStore) RemoveFromCart(userID string, productID int32) *pb.Cart {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, exists := s.carts[userID]
	if !exists {
		return &pb.Cart{
			UserId:     userID,
			Items:      []*pb.CartItem{},
			TotalPrice: 0,
		}
	}

	var newItems []*pb.CartItem
	for _, item := range cart.Items {
		if item.Product.Id != productID {
			newItems = append(newItems, item)
		}
	}

	cart.Items = newItems
	s.updateCartTotal(cart)

	return cart
}

func (s *CartStore) ClearCart(userID string) *pb.Cart {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart := &pb.Cart{
		UserId:     userID,
		Items:      []*pb.CartItem{},
		TotalPrice: 0,
	}
	s.carts[userID] = cart

	return cart
}

func (s *CartStore) updateCartTotal(cart *pb.Cart) {
	var total int32 = 0
	for _, item := range cart.Items {
		total += item.Product.Price * item.Quantity
	}
	cart.TotalPrice = total
}
