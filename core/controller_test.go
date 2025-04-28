package core

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

type MockController struct {
	meta ControllerMeta
}

func (c *MockController) GetMeta() *ControllerMeta {
	return &c.meta
}

func TestControllerRegistration(t *testing.T) {
	controller := &MockController{
		meta: ControllerMeta{
			Path:    "/api",
			Version: "v1",
			Routes:  make([]RouteMeta, 0),
		},
	}

	// Test route registration
	routeOpts := RouteOptions{
		Path:   "/users",
		Method: "GET",
	}

	handler := func(c *fiber.Ctx) error {
		return c.SendString("GET users")
	}

	RegisterRoute(controller, routeOpts, handler)

	if len(controller.GetMeta().Routes) != 1 {
		t.Error("Expected one route to be registered")
	}

	route := controller.GetMeta().Routes[0]
	if route.Path != "/users" {
		t.Errorf("Expected path /users, got %s", route.Path)
	}

	if route.Method != "GET" {
		t.Errorf("Expected method GET, got %s", route.Method)
	}

	// Test guard registration
	guard := &TestGuard{}
	RegisterGuard(controller, guard)

	if len(controller.GetMeta().Guards) != 1 {
		t.Error("Expected one guard to be registered")
	}
}

type TestGuard struct{}

func (g *TestGuard) CanActivate(c *fiber.Ctx) error {
	return c.Next()
}

func TestControllerValidation(t *testing.T) {
	controller := &MockController{
		meta: ControllerMeta{},
	}

	// Test invalid path
	if err := validateControllerMeta(&controller.meta); err == nil {
		t.Error("Expected error for empty path")
	}

	// Test valid path
	controller.meta.Path = "/api"
	if err := validateControllerMeta(&controller.meta); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
} 