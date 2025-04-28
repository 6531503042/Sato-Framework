package core

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers all controller routes
func RegisterRoutes(app *fiber.App) {
	for _, controller := range GetControllers() {
		// Create controller group
		group := app.Group(controller.Path)

		// Apply controller-level guards
		if len(controller.Guards) > 0 {
			group.Use(func(c *fiber.Ctx) error {
				for _, guard := range controller.Guards {
					if err := guard.CanActivate(c); err != nil {
						return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
							"error": "Unauthorized",
						})
					}
				}
				return c.Next()
			})
		}

		// Register routes
		for _, route := range controller.Routes {
			handler := reflect.ValueOf(controller.Instance).MethodByName(route.Handler)
			if !handler.IsValid() {
				continue
			}

			// Apply route-level guards
			routeHandler := func(c *fiber.Ctx) error {
				if len(route.Guards) > 0 {
					for _, guard := range route.Guards {
						if err := guard.CanActivate(c); err != nil {
							return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
								"error": "Unauthorized",
							})
						}
					}
				}

				// Call the handler
				result := handler.Call([]reflect.Value{reflect.ValueOf(c)})
				if len(result) > 0 {
					return result[0].Interface().(error)
				}
				return nil
			}

			// Register route with method
			switch route.Method {
			case GET:
				group.Get(route.Path, routeHandler)
			case POST:
				group.Post(route.Path, routeHandler)
			case PUT:
				group.Put(route.Path, routeHandler)
			case DELETE:
				group.Delete(route.Path, routeHandler)
			case PATCH:
				group.Patch(route.Path, routeHandler)
			}
		}
	}
}
