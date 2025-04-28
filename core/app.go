package core

import (
	"github.com/gofiber/fiber/v2"
)

type App struct {
	server *fiber.App
}

func NewApp() *App {
	app := fiber.New()

	app.Use(GlobalErrorHandler())

	return &App{server: app}
}

func (a *App) GetFiber() *fiber.App {
	return a.server
}

func (a *App) Listen(addr string) error {
	return a.server.Listen(addr)
}
