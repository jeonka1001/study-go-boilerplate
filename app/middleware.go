package app

import (
	"boilerplate/app/core"
	"boilerplate/app/core/exception"
	"boilerplate/app/core/helper/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap/zapcore"
)

// NewMiddleware registers every global middleware, in order:
// panic recover -> cors -> request-id -> request logging.
// 이 함수 하나에서 전역 미들웨어를 조립/등록한다. 개별 미들웨어의 구현 로직은
// core/exception, core/helper/logger 등 각 패키지에 두고 여기서는 등록만 한다.
func NewMiddleware(app *fiber.App, core core.Modules) {
	app.Use(exception.Recover())
	app.Use(cors.New())
	app.Use(requestid.New(requestid.Config{
		ContextKey: logger.RequestID,
	}))
	app.Use(logger.NewMiddleware(logger.Config{
		SkipURIs: []string{"/check_health"},
		Logger:   logger.Zap.Desugar(),
		Fields: []string{
			logger.RequestID,
			logger.Status,
			logger.Latency,
			logger.Error,
			logger.PID,
			logger.IP,
			logger.IPs,
			logger.Method,
			logger.URL,
			logger.Protocol,
		},
		Levels: []zapcore.Level{zapcore.ErrorLevel, zapcore.WarnLevel, zapcore.InfoLevel},
	}))

	// DB 조회가 필요한 인증/크리덴셜 미들웨어를 추가할 때는 여기서 core.Repository.* 를 사용:
	// app.Use(requireUser(core))
}
