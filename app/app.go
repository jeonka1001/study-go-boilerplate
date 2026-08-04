package app

import (
	"boilerplate/config"
	"context"
	"fmt"
	"net/http"

	"boilerplate/app/core/exception"
	"boilerplate/app/core/helper/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"go.uber.org/fx"
)

// NewFiber creates the Fiber app and wires it into the fx lifecycle.
func NewFiber(lc fx.Lifecycle, c *config.Config) *fiber.App {
	app := initializeFiber()

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Zap.Infof("Server is starting on port: %s", c.Server.Port)

			addr := fmt.Sprintf(":%s", c.Server.Port)
			go app.Listen(addr)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})

	return app
}

func initializeFiber() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:       "boilerplate",
		ServerHeader:  "boilerplate",
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		UnescapePath:  true,
		ErrorHandler:  exception.ErrorHandler,
	})

	app.Get("/check_health", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(http.StatusOK)
	})
	app.Get("/swagger/*", swagger.HandlerDefault)

	return app
}
