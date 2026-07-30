package controller

import (
	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"

	"boilerplate/app"
	"boilerplate/app/core"
	"boilerplate/app/core/helper"
	"boilerplate/app/domain/user/dto"
	"boilerplate/app/domain/user/service"
)

type UserController interface {
	Table() []app.Mapping
	GetUser(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
}

type userController struct {
	core   core.Modules
	helper helper.Helper
	svc    service.UserService
}

func NewUserController(core core.Modules, helper helper.Helper, svc service.UserService) UserController {
	return userController{core: core, helper: helper, svc: svc}
}

func (ctrl userController) Table() []app.Mapping {
	return []app.Mapping{
		{Method: fiber.MethodGet, Path: "/v1/users/:id", Handler: ctrl.GetUser},
		{Method: fiber.MethodPost, Path: "/v1/users", Handler: ctrl.CreateUser},
	}
}

// @Summary 유저 조회
// @Tags User
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserRes
// @Router /v1/users/{id} [get]
func (ctrl userController) GetUser(c *fiber.Ctx) error {
	var req dto.GetUserReq
	if err := ctrl.core.Base.Parameter.GetRequest(c, &req); err != nil {
		return errors.WithStack(err)
	}

	var param dto.GetUserParam
	if err := ctrl.core.Base.Parameter.ValidateParams(c, req, &param); err != nil {
		return errors.WithStack(err)
	}

	res, err := ctrl.svc.GetUser(param)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

// @Summary 유저 생성
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.CreateUserReq true "User"
// @Success 201 {object} dto.UserRes
// @Router /v1/users [post]
func (ctrl userController) CreateUser(c *fiber.Ctx) error {
	var req dto.CreateUserReq
	if err := ctrl.core.Base.Parameter.GetRequest(c, &req); err != nil {
		return errors.WithStack(err)
	}

	var param dto.CreateUserParam
	if err := ctrl.core.Base.Parameter.ValidateParams(c, req, &param); err != nil {
		return errors.WithStack(err)
	}

	res, err := ctrl.svc.CreateUser(param)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(res)
}
