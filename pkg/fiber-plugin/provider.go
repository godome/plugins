package fiberPlugin

import (
	"github.com/godome/godome/pkg/component/provider"
	"github.com/gofiber/fiber/v2"
)

const ProviderName = "FiberHandler"

type FiberHandler interface {
	provider.Provider
	AddRoute(func(*fiber.App)) FiberHandler
	LoadRoutes(*fiber.App)
}

func NewFiberHandler() FiberHandler {
	return &fiberHandler{
		Provider: provider.NewProvider(ProviderName),
	}
}

type fiberHandler struct {
	provider.Provider
	routes []func(*fiber.App)
}

func (r *fiberHandler) AddRoute(newRoute func(*fiber.App)) FiberHandler {
	r.routes = append(r.routes, newRoute)
	return r
}

func (r *fiberHandler) LoadRoutes(app *fiber.App) {
	for _, route := range r.routes {
		route(app)
	}
}
