package service

import (
	pb "kakao-shopping/proto"
	"kakao-shopping/server/store"
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
