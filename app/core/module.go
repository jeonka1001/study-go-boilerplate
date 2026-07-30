package core

import (
	"boilerplate/app/core/base"
	"boilerplate/app/core/dao"
	"boilerplate/app/core/repository"

	"go.uber.org/fx"
)

// Modules aggregates every cross-cutting dependency a controller/service can inject.
type Modules struct {
	fx.In

	Base       Base
	Repository Repository
}

type Base struct {
	fx.In

	Parameter base.Parameter
}

var BaseModule = fx.Module(
	"base",
	fx.Provide(
		base.NewGetParameter,
	),
)

// Repository groups every domain repository backed by the DB connection.
type Repository struct {
	fx.In

	User repository.UserRepository
}

var RepositoryModule = fx.Module(
	"repository",
	fx.Provide(
		dao.NewMySQL,
		repository.NewUserRepository,
	),
)
