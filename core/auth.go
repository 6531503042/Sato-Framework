package core

import (
	"github.com/gofiber/fiber/v2"
)

// AuthProvider defines the interface for authentication providers
type AuthProvider interface {
	Authenticate(ctx *fiber.Ctx) (interface{}, error)
}

// JWTProvider implements JWT authentication
type JWTProvider struct {
	Secret string
}

// NewJWTProvider creates a new JWT provider
func NewJWTProvider(secret string) *JWTProvider {
	return &JWTProvider{Secret: secret}
}

// Authenticate implements JWT authentication
func (p *JWTProvider) Authenticate(ctx *fiber.Ctx) (interface{}, error) {
	token := ctx.Get("Authorization")
	if token == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Missing token")
	}

	// TODO: Implement JWT validation
	return nil, nil
}

// Role represents a user role
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(provider AuthProvider) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user, err := provider.Authenticate(ctx)
		if err != nil {
			return err
		}

		ctx.Locals("user", user)
		return ctx.Next()
	}
}

// RoleMiddleware creates a role-based authorization middleware
func RoleMiddleware(roles ...Role) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user := ctx.Locals("user")
		if user == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		// TODO: Implement role checking
		return ctx.Next()
	}
}

// Permission represents a permission
type Permission string

// PermissionMiddleware creates a permission-based authorization middleware
func PermissionMiddleware(permissions ...Permission) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		user := ctx.Locals("user")
		if user == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}

		// TODO: Implement permission checking
		return ctx.Next()
	}
} 