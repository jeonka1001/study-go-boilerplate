package dto

import (
	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
)

type GetUserReq struct {
	ID uint64 `params:"id"`
}

type GetUserParam struct {
	ID uint64 `validate:"required"`
}

func (p *GetUserParam) GenerateParam(ctx *fiber.Ctx, req interface{}) error {
	r, ok := req.(GetUserReq)
	if !ok {
		return errors.New("invalid request type")
	}
	p.ID = r.ID
	return nil
}

type CreateUserReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserParam struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func (p *CreateUserParam) GenerateParam(ctx *fiber.Ctx, req interface{}) error {
	r, ok := req.(CreateUserReq)
	if !ok {
		return errors.New("invalid request type")
	}
	p.Name = r.Name
	p.Email = r.Email
	return nil
}
