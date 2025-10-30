package services

import (
	"context"
	"fmt"
	"learn-grpc/proto"
	"learn-grpc/server/datalayer/actions"
	"learn-grpc/server/datalayer/models"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BlogItemService struct {
	proto.BlogServiceServer
	blogItemAction actions.BlogItemActionInterface
}

func NewBlogItemService(
	blogItemAction actions.BlogItemActionInterface,
) *BlogItemService {
	return &BlogItemService{
		blogItemAction: blogItemAction,
	}
}

func (s *BlogItemService) CreateBlog(ctx context.Context, in *proto.Blog) (*proto.BlogId, error) {
	fmt.Println("CreateBlog: new request")

	data := models.BlogItem{
		AuthorId: uint(in.AuthorId),
		Title:    in.Title,
		Content:  in.Content,
	}

	if err := s.blogItemAction.Save(ctx, &data); err != nil {
		return nil, err
	}

	fmt.Printf("CreateBlog: new blog created with id %d\n", data.ID)

	return &proto.BlogId{Id: int32(data.ID)}, nil
}
func (s *BlogItemService) GetOneBlog(ctx context.Context, in *proto.BlogId) (*proto.Blog, error) {
	return nil, nil
}
func (s *BlogItemService) UpdateBlog(ctx context.Context, in *proto.Blog) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *BlogItemService) DeleteBlog(ctx context.Context, in *proto.BlogId) (*emptypb.Empty, error) {
	return nil, nil
}
func (s *BlogItemService) GetAllBlog(_ *emptypb.Empty, stream grpc.ServerStreamingServer[proto.Blog]) error {
	return nil
}
