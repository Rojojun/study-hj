package store

import (
	"strings"
	"sync"

	pb "kakao-shopping/proto"
)

type ProductStore struct {
	mu       sync.RWMutex
	products map[int32]*pb.Product
	nextID   int32
}

func NewProductStore() *ProductStore {
	return &ProductStore{
		products: make(map[int32]*pb.Product),
		nextID:   1,
	}
}

func (s *ProductStore) AddProduct(p *pb.Product) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if p.Id == 0 {
		p.Id = s.nextID
		s.nextID++
	}
	s.products[p.Id] = p
}

func (s *ProductStore) GetProduct(id int32) (*pb.Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product, exists := s.products[id]
	return product, exists
}

func (s *ProductStore) ListProducts(category string, page, pageSize int32) ([]*pb.Product, int32) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filteredProducts []*pb.Product

	for _, product := range s.products {
		if category == "" || product.Category == category {
			filteredProducts = append(filteredProducts, product)
		}
	}

	totalCount := int32(len(filteredProducts))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= totalCount {
		return []*pb.Product{}, totalCount
	}

	if end > totalCount {
		end = totalCount
	}

	return filteredProducts[start:end], totalCount
}

func (s *ProductStore) SearchProducts(query string, page, pageSize int32) ([]*pb.Product, int32) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var searchResults []*pb.Product
	query = strings.ToLower(query)

	for _, product := range s.products {
		if strings.Contains(strings.ToLower(product.Name), query) ||
			strings.Contains(strings.ToLower(product.Description), query) {
			searchResults = append(searchResults, product)
		}
	}

	totalCount := int32(len(searchResults))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if end >= totalCount {
		return []*pb.Product{}, totalCount
	}

	if end > totalCount {
		end = totalCount
	}

	return searchResults[start:end], totalCount
}

func (s *ProductStore) UpdateStock(productID int32, quantity int32) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	product, exists := s.products[productID]
	if !exists {
		return false
	}

	newStock := product.Stock + quantity
	if newStock < 0 {
		return false
	}

	product.Stock = newStock
	return true
}
