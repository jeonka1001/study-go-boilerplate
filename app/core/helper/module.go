package helper

import (
	"boilerplate/app/core/helper/logger"
	"boilerplate/app/core/helper/validator"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"helper",
	fx.Provide(
		logger.New,
		validator.New,
	),
)

type Helper struct {
	fx.In

	Logger    *logger.Sugared
	Validator *validator.Checker
}
