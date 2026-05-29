// Package response 提供统一的 HTTP JSON 响应封装。
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lagom/lagom-server/internal/pkg/errors"
)

// Response 统一响应结构。
type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// requestIDKey 用于从 gin.Context 读取请求 ID。
const requestIDKey = "request_id"

// extractRequestID 从 gin context 获取 request_id。
func extractRequestID(c *gin.Context) string {
	if rid, exists := c.Get(requestIDKey); exists {
		if s, ok := rid.(string); ok {
			return s
		}
	}
	return c.GetHeader("X-Request-ID")
}

// Success 返回 200 成功响应，携带数据。
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   "success",
		Data:      data,
		RequestID: extractRequestID(c),
	})
}

// OK 返回 200 成功响应，无数据。
func OK(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:      0,
		Message:   "success",
		RequestID: extractRequestID(c),
	})
}

// Created 返回 201 创建成功响应。
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Response{
		Code:      0,
		Message:   "created",
		Data:      data,
		RequestID: extractRequestID(c),
	})
}

// Page 返回分页数据响应。
func Page(c *gin.Context, list any, total int64) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: gin.H{
			"list":  list,
			"total": total,
		},
		RequestID: extractRequestID(c),
	})
}

// Fail 返回错误响应，自动解析 *errors.AppError。
// 如果不是 *AppError，则 fallback 到 500 内部错误。
func Fail(c *gin.Context, err error) {
	if err == nil {
		OK(c)
		return
	}

	if ae, ok := errors.IsAppError(err); ok {
		c.JSON(ae.Status, Response{
			Code:      ae.Code,
			Message:   ae.Message,
			RequestID: extractRequestID(c),
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Code:      errors.ErrInternal.Code,
		Message:   errors.ErrInternal.Message,
		RequestID: extractRequestID(c),
	})
}

// FailWithCode 返回指定业务码和消息的错误响应，HTTP 状态码 400。
func FailWithCode(c *gin.Context, code int, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:      code,
		Message:   message,
		RequestID: extractRequestID(c),
	})
}

// BadRequest 返回 400 参数错误。
func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrBadRequest.Message
	}
	c.JSON(http.StatusBadRequest, Response{
		Code:      errors.ErrBadRequest.Code,
		Message:   message,
		RequestID: extractRequestID(c),
	})
}

// Unauthorized 返回 401 未授权。
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrUnauthorized.Message
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
		Code:      errors.ErrUnauthorized.Code,
		Message:   message,
		RequestID: extractRequestID(c),
	})
}

// Forbidden 返回 403 禁止访问。
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrForbidden.Message
	}
	c.AbortWithStatusJSON(http.StatusForbidden, Response{
		Code:      errors.ErrForbidden.Code,
		Message:   message,
		RequestID: extractRequestID(c),
	})
}

// NotFound 返回 404 资源不存在。
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrNotFound.Message
	}
	c.AbortWithStatusJSON(http.StatusNotFound, Response{
		Code:      errors.ErrNotFound.Code,
		Message:   message,
		RequestID: extractRequestID(c),
	})
}

// Conflict 返回 409 资源冲突。
func Conflict(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrConflict.Message
	}
	c.AbortWithStatusJSON(http.StatusConflict, Response{
		Code:      errors.ErrConflict.Code,
		Message:   message,
		RequestID: extractRequestID(c),
	})
}

// InternalError 返回 500 内部错误。
func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrInternal.Message
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, Response{
		Code:      errors.ErrInternal.Code,
		Message:   message,
		RequestID: extractRequestID(c),
	})
}

// NotImplemented 返回 501 功能未实现。
func NotImplemented(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrNotImplemented.Message
	}
	c.JSON(http.StatusNotImplemented, Response{
		Code:      errors.ErrNotImplemented.Code,
		Message:   message,
		RequestID: extractRequestID(c),
	})
}

// TooManyRequests 返回 429 请求过于频繁。
func TooManyRequests(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrTooManyRequests.Message
	}
	c.AbortWithStatusJSON(http.StatusTooManyRequests, Response{
		Code:      errors.ErrTooManyRequests.Code,
		Message:   message,
		RequestID: extractRequestID(c),
	})
}
