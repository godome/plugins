package fiber

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/godome/godome/pkg/component/exposure"
	"github.com/godome/godome/pkg/component/module"
	"github.com/godome/godome/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

const ExposureName = "FiberExposure"

type FiberExposure interface {
	exposure.Exposure
	ExposeModule(module module.Module) FiberExposure
	Run() error
	Test(req *http.Request, msTimeout ...int) (resp *http.Response, err error)
	loadRoutes() error
}

type fiberExposure struct {
	exposure.Exposure
	modules map[string]module.Module
	port    string
	app     *fiber.App
}

func NewFiberExposure(port string, config *fiber.Config) FiberExposure {
	fiberConfig := fiber.Config{}
	if config != nil {
		fiberConfig = *config
	}

	app := fiber.New(fiberConfig)

	return &fiberExposure{
		Exposure: exposure.NewExposure(ExposureName),
		modules:  make(map[string]module.Module),
		port:     port,
		app:      app,
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
	r.loadRoutes()

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
	// Load Routes
	r.loadRoutes()
	return r.app.Test(req, msTimeout...)
}

func (r *fiberExposure) ExposeModule(module module.Module) FiberExposure {
	r.modules[module.Metadata().GetName()] = module
	return r
}

func (r *fiberExposure) loadRoutes() error {
	for _, module := range r.modules {
		provider := module.GetProvider(ProviderName)
		if provider == nil {
			logger.Debug("FiberHandler provider not found on module: %s \n", module.Metadata().GetName())
			continue
		}
		foundFiberHandler, ok := provider.(FiberHandler)
		if !ok {
			return fmt.Errorf("FiberHandler provider could not been casted on module: %s", module.Metadata().GetName())
		}
		foundFiberHandler.LoadRoutes(r.app)
	}
	return nil
}
