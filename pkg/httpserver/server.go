package httpserver

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 5 * time.Second
	defaultAddr         = ":3120"
)

type Server struct {
	server *fiber.App
	port   string
}

func FiberConfig(status, appName string) fiber.Config {
	config := fiber.Config{
		AppName:                   appName,
		ServerHeader:              appName + " API",
		DisableStartupMessage:     true,
		DisableDefaultContentType: false,
		EnablePrintRoutes:         false,
		ReadTimeout:               defaultReadTimeout,
		WriteTimeout:              defaultWriteTimeout,
		IdleTimeout:               60 * time.Second,
		RequestMethods:            []string{fiber.MethodGet, fiber.MethodHead, fiber.MethodOptions},
	}

	if status == "dev" {
		config.DisableStartupMessage = false
		config.EnablePrintRoutes = true
	}

	return config
}

func New(app *fiber.App, opts ...Option) *Server {
	s := &Server{
		server: app,
		port:   defaultAddr,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.start()
	return s
}

func (s *Server) start() {
	go func() {
		if err := s.server.Listen(s.port); err != nil {
			log.Fatalf("server shutdown error - %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c
	log.Println("Gracefully shutting down...")
	_ = s.server.Shutdown()
	log.Println("App was successful shutdown.")
}
