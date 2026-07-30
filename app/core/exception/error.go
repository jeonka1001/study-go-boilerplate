package exception

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Code 는 "<httpStatus>.<detail>" 형식의 애플리케이션 에러 코드다.
// 앞자리에서 HTTP status 를 파생시키므로 새 코드를 추가할 때 형식을 지켜야 한다.
type Code string

const (
	CodeBadRequest          Code = "400.0001"
	CodeInvalidParameter    Code = "400.0002"
	CodeUnauthorized        Code = "401.0001"
	CodeForbidden           Code = "403.0001"
	CodeNotFound            Code = "404.0001"
	CodeConflict            Code = "409.0001"
	CodeInternalServerError Code = "500.0001"
)

// defaultMessages 는 코드별 기본 사용자 노출 메시지다.
// 내부 에러 상세(err.Error())는 5xx 응답에 노출하지 않는다.
var defaultMessages = map[Code]string{
	CodeBadRequest:          "Bad Request",
	CodeInvalidParameter:    "Invalid Parameter",
	CodeUnauthorized:        "Unauthorized",
	CodeForbidden:           "Forbidden",
	CodeNotFound:            "Not Found",
	CodeConflict:            "Conflict",
	CodeInternalServerError: "Internal Server Error",
}

// Error 는 HTTP 응답으로 변환 가능한 애플리케이션 에러다.
type Error struct {
	Code    Code
	Message string
	// Data 는 클라이언트에게 돌려줄 부가 정보(검증 실패 필드 등). 민감 정보를 담지 말 것.
	Data any
	// cause 는 원인 에러. errors.Is/As 로 추적 가능하며 응답에는 직접 노출되지 않는다.
	cause error
}

func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.cause
}

// StatusCode 는 Code 앞자리에서 HTTP status 를 파생시킨다.
func (e *Error) StatusCode() int {
	parts := strings.SplitN(string(e.Code), ".", 2)
	if len(parts) > 0 {
		if code, err := strconv.Atoi(parts[0]); err == nil && code >= 100 && code <= 599 {
			return code
		}
	}
	return http.StatusInternalServerError
}

// New 는 원인 에러 없이 새 Error 를 만든다.
func New(code Code, message string) *Error {
	if message == "" {
		message = messageOf(code)
	}
	return &Error{Code: code, Message: message}
}

// Wrap 은 원인 에러를 코드와 함께 감싼다.
func Wrap(code Code, err error) *Error {
	return &Error{Code: code, Message: messageOf(code), cause: err}
}

// WithData 는 원인 에러를 감싸면서 클라이언트에 전달할 부가 정보를 함께 담는다.
func WithData(code Code, err error, data any) *Error {
	return &Error{Code: code, Message: messageOf(code), Data: data, cause: err}
}

func messageOf(code Code) string {
	if msg, ok := defaultMessages[code]; ok {
		return msg
	}
	return "Unknown Error"
}
