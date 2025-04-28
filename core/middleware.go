package core

import (
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// MiddlewareFunc is a function that can be used as middleware
type MiddlewareFunc func(*fiber.Ctx) error

// MiddlewareChain represents a chain of middleware functions
type MiddlewareChain struct {
	middlewares []MiddlewareFunc
	mu          sync.RWMutex
}

// NewMiddlewareChain creates a new middleware chain
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]MiddlewareFunc, 0),
	}
}

// Use adds a middleware to the chain
func (c *MiddlewareChain) Use(middleware MiddlewareFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.middlewares = append(c.middlewares, middleware)
}

// Then executes the middleware chain
func (c *MiddlewareChain) Then(handler fiber.Handler) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Execute middleware chain
		for _, middleware := range c.middlewares {
			if err := middleware(ctx); err != nil {
				return err
			}
		}

		// Execute the final handler
		return handler(ctx)
	}
}

// Compose creates a new middleware chain from multiple middleware functions
func Compose(middlewares ...MiddlewareFunc) *MiddlewareChain {
	chain := NewMiddlewareChain()
	for _, middleware := range middlewares {
		chain.Use(middleware)
	}
	return chain
}

// GlobalErrorHandler creates a global error handling middleware
func GlobalErrorHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()
		if err != nil {
			// Log the error
			fmt.Printf("Error: %v\n", err)

			// Return appropriate error response
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return nil
	}
}

// ApplyMiddleware applies middleware to a fiber app
func ApplyMiddleware(app *fiber.App, chain *MiddlewareChain) {
	app.Use(chain.Then(func(ctx *fiber.Ctx) error {
		return ctx.Next()
	}))
}

// CreateMiddleware creates a new middleware function
func CreateMiddleware(fn func(*fiber.Ctx) error) MiddlewareFunc {
	return fn
}
