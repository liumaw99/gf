package pagination

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	p := Default()
	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 20, p.PageSize)
}

func TestFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		query    string
		wantPage int
		wantSize int
	}{
		{"default", "", 1, 20},
		{"normal", "?page=2&page_size=10", 2, 10},
		{"per_page_compat", "?page=3&per_page=15", 3, 15},
		{"zero_page_fallback", "?page=0&page_size=5", 1, 5},
		{"negative_fallback", "?page=-1&page_size=-5", 1, 20},
		{"max_limit", "?page=1&page_size=999", 1, 100},
		{"invalid_str", "?page=abc&page_size=xyz", 1, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/"+tt.query, nil)

			p := FromContext(c)
			assert.Equal(t, tt.wantPage, p.Page)
			assert.Equal(t, tt.wantSize, p.PageSize)
		})
	}
}

func TestParam_LimitOffset(t *testing.T) {
	p := &Param{Page: 3, PageSize: 10}
	assert.Equal(t, 10, p.Limit())
	assert.Equal(t, 20, p.Offset())
}

func TestNewResult(t *testing.T) {
	type Item struct {
		ID int `json:"id"`
	}

	data := []Item{{ID: 1}, {ID: 2}, {ID: 3}}
	param := &Param{Page: 1, PageSize: 10}
	result := NewResult(data, int64(25), param)

	assert.Equal(t, 3, len(result.Data))
	assert.Equal(t, int64(25), result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)
	assert.Equal(t, 3, result.TotalPages)
	assert.True(t, result.HasNext)

	// 最后一页
	param2 := &Param{Page: 3, PageSize: 10}
	result2 := NewResult(data, int64(25), param2)
	assert.False(t, result2.HasNext)

	// 空结果
	param3 := &Param{Page: 1, PageSize: 10}
	result3 := NewResult([]Item{}, int64(0), param3)
	assert.Equal(t, 0, result3.TotalPages)
	assert.False(t, result3.HasNext)
}
