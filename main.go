package main

import (
	"boilerplate/app"
	"boilerplate/app/config"
	"boilerplate/app/core"
	"boilerplate/app/core/helper"
	"boilerplate/app/core/helper/logger"
	"boilerplate/app/domain/user"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// @title boilerplate API
// @version 1.0
// @contact.name TODO
// @contact.email TODO
// @host localhost:8080
// @accept application/json
// @produce application/json
func main() {
	fx.New(
		config.Module,
		helper.Module,

		core.BaseModule,
		core.RepositoryModule,

		user.ControllerModule,
		user.ServiceModule,

		fx.Provide(
			app.NewFiber,
			fx.Annotate(
				app.NewRouter,
				fx.ParamTags(``, `group:"routes"`),
			),
		),
		fx.Invoke(
			app.NewMiddleware,
			func(*logger.Sugared) {},
			func(fiber.Router) {},
			func(*fiber.App) {},
		),
	).Run()
}
