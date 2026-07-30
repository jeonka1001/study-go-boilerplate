package user

import (
	"boilerplate/app"
	"boilerplate/app/domain/user/controller"
	"boilerplate/app/domain/user/service"

	"go.uber.org/fx"
)

var ControllerModule = fx.Module(
	"user-controller",
	fx.Provide(
		app.AsRoute(controller.NewUserController),
	),
)

var ServiceModule = fx.Module(
	"user-service",
	fx.Provide(
		service.NewUserService,
	),
)
