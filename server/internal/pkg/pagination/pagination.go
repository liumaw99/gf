// Package pagination 提供统一的分页参数解析和分页结果封装。
package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	maxPageSize     = 100
)

// Param 分页请求参数。
type Param struct {
	Page     int `json:"page"`      // 当前页，从 1 开始
	PageSize int `json:"page_size"` // 每页条数
}

// Default 返回默认分页参数。
func Default() *Param {
	return &Param{
		Page:     defaultPage,
		PageSize: defaultPageSize,
	}
}

// FromContext 从 gin query 参数解析分页信息。
// 支持 query key：page / page_size（或 per_page）。
func FromContext(c *gin.Context) *Param {
	p := Default()

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			p.Page = page
		}
	}

	// 优先 page_size，兼容 per_page
	pageSizeStr := c.Query("page_size")
	if pageSizeStr == "" {
		pageSizeStr = c.Query("per_page")
	}
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			p.PageSize = ps
		}
	}

	// 限制最大页大小
	if p.PageSize > maxPageSize {
		p.PageSize = maxPageSize
	}

	return p
}

// Limit 返回 SQL LIMIT 值（即 PageSize）。
func (p *Param) Limit() int {
	return p.PageSize
}

// Offset 返回 SQL OFFSET 值。
func (p *Param) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// --- 分页结果 ---

// Result 分页查询结果（泛型）。
type Result[T any] struct {
	Data       []T   `json:"data"`        // 数据列表
	Total      int64 `json:"total"`       // 总记录数
	Page       int   `json:"page"`        // 当前页
	PageSize   int   `json:"page_size"`   // 每页条数
	TotalPages int   `json:"total_pages"` // 总页数
	HasNext    bool  `json:"has_next"`    // 是否有下一页
}

// NewResult 创建分页结果。
func NewResult[T any](data []T, total int64, param *Param) *Result[T] {
	totalPages := 0
	if param.PageSize > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(param.PageSize)))
	}

	return &Result[T]{
		Data:       data,
		Total:      total,
		Page:       param.Page,
		PageSize:   param.PageSize,
		TotalPages: totalPages,
		HasNext:    param.Page < totalPages,
	}
}
