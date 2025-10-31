package main

import (
	"learn-grpc/pb"
	database "learn-grpc/server/datalayer"
	"learn-grpc/server/datalayer/actions"
	"learn-grpc/server/middleware"
	"learn-grpc/server/services"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	address   string = "localhost:3002"
	secretKey string = "your-secret-key-change-in-production"
)

func main() {
	database.InitDB()
	defer database.CloseDB()
	database.Migrate()
	db := database.DB

	// Actions
	blogItemAction := actions.NewBlogItemAction(db)
	userAction := actions.NewUserAction(db)

	// Services
	blogItemService := services.NewBlogItemService(blogItemAction)
	userService := services.NewUserService(userAction)

	// Create auth interceptor
	authInterceptor := middleware.NewAuthInterceptor(secretKey)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed listening address: %v\n", err)
	}

	log.Printf("Listening on address: %v\n", address)

	// Create gRPC server with interceptors
	server := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.UnaryInterceptor),
		grpc.StreamInterceptor(authInterceptor.StreamInterceptor),
	)

	// Register services
	pb.RegisterBlogServiceServer(server, blogItemService)
	pb.RegisterUserServiceServer(server, userService)

	if err = server.Serve(listener); err != nil {
		log.Fatalf("Failed serving: %v\n", err)
	}
}
