package app

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// NewRouter registers every Route's Mapping table onto the Fiber app.
func NewRouter(app *fiber.App, routes []Route) fiber.Router {
	for _, route := range routes {
		for _, mapping := range route.Table() {
			switch mapping.Method {
			case fiber.MethodGet:
				app.Get(mapping.Path, mapping.Handler)
			case fiber.MethodPost:
				app.Post(mapping.Path, mapping.Handler)
			case fiber.MethodPut:
				app.Put(mapping.Path, mapping.Handler)
			case fiber.MethodPatch:
				app.Patch(mapping.Path, mapping.Handler)
			case fiber.MethodDelete:
				app.Delete(mapping.Path, mapping.Handler)
			}
		}
	}

	return nil
}

// AsRoute annotates a controller constructor so fx collects it into the "routes" group.
func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}

type Route interface {
	Table() []Mapping
}

type Mapping struct {
	Method  string
	Path    string
	Handler func(ctx *fiber.Ctx) error
}
