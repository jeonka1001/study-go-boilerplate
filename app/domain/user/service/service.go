package service

import (
	"github.com/cockroachdb/errors"

	"boilerplate/app/core"
	"boilerplate/app/core/repository"
	"boilerplate/app/domain/user/dto"
)

type UserService interface {
	GetUser(param dto.GetUserParam) (*dto.UserRes, error)
	CreateUser(param dto.CreateUserParam) (*dto.UserRes, error)
}

type userService struct {
	core core.Modules
}

func NewUserService(core core.Modules) UserService {
	return userService{core: core}
}

func (s userService) GetUser(param dto.GetUserParam) (*dto.UserRes, error) {
	user, err := s.core.Repository.User.FindByID(param.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return toUserRes(user), nil
}

func (s userService) CreateUser(param dto.CreateUserParam) (*dto.UserRes, error) {
	user := &repository.User{Name: param.Name, Email: param.Email}
	if err := s.core.Repository.User.Create(user); err != nil {
		return nil, errors.WithStack(err)
	}
	return toUserRes(user), nil
}

func toUserRes(user *repository.User) *dto.UserRes {
	return &dto.UserRes{ID: user.ID, Name: user.Name, Email: user.Email}
}
