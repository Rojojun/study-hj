package main

import (
	"context"
	pb "kakao-shopping/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func main() {
	conn, err := grpc.NewClient("localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                60 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: false,
		}),
	)

	if err != nil {
		log.Fatalf("서버 연결 실패: %v", err)
	}
	defer conn.Close()

	productClient := pb.NewProductServiceClient(conn)
	cartClient := pb.NewCartServiceClient(conn)
	orderClient := pb.NewOrderServiceClient(conn)

	ctx := context.Background()
	userID := "test_user"

	log.Println("카카오 쇼핑 클라이언트 시작")

	listProductsTest(productClient, ctx)
	searchProductTest(productClient, ctx)
	addProductTest(cartClient, ctx, userID)
	getCartTest(cartClient, ctx, userID)
	order := createOrderTest(orderClient, ctx, userID)
	getOrderTest(orderClient, ctx, order)
	listUserOrdersTest(orderClient, ctx, userID)
}

func listProductsTest(productClient pb.ProductServiceClient, ctx context.Context) {
	log.Println("\n1. 상품 목록 조회")
	listResp, err := productClient.ListProducts(ctx, &pb.ListProductsRequest{
		Page:     1,
		PageSize: 5,
	})
	if err != nil {
		log.Fatalf("상품 목록 조회 실패: %v", err)
	}

	log.Printf("총 %d개 상품 중 %d개 조회", listResp.TotalCount, len(listResp.Products))
	for _, product := range listResp.Products {
		log.Printf("- [%d] %s: %d원 (재고: %d개)",
			product.Id, product.Name, product.Price, product.Stock)
	}
}

func searchProductTest(productClient pb.ProductServiceClient, ctx context.Context) {
	log.Printf("\n2. 상품 검색 (아이폰)")
	searchResp, err := productClient.SearchProducts(ctx, &pb.SearchProductsRequest{
		Query:    "아이폰",
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Fatalf("상품 검색 실패: %v", err)
	}
	log.Printf("'아이폰' 검색 결과: %d개", len(searchResp.Products))
	for _, product := range searchResp.Products {
		log.Printf("- %s: %d원", product.Name, product.Price)
	}
}

func addProductTest(cartClient pb.CartServiceClient, ctx context.Context, userID string) {
	log.Printf("\n3. 장바구니에 상품 추가")

	addResp1, err := cartClient.AddToCart(ctx, &pb.AddToCartRequest{
		UserId:    userID,
		ProductId: 1,
		Quantity:  1,
	})
	if err != nil {
		log.Fatalf("장바구니 추가 실패: %v", err)
	}

	log.Printf("성공: %s", addResp1.Message)

	addResp2, err := cartClient.AddToCart(ctx, &pb.AddToCartRequest{
		UserId:    userID,
		ProductId: 3,
		Quantity:  2,
	})
	if err != nil {
		log.Fatalf("장바구니 추가 실패: %v", err)
	}

	log.Printf("성공: %s", addResp2.Message)
}

func getCartTest(cartClient pb.CartServiceClient, ctx context.Context, userID string) {
	log.Printf("\n4. 장바구니 조회")
	cartResp, err := cartClient.GetCart(ctx, &pb.GetCartRequest{
		UserId: userID,
	})
	if err != nil {
		log.Fatalf("장바구니 조회 실패: %v", err)
	}

	log.Printf("장바구니 상품 수: %d개", len(cartResp.Cart.Items))
	log.Printf("총 금액: %d원", cartResp.Cart.TotalPrice)
	for _, item := range cartResp.Cart.Items {
		log.Printf("- %s x%d = %d원",
			item.Product.Name, item.Quantity, item.Product.Price*item.Quantity)
	}
}

func createOrderTest(orderClient pb.OrderServiceClient, ctx context.Context, userID string) *pb.Order {
	order, err := orderClient.CreateOrder(ctx, &pb.CreateOrderRequest{
		UserId: userID,
	})
	if err != nil {
		log.Fatalf("주문 생성 실패: %v", err)
	}

	log.Printf("주문 생성 완료!")
	log.Printf("주문 ID: %s", order.OrderId)
	log.Printf("주문 상태: %s", order.Status)
	log.Printf("주문 금액: %d원", order.TotalPrice)
	log.Printf("주문 시간: %s", order.CreatedAt)

	return order
}

func getOrderTest(orderClient pb.OrderServiceClient, ctx context.Context, order *pb.Order) {
	getOrderResp, err := orderClient.GetOrder(ctx, &pb.GetOrderRequest{
		OrderId: order.OrderId,
	})
	if err != nil {
		log.Fatalf("주문 조회 실패: %v", err)
	}

	log.Printf("주문 상세 정보:")
	log.Printf("- 주문 ID: %s", getOrderResp.OrderId)
	log.Printf("- 상태: %s", getOrderResp.Status)
	log.Printf("- 상품 수: %d개", len(getOrderResp.Items))
}

func listUserOrdersTest(orderClient pb.OrderServiceClient, ctx context.Context, userID string) {
	ordersResp, err := orderClient.ListUserOrders(ctx, &pb.ListUserOrdersRequest{
		UserId:   userID,
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Fatalf("주문 목록 조회 실패: %v", err)
	}

	log.Printf("총 주문 수: %d개", ordersResp.TotalCount)
	for _, userOrder := range ordersResp.Orders {
		log.Printf("- %s (%s): %d원",
			userOrder.OrderId, userOrder.Status, userOrder.TotalPrice)
	}

	log.Println("\n모든 테스트 완료!")
}
