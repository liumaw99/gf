// Package errors 提供业务错误体系，支持错误码、HTTP 状态码、内部错误包装。
//
// 业务错误码规则：
//   - 0        = 成功
//   - 1001-1099 = 通用错误（参数、认证、权限等）
//   - 各业务模块自行扩展（如 2001+ 认证模块）
package errors

import (
	"errors"
	"fmt"
)

// AppError 业务错误，包含对外展示信息和内部错误链。
type AppError struct {
	Code    int    // 业务错误码
	Status  int    // HTTP 状态码
	Message string // 对外展示的消息
	err     error  // 内部错误，不暴露给客户端
}

// Error 实现 error 接口。
func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 返回内部错误，支持 errors.Is / errors.As 链式追踪。
func (e *AppError) Unwrap() error {
	return e.err
}

// Internal 返回内部原始错误（供日志记录）。
func (e *AppError) Internal() error {
	return e.err
}

// New 创建一个新的业务错误。
func New(code, status int, message string) *AppError {
	return &AppError{
		Code:    code,
		Status:  status,
		Message: message,
	}
}

// Wrap 将底层 error 包装为业务错误。
func Wrap(err error, code, status int, message string) *AppError {
	return &AppError{
		Code:    code,
		Status:  status,
		Message: message,
		err:     err,
	}
}

// WrapInternal 将任意 error 包装为 500 内部错误。
func WrapInternal(err error) *AppError {
	return Wrap(err, ErrInternal.Code, ErrInternal.Status, ErrInternal.Message)
}

// WithMessage 基于现有错误创建副本，替换消息。
func WithMessage(err *AppError, message string) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		Code:    err.Code,
		Status:  err.Status,
		Message: message,
		err:     err.err,
	}
}

// WithStatus 基于现有错误创建副本，替换 HTTP 状态码。
func WithStatus(err *AppError, status int) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		Code:    err.Code,
		Status:  status,
		Message: err.Message,
		err:     err.err,
	}
}

// IsAppError 判断 err 是否为 *AppError。
func IsAppError(err error) (*AppError, bool) {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

// 预定义通用业务错误。
var (
	ErrBadRequest      = New(1001, 400, "请求参数错误")
	ErrUnauthorized    = New(1002, 401, "未授权")
	ErrForbidden       = New(1003, 403, "禁止访问")
	ErrNotFound        = New(1004, 404, "资源不存在")
	ErrConflict        = New(1005, 409, "资源冲突")
	ErrInternal        = New(1006, 500, "服务器内部错误")
	ErrNotImplemented  = New(1007, 501, "功能尚未实现")
	ErrTooManyRequests = New(1008, 429, "请求过于频繁")
)
