package core

import "github.com/gofiber/fiber/v2"

type (
	Guard interface {
		CanActivate(c *fiber.Ctx) error
	}
)