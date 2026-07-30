package exception

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"

	"boilerplate/app/core/helper/logger"
)

// ErrorHandler converts any error returned from a handler/middleware into a JSON
// response. 5xx 는 서버 로그에 상세를 남기고, 응답에는 내부 상세를 노출하지 않는다.
func ErrorHandler(c *fiber.Ctx, err error) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	statusCode, body := determineErrorResponse(err)
	if statusCode >= fiber.StatusInternalServerError && logger.Zap != nil {
		logger.Zap.Errorw("unhandled server error",
			"method", c.Method(),
			"path", c.Path(),
			"requestId", c.GetRespHeader(fiber.HeaderXRequestID),
			"error", err.Error(),
		)
	}

	return c.Status(statusCode).JSON(body)
}

func determineErrorResponse(err error) (int, fiber.Map) {
	var appErr *Error
	if errors.As(err, &appErr) {
		body := fiber.Map{
			"code":  appErr.Code,
			"error": appErr.Message,
		}
		if appErr.Data != nil {
			body["data"] = appErr.Data
		}
		return appErr.StatusCode(), body
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return fiberErr.Code, fiber.Map{
			"code":  CodeFromStatus(fiberErr.Code),
			"error": fiberErr.Message,
		}
	}

	// 정체를 모르는 에러는 내부 상세를 감추고 500 으로 응답한다.
	return fiber.StatusInternalServerError, fiber.Map{
		"code":  CodeInternalServerError,
		"error": messageOf(CodeInternalServerError),
	}
}

// CodeFromStatus maps a bare HTTP status onto the closest application Code.
func CodeFromStatus(status int) Code {
	switch status {
	case fiber.StatusBadRequest:
		return CodeBadRequest
	case fiber.StatusUnauthorized:
		return CodeUnauthorized
	case fiber.StatusForbidden:
		return CodeForbidden
	case fiber.StatusNotFound:
		return CodeNotFound
	case fiber.StatusConflict:
		return CodeConflict
	default:
		return CodeInternalServerError
	}
}

// Recover is a panic-recovery middleware. 패닉을 error 로 변환해 ErrorHandler 로 넘긴다.
func Recover(config ...fiberRecover.Config) fiber.Handler {
	cfg := fiberRecover.ConfigDefault
	if len(config) > 0 {
		cfg = config[0]
	}
	if cfg.EnableStackTrace && cfg.StackTraceHandler == nil {
		cfg.StackTraceHandler = defaultStackTraceHandler
	}

	return func(c *fiber.Ctx) (err error) {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		defer func() {
			if r := recover(); r != nil {
				if cfg.EnableStackTrace {
					cfg.StackTraceHandler(c, r)
				}

				cause, ok := r.(error)
				if !ok {
					cause = fmt.Errorf("%v", r)
				}

				if logger.Zap != nil {
					logger.Zap.Errorw("recovered from panic",
						"method", c.Method(),
						"path", c.Path(),
						"panic", cause.Error(),
						"stack", string(debug.Stack()),
					)
				}

				err = Wrap(CodeInternalServerError, cause)
			}
		}()

		return c.Next()
	}
}

func defaultStackTraceHandler(_ *fiber.Ctx, e interface{}) {
	_, _ = os.Stderr.WriteString(fmt.Sprintf("panic: %v\n%s\n", e, debug.Stack())) //nolint:errcheck
}
