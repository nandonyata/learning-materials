package models

import (
	"learn-grpc/proto"

	"gorm.io/gorm"
)

type BlogItem struct {
	gorm.Model
	AuthorId uint
	Title    string
	Content  string
}

func ToBlog(b *BlogItem) *proto.Blog {
	return &proto.Blog{
		Id:       int32(b.ID),
		AuthorId: int32(b.AuthorId),
		Title:    b.Title,
		Content:  b.Content,
	}
}
