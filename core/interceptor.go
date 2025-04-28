package core

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// Interceptor is a function that can intercept requests
type Interceptor func(*fiber.Ctx, fiber.Handler) error

// InterceptorOptions defines interceptor configuration
type InterceptorOptions struct {
	Path string
}

// InterceptorMeta stores interceptor metadata
type InterceptorMeta struct {
	Path        string
	Interceptor Interceptor
}

var interceptors []InterceptorMeta

// UseInterceptor decorator for method
func UseInterceptor(interceptor Interceptor, options InterceptorOptions) func(interface{}, string) {
	return func(target interface{}, propertyKey string) {
		interceptors = append(interceptors, InterceptorMeta{
			Path:        options.Path,
			Interceptor: interceptor,
		})
	}
}

// GetInterceptors returns all registered interceptors
func GetInterceptors() []InterceptorMeta {
	return interceptors
}

// LoggingInterceptor example interceptor
func LoggingInterceptor() Interceptor {
	return func(c *fiber.Ctx, next fiber.Handler) error {
		// Log request
		log.Printf("Request: %s %s", c.Method(), c.Path())

		// Call next handler
		err := next(c)
		if err != nil {
			log.Printf("Error: %v", err)
		}

		// Log response
		log.Printf("Response: %d", c.Response().StatusCode())

		return err
	}
}

// ErrorInterceptor example interceptor
func ErrorInterceptor() Interceptor {
	return func(c *fiber.Ctx, next fiber.Handler) error {
		err := next(c)
		if err != nil {
			// Handle error
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return nil
	}
}

// AuthInterceptor example interceptor
func AuthInterceptor() Interceptor {
	return func(c *fiber.Ctx, next fiber.Handler) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Validate token
		// ...

		return next(c)
	}
} 