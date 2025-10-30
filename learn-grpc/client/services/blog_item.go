package services

import (
	"context"
	"fmt"
	"io"
	"learn-grpc/proto"
	"log"
	"time"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BlogItemClientService struct {
	client proto.BlogServiceClient
}

func NewBlogItemClientService(client proto.BlogServiceClient) *BlogItemClientService {
	return &BlogItemClientService{
		client: client,
	}
}

func (s *BlogItemClientService) CreateBlog() {
	res, err := s.client.CreateBlog(context.Background(), &proto.Blog{
		AuthorId: 1,
		Title:    fmt.Sprintf("Random Title %d", time.Now().Unix()),
		Content:  fmt.Sprintf("Random Contnet %d", time.Now().Unix()),
	})
	if err != nil {
		_, ok := status.FromError(err)
		if ok {
			log.Printf("An error occured from server: %v", err)
			return
		} else {
			log.Printf("A non gRPC err: %v", err)
			return
		}
	}

	log.Printf("Success create blog, result id: %v\n", res)
}

func (s *BlogItemClientService) DeleteBlog(id int32) {
	_, err := s.client.DeleteBlog(context.Background(), &proto.BlogId{Id: id})
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	log.Println("Success delete blog")
}

func (s *BlogItemClientService) GetAllBlog() {
	res, err := s.client.GetAllBlog(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Err from server: %v\n", err)
	}

	for {
		stream, err := res.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Err receiving server: %v\n", err)
		}

		log.Printf("Success: %v", stream)
	}
}

func (s *BlogItemClientService) GetOneBlog(id int32) {
	log.Print("GetOneBlog was invoked ,,,")

	blog, err := s.client.GetOneBlog(context.Background(), &proto.BlogId{Id: id})
	if err != nil {
		_, ok := status.FromError(err)

		if ok {
			log.Printf("Error from server: %v\n", err)
		} else {
			log.Printf("Non gRPC err: %v", err)
		}
	}

	log.Printf("Success get one: %+v\n", blog)
}

func (s *BlogItemClientService) UpdateBlog(blogItem *proto.Blog) {
	log.Print("udpate was invoked bro..")

	res, err := s.client.UpdateBlog(context.Background(), blogItem)
	if err != nil {
		_, ok := status.FromError(err)
		if ok {
			log.Printf("Error from server %v\n", err)
		} else {
			log.Printf("Non grpc error %v", err)
		}
	}

	log.Printf("Success update one: %+v\n", res)
}
