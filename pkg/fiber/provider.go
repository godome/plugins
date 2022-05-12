package fiber

import (
	"github.com/godome/godome/pkg/module"
	"github.com/godome/godome/pkg/provider"
	"github.com/gofiber/fiber/v2"
)

const ProviderType provider.ProviderType = "FiberHandler"

type FiberHandler interface {
	provider.Provider
	AddRoute(func(*fiber.App))
	LoadRoutes(*fiber.App)
}

type fiberHandler struct {
	routes       []func(*fiber.App)
	module       module.Module
	providerType provider.ProviderType
}

func NewFiberHandler(m module.Module) FiberHandler {
	return &fiberHandler{
		module:       m,
		providerType: ProviderType,
	}
}

func (r *fiberHandler) GetType() provider.ProviderType {
	return r.providerType
}

func (r *fiberHandler) AddRoute(newRoute func(*fiber.App)) {
	r.routes = append(r.routes, newRoute)
}

func (r *fiberHandler) LoadRoutes(app *fiber.App) {
	for _, route := range r.routes {
		route(app)
	}
}
