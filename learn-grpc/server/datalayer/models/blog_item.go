package models

import (
	"learn-grpc/pb"

	"gorm.io/gorm"
)

type BlogItem struct {
	gorm.Model
	AuthorId uint
	Title    string
	Content  string
}

func ToBlog(b *BlogItem) *pb.Blog {
	return &pb.Blog{
		Id:       int32(b.ID),
		AuthorId: int32(b.AuthorId),
		Title:    b.Title,
		Content:  b.Content,
	}
}
