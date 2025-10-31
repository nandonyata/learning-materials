package main

import (
	"learn-grpc/pb"
	database "learn-grpc/server/datalayer"
	"learn-grpc/server/datalayer/actions"
	"learn-grpc/server/services"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	pb.BlogServiceServer
}

var (
	address string = "localhost:3002"
)

func main() {
	database.InitDB()
	defer database.CloseDB()
	database.Migrate()
	db := database.DB

	// Actions
	blogItemAction := actions.NewBlogItemAction(db)

	// Services
	blogItemService := services.NewBlogItemService(blogItemAction)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed listening address: %v\n", err)
	}

	log.Printf("Listening on address: %v\n", address)

	server := grpc.NewServer()
	pb.RegisterBlogServiceServer(server, blogItemService)

	if err = server.Serve(listener); err != nil {
		log.Fatalf("Failed serving: %v\n", err)
	}

}
