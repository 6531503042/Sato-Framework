package core

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

// TestController for testing
type TestController struct {
	meta ControllerMeta
}

func (c *TestController) GetMeta() *ControllerMeta {
	return &c.meta
}

func TestRouteRegistration(t *testing.T) {
	app := fiber.New()
	controller := &TestController{
		meta: ControllerMeta{
			Path:    "/test",
			Version: "v1",
			Routes:  make([]RouteMeta, 0),
		},
	}

	// Add test routes
	addRoute(controller, RouteMeta{
		Path:    "/",
		Method:  "GET",
		Handler: "GetHandler",
	})

	addRoute(controller, RouteMeta{
		Path:    "/",
		Method:  "POST",
		Handler: "PostHandler",
	})

	// Register routes
	RegisterRoutes(app)

	// Test route registration
	routes := app.GetRoutes()
	if len(routes) == 0 {
		t.Error("No routes were registered")
	}

	// Verify GET route
	found := false
	for _, route := range routes {
		if route.Path == "/test" && route.Method == "GET" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GET /test route was not registered")
	}

	// Verify POST route
	found = false
	for _, route := range routes {
		if route.Path == "/test" && route.Method == "POST" {
			found = true
			break
		}
	}
	if !found {
		t.Error("POST /test route was not registered")
	}
}

func TestGuardRegistration(t *testing.T) {
	app := fiber.New()
	controller := &TestController{
		meta: ControllerMeta{
			Path:    "/test",
			Version: "v1",
			Routes:  make([]RouteMeta, 0),
			Guards:  make([]Guard, 0),
		},
	}

	// Create a test guard
	testGuard := &TestGuard{}
	controller.meta.Guards = append(controller.meta.Guards, testGuard)

	// Add test route
	addRoute(controller, RouteMeta{
		Path:    "/",
		Method:  "GET",
		Handler: "GetHandler",
	})

	// Register routes
	RegisterRoutes(app)

	// Test that routes are registered
	routes := app.GetRoutes()
	if len(routes) == 0 {
		t.Error("No routes were registered with guards")
	}
} 