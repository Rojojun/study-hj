package main

import (
	pb "kakao-shopping/proto"
	"kakao-shopping/server/service"
	"kakao-shopping/server/store"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	productStore := store.NewProductStore()
	cartStore := store.NewCartStore()
	orderStore := store.NewOrderStore()

	initSampleData(productStore)

	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     15 * time.Second,
			MaxConnectionAge:      30 * time.Second,
			MaxConnectionAgeGrace: 5 * time.Second,
			Time:                  5 * time.Second,
			Timeout:               1 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second,
			PermitWithoutStream: true,
		}),
	)

	pb.RegisterProductServiceServer(server, service.NewProductService(productStore))
	pb.RegisterCartServiceServer(server, service.NewCartService(cartStore, productStore))
	pb.RegisterOrderServiceServer(server, service.NewOrderService(orderStore, cartStore))
}

func initSampleData(productStore *store.ProductStore) {
	products := []*pb.Product{
		{
			Id:          1,
			Name:        "iPhone 15 Pro",
			Description: "최신 아이폰 15 프로 모델",
			Price:       1200000,
			Category:    "전자제품",
			Stock:       50,
			ImageUrl:    "https://example.com/iphone15pro.jpg",
		},
		{
			Id:          2,
			Name:        "MacBook Air M3",
			Description: "애플 맥북 에어 M3 칩셋",
			Price:       1500000,
			Category:    "전자제품",
			Stock:       30,
			ImageUrl:    "https://example.com/macbook-air-m3.jpg",
		},
		{
			Id:          3,
			Name:        "나이키 에어맥스",
			Description: "편안한 운동화",
			Price:       150000,
			Category:    "신발",
			Stock:       100,
			ImageUrl:    "https://example.com/nike-airmax.jpg",
		},
		{
			Id:          4,
			Name:        "삼성 갤럭시 S24",
			Description: "삼성 최신 스마트폰",
			Price:       1100000,
			Category:    "전자제품",
			Stock:       40,
			ImageUrl:    "https://example.com/galaxy-s24.jpg",
		},
		{
			Id:          5,
			Name:        "아디다스 운동복",
			Description: "편안한 트레이닝복",
			Price:       80000,
			Category:    "의류",
			Stock:       200,
			ImageUrl:    "https://example.com/adidas-tracksuit.jpg",
		},
	}

	for _, product := range products {
		productStore.AddProduct(product)
	}
	log.Printf("샘플 상품 %d개 추가 완료", len(products))
}
