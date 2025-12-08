package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pb "kakao-shopping/proto"
)

func main() {
	conn, err := grpc.("localhost:8080")
}
