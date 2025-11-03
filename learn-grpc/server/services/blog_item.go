package services

import (
	"context"
	"fmt"
	"learn-grpc/pb"
	"learn-grpc/server/datalayer/actions"
	"learn-grpc/server/datalayer/models"
	"learn-grpc/server/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BlogItemService struct {
	pb.UnimplementedBlogServiceServer
	blogItemAction actions.BlogItemActionInterface
}

func NewBlogItemService(
	blogItemAction actions.BlogItemActionInterface,
) *BlogItemService {
	return &BlogItemService{
		blogItemAction: blogItemAction,
	}
}

func (s *BlogItemService) CreateBlog(ctx context.Context, in *pb.Blog) (*pb.BlogId, error) {
	fmt.Println("CreateBlog: new request")

	// Get authenticated user ID from context
	userID, ok := ctx.Value(middleware.UserIDKey).(int32)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// Use authenticated user as author (ignore in.AuthorId from request)
	data := models.BlogItem{
		AuthorId: uint(userID),
		Title:    in.Title,
		Content:  in.Content,
	}

	if err := s.blogItemAction.Save(ctx, &data); err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &pb.BlogId{Id: int32(data.ID)}, nil
}

func (s *BlogItemService) GetOneBlog(ctx context.Context, in *pb.BlogId) (*pb.Blog, error) {
	fmt.Println("GetOneBlog: new request")

	blogItem, err := s.blogItemAction.FetchById(ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return models.ToBlog(blogItem), nil
}

func (s *BlogItemService) UpdateBlog(ctx context.Context, in *pb.Blog) (*emptypb.Empty, error) {
	fmt.Println("UpdateBlog: new request")

	// Get authenticated user ID
	userID, ok := ctx.Value(middleware.UserIDKey).(int32)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	blogItem, err := s.blogItemAction.FetchById(ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	// Check if the user owns this blog
	if blogItem.AuthorId != uint(userID) {
		return nil, status.Errorf(codes.PermissionDenied, "you can only update your own blogs")
	}

	blogItem.Title = in.Title
	blogItem.Content = in.Content

	if err := s.blogItemAction.Save(ctx, blogItem); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *BlogItemService) DeleteBlog(ctx context.Context, in *pb.BlogId) (*emptypb.Empty, error) {
	fmt.Println("DeleteBlog: new request")

	// Get authenticated user ID
	userID, ok := ctx.Value(middleware.UserIDKey).(int32)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	blogItem, err := s.blogItemAction.FetchById(ctx, uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	// Check if the user owns this blog
	if blogItem.AuthorId != uint(userID) {
		return nil, status.Errorf(codes.PermissionDenied, "you can only delete your own blogs")
	}

	err = s.blogItemAction.SoftDelete(ctx, blogItem)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *BlogItemService) GetAllBlog(_ *emptypb.Empty, stream grpc.ServerStreamingServer[pb.Blog]) error {
	fmt.Println("GetAllBlog: new request")

	blogItems, err := s.blogItemAction.FetchList(context.Background())
	if err != nil {
		return status.Errorf(codes.Internal, "%v", err)
	}

	for _, bi := range blogItems {
		stream.Send(models.ToBlog(bi))
	}

	return nil
}
