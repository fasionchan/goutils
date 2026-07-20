package types

import (
	"github.com/fasionchan/goutils/stl"
)

const (
	DefaultPageSize = 20
)

type Pagination struct {
	PageSize int64 `json:"PageSize"`
	Pages    int64 `json:"Pages"`
	Items    int64 `json:"Items"`
	PageNo   int64 `json:"PageNo"`
}

func NewPagination(pageSize, pageNo int64) *Pagination {
	return &Pagination{
		PageSize: pageSize,
		PageNo:   pageNo,
	}
}

func NewSinglePagination(n int64) *Pagination {
	return NewPagination(n, 1).WithItems(n)
}

func NewPaginationFromTotalNext(total, next, pageSize int64) *Pagination {
	return &Pagination{
		PageSize: pageSize,
		Pages:    (total + pageSize - 1) / pageSize,
		Items:    total,
		PageNo:   next / pageSize,
	}
}

func (pagination *Pagination) Dup() *Pagination {
	return stl.Dup(pagination)
}

func (pagination *Pagination) WithDefault() *Pagination {
	if pagination == nil {
		return &Pagination{
			PageSize: DefaultPageSize,
			PageNo:   1,
		}
	}

	if pagination.PageSize == 0 {
		pagination.PageSize = DefaultPageSize
	}
	if pagination.PageNo == 0 {
		pagination.PageNo = 1
	}

	return pagination
}

func (pagination *Pagination) WithItems(items int64) *Pagination {
	pagination = pagination.WithDefault()
	pagination.Pages = (items + pagination.PageSize - 1) / pagination.PageSize
	pagination.Items = items

	if pagination.PageNo > pagination.Pages {
		pagination.PageNo = pagination.Pages
	}

	return pagination
}

func (pagination *Pagination) GetSkip() int64 {
	return (pagination.PageNo - 1) * pagination.PageSize
}

func (pagination *Pagination) GetLimit() int64 {
	return pagination.PageSize
}

func (pagination *Pagination) ItemsUpToCurrentPage() int64 {
	return pagination.PageSize * pagination.PageNo
}
