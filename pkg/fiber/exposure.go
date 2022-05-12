package fiber

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/godome/godome/pkg/exposure"
	"github.com/godome/godome/pkg/logger"
	"github.com/godome/godome/pkg/module"
	"github.com/gofiber/fiber/v2"
)

const ExposureType exposure.ExposureType = "FiberExposure"

type FiberExposure interface {
	exposure.Exposure
	ExposeModule(module module.Module) FiberExposure
	Run() error
	Test(req *http.Request, msTimeout ...int) (resp *http.Response, err error)
}

type fiberExposure struct {
	modules      map[string]module.Module
	exposureType exposure.ExposureType
	port         string
	app          *fiber.App
}

func NewFiberExposure(port string) FiberExposure {
	app := fiber.New(fiber.Config{
		// DisableStartupMessage: true,
	})

	return &fiberExposure{
		modules:      make(map[string]module.Module),
		exposureType: ExposureType,
		port:         port,
		app:          app,
	}
}

func (r *fiberExposure) Run() error {
	// Shutdown gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logger.Info("\nshutting down...")
		err := r.app.Shutdown()
		if err != nil {
			os.Exit(0)
		}
	}()

	// Load Routes
	for _, module := range r.modules {
		provider := module.GetProvider(ProviderType)
		if provider == nil {
			logger.Debug("FiberHandler provider not found on module: %s \n", module.GetName())
			continue
		}

		foundFiberHandler, ok := provider.(FiberHandler)
		if !ok {
			return fmt.Errorf("FiberHandler provider could not been casted on module: %s", module.GetName())
		}

		foundFiberHandler.LoadRoutes(r.app)
	}

	// Startup message
	// startupMessage := fmt.Sprint(
	// 	fmt.Sprintln("                                                       "),
	// 	fmt.Sprintln("\033[36m┌───────────────────────────────────────────────┐"),
	// 	fmt.Sprintln("│                                               │"),
	// 	fmt.Sprintf("│             http://127.0.0.1:%s             │\n", r.port),
	// 	fmt.Sprintf("│     (bound on host 0.0.0.0 and port %s)     │\n", r.port),
	// 	fmt.Sprintln("│                                               │"),
	// 	fmt.Sprintln("│               To exit, press ^C               │"),
	// 	fmt.Sprintln("└───────────────────────────────────────────────┘\033[97m"),
	// )
	// logger.Info(startupMessage)

	// Start
	if err := r.app.Listen(":" + r.port); err != nil {
		return err
	}

	return nil
}

func (r *fiberExposure) Test(req *http.Request, msTimeout ...int) (resp *http.Response, err error) {
	return r.app.Test(req, msTimeout...)
}

func (r *fiberExposure) GetType() exposure.ExposureType {
	return r.exposureType
}

func (r *fiberExposure) ExposeModule(module module.Module) FiberExposure {
	r.modules[module.GetName()] = module
	return r
}
