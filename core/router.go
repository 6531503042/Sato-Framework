package core

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	for _, controller := range GetControllers() {
		typ := reflect.TypeOf(controller.Instance)
		val := reflect.ValueOf(controller.Instance)

		group := app.Group(controller.RoutePrefix)

		for i := 0; i < val.NumMethod(); i++ {
			method := typ.Method(i)
			path := "/" + method.Name

			group.All(path, func(c *fiber.Ctx) error {
				for _, guard := range controller.Guards {
					if err := guard.CanActivate(c); err != nil {
						return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
							"error": "Unauthorized",
						})
					}
				}

				result := val.MethodByName(method.Name).Call([]reflect.Value{reflect.ValueOf(c)})
				if len(result) > 0 {
					return result[0].Interface().(error)
				}
				return nil
			})
		}
	}
}
