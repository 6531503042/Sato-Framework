package core

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

func Learn(model interface{}, prefix string, app *fiber.App) {
	modelName := reflect.TypeOf(model).Name()
	group := app.Group(prefix + "/" + modelName)

	group.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("Created " + modelName)
	})
	group.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Listed " + modelName)
	})
	group.Get("/:id", func(c *fiber.Ctx) error {
		return c.SendString("Single " + modelName)
	})
	group.Put("/:id", func(c *fiber.Ctx) error {
		return c.SendString("Updated " + modelName)
	})
	group.Delete("/:id", func(c *fiber.Ctx) error {
		return c.SendString("Deleted " + modelName)
	})
}
