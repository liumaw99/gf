// Package response 提供统一的 HTTP JSON 响应封装。
//
// 设计原则：所有响应（包括错误）统一返回 HTTP 200，业务状态通过 body 中的 code 区分。
// 业务码规则：
//   - 0        = 成功
//   - 1001-1099 = 通用错误
//   - 各模块自行扩展
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

// write 统一写入响应（HTTP 200 + JSON body）。
func write(c *gin.Context, code int, message string, data any) {
	c.JSON(http.StatusOK, Response{
		Code:      code,
		Message:   message,
		Data:      data,
		RequestID: extractRequestID(c),
	})
}

// --- 成功响应 ---

// Success 返回成功响应，携带数据。
func Success(c *gin.Context, data any) {
	write(c, 0, "success", data)
}

// OK 返回成功响应，无数据。
func OK(c *gin.Context) {
	write(c, 0, "success", nil)
}

// Created 返回创建成功响应。
func Created(c *gin.Context, data any) {
	write(c, 0, "created", data)
}

// Page 返回分页数据响应。
func Page(c *gin.Context, list any, total int64) {
	write(c, 0, "success", gin.H{
		"list":  list,
		"total": total,
	})
}

// --- 错误响应 ---

// Fail 返回错误响应，自动解析 *errors.AppError。
// 如果不是 *AppError，则 fallback 到内部错误。
func Fail(c *gin.Context, err error) {
	if err == nil {
		OK(c)
		return
	}

	if ae, ok := errors.IsAppError(err); ok {
		write(c, ae.Code, ae.Message, nil)
		return
	}

	write(c, errors.ErrInternal.Code, errors.ErrInternal.Message, nil)
}

// FailWithCode 返回指定业务码和消息的错误响应。
func FailWithCode(c *gin.Context, code int, message string) {
	write(c, code, message, nil)
}

// BadRequest 返回参数错误。
func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrBadRequest.Message
	}
	write(c, errors.ErrBadRequest.Code, message, nil)
}

// Unauthorized 返回未授权错误。
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrUnauthorized.Message
	}
	write(c, errors.ErrUnauthorized.Code, message, nil)
}

// Forbidden 返回禁止访问错误。
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrForbidden.Message
	}
	write(c, errors.ErrForbidden.Code, message, nil)
}

// NotFound 返回资源不存在错误。
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrNotFound.Message
	}
	write(c, errors.ErrNotFound.Code, message, nil)
}

// Conflict 返回资源冲突错误。
func Conflict(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrConflict.Message
	}
	write(c, errors.ErrConflict.Code, message, nil)
}

// InternalError 返回内部错误。
func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrInternal.Message
	}
	write(c, errors.ErrInternal.Code, message, nil)
}

// NotImplemented 返回功能未实现错误。
func NotImplemented(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrNotImplemented.Message
	}
	write(c, errors.ErrNotImplemented.Code, message, nil)
}

// TooManyRequests 返回请求过于频繁错误。
func TooManyRequests(c *gin.Context, message string) {
	if message == "" {
		message = errors.ErrTooManyRequests.Message
	}
	write(c, errors.ErrTooManyRequests.Code, message, nil)
}
