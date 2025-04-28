package core

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// HTTP Methods
const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
	PATCH  = "PATCH"
)

// ControllerOptions defines controller configuration
type ControllerOptions struct {
	Path    string
	Guards  []Guard
	Version string
}

// RouteOptions defines route configuration
type RouteOptions struct {
	Path   string
	Method string
}

// ControllerMeta stores controller metadata
type ControllerMeta struct {
	Instance interface{}
	Path     string
	Guards   []Guard
	Version  string
	Routes   []RouteMeta
}

// RouteMeta stores route metadata
type RouteMeta struct {
	Path     string
	Method   string
	Handler  string
	Guards   []Guard
	Pipes    []PipeMeta
}

var controllers []ControllerMeta

// Controller decorator for class
func Controller(options ControllerOptions) func(c interface{}) interface{} {
	return func(c interface{}) interface{} {
		path := options.Path
		if path == "" {
			path = "/" + strings.ToLower(reflect.TypeOf(c).Name())
		}

		controllers = append(controllers, ControllerMeta{
			Instance: c,
			Path:     path,
			Guards:   options.Guards,
			Version:  options.Version,
			Routes:   make([]RouteMeta, 0),
		})

		return c
	}
}

// mergeRouteOptions merges route options with defaults
func mergeRouteOptions(path string, method string, options ...RouteOptions) RouteOptions {
	opts := RouteOptions{
		Path:   path,
		Method: method,
	}
	if len(options) > 0 {
		// Merge with provided options
		if options[0].Path != "" {
			opts.Path = options[0].Path
		}
		if options[0].Method != "" {
			opts.Method = options[0].Method
		}
	}
	return opts
}

// Get registers a GET route
func Get(path string, options ...RouteOptions) func(interface{}, string, fiber.Handler) {
	return func(controller interface{}, handlerName string, handlerFunc fiber.Handler) {
		opts := mergeRouteOptions(path, "GET", options...)
		addRoute(controller, RouteMeta{
			Path:    opts.Path,
			Method:  GET,
			Handler: handlerName,
		})
	}
}

// Post registers a POST route
func Post(path string, options ...RouteOptions) func(interface{}, string, fiber.Handler) {
	return func(controller interface{}, handlerName string, handlerFunc fiber.Handler) {
		opts := mergeRouteOptions(path, "POST", options...)
		addRoute(controller, RouteMeta{
			Path:    opts.Path,
			Method:  POST,
			Handler: handlerName,
		})
	}
}

// Put registers a PUT route
func Put(path string, options ...RouteOptions) func(interface{}, string, fiber.Handler) {
	return func(controller interface{}, handlerName string, handlerFunc fiber.Handler) {
		opts := mergeRouteOptions(path, "PUT", options...)
		addRoute(controller, RouteMeta{
			Path:    opts.Path,
			Method:  PUT,
			Handler: handlerName,
		})
	}
}

// Delete registers a DELETE route
func Delete(path string, options ...RouteOptions) func(interface{}, string, fiber.Handler) {
	return func(controller interface{}, handlerName string, handlerFunc fiber.Handler) {
		opts := mergeRouteOptions(path, "DELETE", options...)
		addRoute(controller, RouteMeta{
			Path:    opts.Path,
			Method:  DELETE,
			Handler: handlerName,
		})
	}
}

// Patch registers a PATCH route
func Patch(path string, options ...RouteOptions) func(interface{}, string, fiber.Handler) {
	return func(controller interface{}, handlerName string, handlerFunc fiber.Handler) {
		opts := mergeRouteOptions(path, "PATCH", options...)
		addRoute(controller, RouteMeta{
			Path:    opts.Path,
			Method:  PATCH,
			Handler: handlerName,
		})
	}
}

// UseGuards decorator for method
func UseGuards(guards ...Guard) func(interface{}, string) {
	return func(target interface{}, propertyKey string) {
		addGuards(target, propertyKey, guards)
	}
}

// UsePipes decorator for method
func UsePipes(pipes ...PipeMeta) func(interface{}, string) {
	return func(target interface{}, propertyKey string) {
		addPipes(target, propertyKey, pipes)
	}
}

func addRoute(target interface{}, route RouteMeta) {
	for i, c := range controllers {
		if c.Instance == target {
			controllers[i].Routes = append(controllers[i].Routes, route)
			break
		}
	}
}

func addGuards(target interface{}, handler string, guards []Guard) {
	for i, c := range controllers {
		if c.Instance == target {
			for j, r := range c.Routes {
				if r.Handler == handler {
					controllers[i].Routes[j].Guards = guards
					break
				}
			}
			break
		}
	}
}

func addPipes(target interface{}, handler string, pipes []PipeMeta) {
	for i, c := range controllers {
		if c.Instance == target {
			for j, r := range c.Routes {
				if r.Handler == handler {
					controllers[i].Routes[j].Pipes = pipes
					break
				}
			}
			break
		}
	}
}

// GetControllers returns all registered controllers
func GetControllers() []ControllerMeta {
	return controllers
}

// ApplyPipes applies pipes to a value
func ApplyPipes(value interface{}, pipes []PipeMeta) (interface{}, error) {
	var err error
	result := value
	
	for _, pipe := range pipes {
		result, err = pipe.Pipe(result)
		if err != nil {
			return nil, err
		}
	}
	
	return result, nil
}

// RegisterRoute adds a new route to the controller's metadata
func RegisterRoute(controller interface{}, opts RouteOptions, handler fiber.Handler) {
	if c, ok := controller.(interface{ GetMeta() *ControllerMeta }); ok {
		meta := c.GetMeta()
		route := RouteMeta{
			Path:    opts.Path,
			Method:  opts.Method,
			Handler: runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(),
			Guards:  make([]Guard, 0),
		}
		meta.Routes = append(meta.Routes, route)
	}
}

// RegisterGuard adds a guard to the controller's metadata
func RegisterGuard(controller interface{}, guard Guard) {
	if c, ok := controller.(interface{ GetMeta() *ControllerMeta }); ok {
		meta := c.GetMeta()
		meta.Guards = append(meta.Guards, guard)
	}
}

// validateControllerMeta validates the controller metadata
func validateControllerMeta(meta *ControllerMeta) error {
	if meta.Path == "" {
		return fmt.Errorf("controller path cannot be empty")
	}
	if !strings.HasPrefix(meta.Path, "/") {
		return fmt.Errorf("controller path must start with /")
	}
	return nil
}

// findRoute finds a route by path and method
func findRoute(routes []RouteMeta, path string, method string) *RouteMeta {
	for _, r := range routes {
		if r.Path == path && r.Method == method {
			return &r
		}
	}
	return nil
}

