package fiber

import (
	"github.com/godome/godome/pkg/component/module"
	"github.com/godome/godome/pkg/component/provider"
	"github.com/gofiber/fiber/v2"
)

const ProviderName = "FiberHandler"

type FiberHandler interface {
	provider.Provider
	AddRoute(func(*fiber.App))
	LoadRoutes(*fiber.App)
}

func NewFiberHandler(m module.Module) FiberHandler {
	return &fiberHandler{
		Provider: provider.NewProvider(ProviderName),
		module:   m,
	}
}

type fiberHandler struct {
	provider.Provider
	routes []func(*fiber.App)
	module module.Module
}

func (r *fiberHandler) AddRoute(newRoute func(*fiber.App)) {
	r.routes = append(r.routes, newRoute)
}

func (r *fiberHandler) LoadRoutes(app *fiber.App) {
	for _, route := range r.routes {
		route(app)
	}
}
