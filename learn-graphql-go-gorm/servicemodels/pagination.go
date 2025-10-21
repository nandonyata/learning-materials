package servicemodels

import (
	"gorm.io/gorm"
)

type PaginationResponse[T any] struct {
	Data      []T   `json:"data"`
	TotalPage uint  `json:"total_page"`
	TotalData *uint `json:"total_data"`
	NextPage  bool  `json:"next_page"`
}

type PaginationRequest struct {
	PageNo    int   `form:"page_no" json:"page_no"`
	PerPage   int   `form:"per_page" json:"per_per"`
	TotalData int64 `form:"-" json:"-"`
}

func (p PaginationRequest) IsPaginationApplicable() bool {
	return p.PageNo >= 1 && p.PerPage >= 1
}

func (p PaginationRequest) GetLimit() int {
	return p.PerPage
}

func (p PaginationRequest) GetOffset() int {
	return (p.PageNo - 1) * p.PerPage
}

func WithPagination(data PaginationRequest) func(db *gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		if data.IsPaginationApplicable() {
			query = query.Offset(data.GetOffset()).Limit(data.GetLimit())
		}
		return query
	}
}
