package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/sigame/auth/proto"
)

func main() {
	// Подключаемся к gRPC серверу
	conn, err := grpc.Dial("localhost:50051", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	// Тест 1: ValidateToken с неправильным токеном
	fmt.Println("=== Test 1: ValidateToken ===")
	resp1, err := client.ValidateToken(context.Background(), &pb.ValidateTokenRequest{
		Token: "invalid-token",
	})
	if err != nil {
		log.Fatalf("ValidateToken failed: %v", err)
	}
	fmt.Printf("Valid: %v, Error: %s\n", resp1.Valid, resp1.Error)

	// Тест 2: GetUserInfo с неправильным ID
	fmt.Println("\n=== Test 2: GetUserInfo ===")
	resp2, err := client.GetUserInfo(context.Background(), &pb.GetUserInfoRequest{
		UserId: "00000000-0000-0000-0000-000000000000",
	})
	if err != nil {
		log.Fatalf("GetUserInfo failed: %v", err)
	}
	fmt.Printf("Error: %s\n", resp2.Error)

	fmt.Println("\n✓ All gRPC tests passed!")
}

