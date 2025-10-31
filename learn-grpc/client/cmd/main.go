package main

import (
	"learn-grpc/client/services"
	"learn-grpc/pb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var address string = "localhost:3002"

func main() {
	c, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed conecting address: %v\n", err)
	}

	defer c.Close()
	client := pb.NewBlogServiceClient(c)

	blogItemClient := services.NewBlogItemClientService(client)

	blogItemClient.CreateBlog()
	// blogItemClient.GetOneBlog(1)
	// blogItemClient.GetAllBlog()
	// blogItemClient.UpdateBlog(&pb.Blog{
	// 	Id:       1,
	// 	AuthorId: 10,
	// 	Title:    "Updated Title",
	// 	Content:  "Updated Content",
	// })
	// blogItemClient.DeleteBlog(1)

}
