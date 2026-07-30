package base

import (
	"github.com/gofiber/fiber/v2"

	"boilerplate/app/core/exception"
	"boilerplate/app/core/helper"
)

type Parameter interface {
	GetRequest(ctx *fiber.Ctx, param interface{}) error
	ValidateParams(ctx *fiber.Ctx, req interface{}, param Param) error
}

type getParameter struct {
	helper helper.Helper
}

func NewGetParameter(helper helper.Helper) Parameter {
	return getParameter{helper: helper}
}

// GetRequest parses path params, headers, query string and (optionally) the body into param.
func (gp getParameter) GetRequest(ctx *fiber.Ctx, param interface{}) error {
	if err := ctx.ParamsParser(param); err != nil {
		return exception.WithData(exception.CodeInvalidParameter, err, fiber.Map{"error": err.Error()})
	}
	if err := ctx.ReqHeaderParser(param); err != nil {
		return exception.WithData(exception.CodeInvalidParameter, err, fiber.Map{"error": err.Error()})
	}
	if err := ctx.QueryParser(param); err != nil {
		return exception.WithData(exception.CodeInvalidParameter, err, fiber.Map{"error": err.Error()})
	}
	if len(ctx.Request().Body()) > 0 {
		if err := ctx.BodyParser(param); err != nil {
			return exception.WithData(exception.CodeInvalidParameter, err, fiber.Map{"error": err.Error()})
		}
	}
	return nil
}

// ValidateParams converts the raw request DTO into the domain Param and validates it.
func (gp getParameter) ValidateParams(ctx *fiber.Ctx, req interface{}, param Param) error {
	if err := param.GenerateParam(ctx, req); err != nil {
		return exception.Wrap(exception.CodeInternalServerError, err)
	}
	if err := gp.helper.Validator.Struct(param); err != nil {
		return exception.WithData(exception.CodeBadRequest, err, fiber.Map{"error": err.Error()})
	}
	return nil
}

type Param interface {
	GenerateParam(ctx *fiber.Ctx, req interface{}) error
}
